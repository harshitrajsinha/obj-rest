// Package models defines data structures and functions that are used across the application
package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims embeds jwt.RegisteredClaims and add user role for jwt payload
type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAuthToken creates a JWT token for authentication with user role as payload
func GenerateAuthToken(role string, secretKey string) (string, error) {

	expiration := time.Now().Add(2 * time.Minute).UTC()
	claims := &CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

// VerifyAuthToken validate and verify authenticity of token
func VerifyAuthToken(token string, authSecretKey string) (string, error) {

	var parsedClaims CustomClaims

	parsedToken, err := jwt.ParseWithClaims(token, &parsedClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(authSecretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if !parsedToken.Valid {
		return "", errors.New("token not valid")
	}

	return parsedClaims.Role, nil

}
