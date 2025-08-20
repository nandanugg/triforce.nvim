package apitest

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v4"
)

var (
	// JWTPrivateKey pairs with JWTPublicKey, used for signing HTTP
	// Authorization header in tests.
	JWTPrivateKey *rsa.PrivateKey

	// JWTPublicKey pairs with JWTPrivateKey, used for verifying HTTP
	// Authorization header in tests.
	JWTPublicKey *rsa.PublicKey
)

func init() {
	JWTPrivateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	JWTPublicKey = &JWTPrivateKey.PublicKey
}

// GenerateAuthHeader generates HTTP Authorization header for use in tests.
func GenerateAuthHeader(userID int64, role ...string) string {
	claims := jwt.MapClaims{"user_id": userID}
	if len(role) > 0 {
		claims["role"] = role[0]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, _ := token.SignedString(JWTPrivateKey)
	return "Bearer " + tokenString
}
