package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWT struct {
	SecretKey []byte
}

func NewJWT(secret string) *JWT {
	return &JWT{SecretKey: []byte(secret)}
}

func (j *JWT) Create(role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(j.SecretKey)
}

func (j *JWT) ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.SecretKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return "", fmt.Errorf("token expired")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("role missing")
	}

	return role, nil
}
