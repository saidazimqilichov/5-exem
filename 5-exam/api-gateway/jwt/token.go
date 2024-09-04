package jwt

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("said005")

type TokenClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func VerifyToken(bearerToken string) (bool, string, error) {
	tokenString := strings.TrimPrefix(bearerToken, "Bearer ")
	if tokenString == "" {
		return false, "", nil
	}

	
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false, "", err
	}

	
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return false, "", nil
	}


	if time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
		return false, "", nil
	}

	return true, claims.Email, nil
}
