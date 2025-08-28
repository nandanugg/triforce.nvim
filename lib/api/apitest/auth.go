package apitest

import (
	"crypto/rand"
	"crypto/rsa"
	"strconv"

	"github.com/golang-jwt/jwt/v5"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

var (
	jwtPrivateKey *rsa.PrivateKey

	// Keyfunc encapsulate jwt.Keyfunc, used for verifying HTTP
	// Authorization header in tests.
	Keyfunc *api.Keyfunc
)

func init() {
	jwtPrivateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	Keyfunc = &api.Keyfunc{
		Keyfunc:  func(*jwt.Token) (any, error) { return &jwtPrivateKey.PublicKey, nil },
		Audience: "testing",
	}
}

// GenerateAuthHeader generates HTTP Authorization header for use in tests.
func GenerateAuthHeader(userID int64, role ...string) string {
	h := strconv.FormatInt(userID, 10)
	if len(role) > 0 {
		h += ":" + role[0]
	}
	return h

	// claims := jwt.MapClaims{"user_id": userID, "aud": "testing"}
	// if len(role) > 0 {
	// 	claims["role"] = role[0]
	// }
	//
	// token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// tokenString, _ := token.SignedString(jwtPrivateKey)
	// return "Bearer " + tokenString
}
