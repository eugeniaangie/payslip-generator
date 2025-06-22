package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtKey is the secret key used for signing JWT tokens.
var jwtKey []byte

// SetJWTKey allows setting a custom secret key for JWT generation.
func SetJWTKey(key []byte) {
	jwtKey = key
}

// GenerateJWT creates a signed JWT token containing user_id and username.
// Expiration is set to 24 hours.
func GenerateJWT(userID, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseJWT parses and validates a JWT token.
// Returns claims if valid, error otherwise.
func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
