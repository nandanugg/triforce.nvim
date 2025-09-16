package auth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/config"
)

type service struct {
	repo       *repository
	keycloak   config.Keycloak
	client     *http.Client
	privateKey *rsa.PrivateKey
	keyfunc    jwt.Keyfunc
}

func newService(repo *repository, keycloak config.Keycloak, client *http.Client, privateKey *rsa.PrivateKey, keyfunc jwt.Keyfunc) *service {
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

	query := authURL.Query()
	query.Set("client_id", s.keycloak.ClientID)
	query.Set("response_type", "code")
	query.Set("scope", "openid")
	query.Set("redirect_uri", redirectURI)
	query.Set("prompt", "login")

	authURL.RawQuery = query.Encode()
	return authURL.String(), nil
}

func (s *service) generateLogoutURL(idTokenHint, postLogoutRedirectURI string) (string, error) {
	logoutURL, err := url.Parse(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", s.keycloak.PublicHost, s.keycloak.Realm))
	if err != nil {
		return "", fmt.Errorf("url parse: %w", err)
	}

	query := logoutURL.Query()
	query.Set("id_token_hint", idTokenHint)
	query.Set("post_logout_redirect_uri", postLogoutRedirectURI)

	logoutURL.RawQuery = query.Encode()
	return logoutURL.String(), nil
}

func (s *service) exchangeToken(ctx context.Context, code, redirectURI string) (*token, error) {
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

	if token.AccessToken, err = s.enrichTokenWithAdditionalUserData(ctx, token.AccessToken); err != nil {
		return nil, fmt.Errorf("enrich token: %w", err)
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

	if token.AccessToken, err = s.enrichTokenWithAdditionalUserData(ctx, token.AccessToken); err != nil {
		return nil, fmt.Errorf("enrich token: %w", err)
	}

	return token, nil
}

func (s *service) enrichTokenWithAdditionalUserData(ctx context.Context, accessToken string) (string, error) {
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(accessToken, &claims, s.keyfunc); err != nil {
		return "", fmt.Errorf("parse jwt: %w", err)
	}

	var user *user
	if zimbraID, ok := claims["zimbra_id"].(string); ok {
		if id, err := uuid.Parse(zimbraID); err == nil {
			if user, err = s.repo.getUser(ctx, id, sourceZimbra); err != nil {
				return "", fmt.Errorf("get user zimbra: %w", err)
			}
		}
	}

	// fallback using keycloak_id
	if user == nil {
		if keycloakID, err := claims.GetSubject(); err == nil {
			if id, err := uuid.Parse(keycloakID); err == nil {
				if user, err = s.repo.getUser(ctx, id, sourceKeycloak); err != nil {
					return "", fmt.Errorf("get user keycloak: %w", err)
				}
			}
		}
	}

	if user == nil {
		return "", errUserNotFound
	}

	claims["nip"] = user.nip
	if len(user.roles) > 0 {
		claims["roles"] = user.roles
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.keycloak.KID

	signed, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}
