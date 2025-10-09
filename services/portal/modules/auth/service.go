package auth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/config"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
)

type repository interface {
	GetUserNIPByIDAndSource(ctx context.Context, arg sqlc.GetUserNIPByIDAndSourceParams) (string, error)
	ListUserRoleByNIP(ctx context.Context, nip string) ([]sqlc.ListUserRoleByNIPRow, error)
	UpdateLastLoginAt(ctx context.Context, arg sqlc.UpdateLastLoginAtParams) error
}

type service struct {
	repo       repository
	keycloak   config.Keycloak
	client     *http.Client
	privateKey *rsa.PrivateKey
	keyfunc    jwt.Keyfunc
}

func newService(repo repository, keycloak config.Keycloak, client *http.Client, privateKey *rsa.PrivateKey, keyfunc jwt.Keyfunc) *service {
	return &service{
		repo:       repo,
		keycloak:   keycloak,
		client:     client,
		privateKey: privateKey,
		keyfunc:    keyfunc,
	}
}

func (s *service) generateAuthURL(redirectURI string) (string, error) {
	authURL, err := url.Parse(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth", s.keycloak.PublicHost, s.keycloak.Realm))
	if err != nil {
		return "", fmt.Errorf("url parse: %w", err)
	}

	if redirectURI == "" {
		redirectURI = s.keycloak.RedirectURI
	}

	urlParams := authURL.Query() // nosemgrep: rules.go.sql.go_sql_rule-concat-sqli - false positive, this is not SQL
	urlParams.Set("client_id", s.keycloak.ClientID)
	urlParams.Set("response_type", "code")
	urlParams.Set("scope", "openid")
	urlParams.Set("redirect_uri", redirectURI)
	urlParams.Set("prompt", "login")

	authURL.RawQuery = urlParams.Encode()
	return authURL.String(), nil
}

func (s *service) generateLogoutURL(idTokenHint, postLogoutRedirectURI string) (string, error) {
	logoutURL, err := url.Parse(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", s.keycloak.PublicHost, s.keycloak.Realm))
	if err != nil {
		return "", fmt.Errorf("url parse: %w", err)
	}

	if postLogoutRedirectURI == "" {
		postLogoutRedirectURI = s.keycloak.PostLogoutRedirectURI
	}

	urlParams := logoutURL.Query() // nosemgrep: rules.go.sql.go_sql_rule-concat-sqli - false positive, this is not SQL
	urlParams.Set("id_token_hint", idTokenHint)
	urlParams.Set("post_logout_redirect_uri", postLogoutRedirectURI)

	logoutURL.RawQuery = urlParams.Encode()
	return logoutURL.String(), nil
}

func (s *service) exchangeToken(ctx context.Context, code, redirectURI string) (*token, error) {
	if redirectURI == "" {
		redirectURI = s.keycloak.RedirectURI
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.keycloak.Host, s.keycloak.Realm)
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", s.keycloak.ClientID)
	data.Set("client_secret", s.keycloak.ClientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("http newRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpClient do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &httpStatusError{code: resp.StatusCode, message: body}
	}

	var token *token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	var user *user
	if token.AccessToken, user, err = s.enrichTokenWithAdditionalUserData(ctx, token.AccessToken); err != nil {
		return nil, fmt.Errorf("enrich token: %w", err)
	}

	if err := s.repo.UpdateLastLoginAt(ctx, sqlc.UpdateLastLoginAtParams{
		ID:     user.id,
		Source: user.source,
	}); err != nil {
		return nil, fmt.Errorf("update last login: %w", err)
	}

	return token, nil
}

func (s *service) refreshToken(ctx context.Context, refreshToken string) (*token, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.keycloak.Host, s.keycloak.Realm)
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", s.keycloak.ClientID)
	data.Set("client_secret", s.keycloak.ClientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("http newRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpClient do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &httpStatusError{code: resp.StatusCode, message: body}
	}

	var token *token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	if token.AccessToken, _, err = s.enrichTokenWithAdditionalUserData(ctx, token.AccessToken); err != nil {
		return nil, fmt.Errorf("enrich token: %w", err)
	}

	return token, nil
}

func (s *service) enrichTokenWithAdditionalUserData(ctx context.Context, accessToken string) (string, *user, error) {
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(accessToken, &claims, s.keyfunc); err != nil {
		return "", nil, fmt.Errorf("parse jwt: %w", err)
	}

	var id pgtype.UUID
	var user *user
	if zimbraID, ok := claims["zimbra_id"].(string); ok {
		if err := id.Scan(zimbraID); err == nil {
			if user, err = s.getUser(ctx, id, sourceZimbra); err != nil {
				return "", nil, fmt.Errorf("get user zimbra: %w", err)
			}
		}
	}

	// fallback using keycloak_id
	if user == nil {
		if keycloakID, err := claims.GetSubject(); err == nil {
			if err := id.Scan(keycloakID); err == nil {
				if user, err = s.getUser(ctx, id, sourceKeycloak); err != nil {
					return "", nil, fmt.Errorf("get user keycloak: %w", err)
				}
			}
		}
	}

	if user == nil {
		return "", nil, errUserNotFound
	}

	claims["nip"] = user.nip
	if len(user.roles) > 0 {
		claims["roles"] = user.roles
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.keycloak.KID

	signed, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", nil, fmt.Errorf("sign token: %w", err)
	}
	return signed, user, nil
}

func (s *service) getUser(ctx context.Context, id pgtype.UUID, source string) (*user, error) {
	nip, err := s.repo.GetUserNIPByIDAndSource(ctx, sqlc.GetUserNIPByIDAndSourceParams{
		ID:     id,
		Source: source,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo get user nip: %w", err)
	}

	rows, err := s.repo.ListUserRoleByNIP(ctx, nip)
	if err != nil {
		return nil, fmt.Errorf("repo list user role: %w", err)
	}

	roles := make(map[string]string, len(rows))
	for _, row := range rows {
		roles[row.Service.String] = row.Nama
	}

	return &user{
		id:     id,
		source: source,
		nip:    nip,
		roles:  roles,
	}, nil
}
