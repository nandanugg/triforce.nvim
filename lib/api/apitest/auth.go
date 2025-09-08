package apitest

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

var (
	JwtPrivateKey *rsa.PrivateKey

	// Keyfunc encapsulate jwt.Keyfunc, used for verifying HTTP
	// Authorization header in tests.
	Keyfunc *api.Keyfunc
)

func init() {
	JwtPrivateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	Keyfunc = &api.Keyfunc{
		Keyfunc:  func(*jwt.Token) (any, error) { return &JwtPrivateKey.PublicKey, nil },
		Audience: "testing",
	}
}

func GenerateAuthHeader(service, nip string, role ...string) string {
	claims := jwt.MapClaims{"nip": nip, "aud": "testing"}
	if len(role) > 0 {
		claims["roles"] = map[string]string{
			service: role[0],
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, _ := token.SignedString(JwtPrivateKey)
	return "Bearer " + tokenString
}
