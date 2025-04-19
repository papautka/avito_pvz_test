package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// JWT структура с приватным ключом
type JWT struct {
	Secret     string
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

// NewJWT инициализирует структуру JWT
func NewJWT(secret string) *JWT {
	// Генерация ключа на основе секретной строки (это можно адаптировать под ваш случай)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return nil
	}

	// Извлекаем публичный ключ из приватного
	publicKey := &privateKey.PublicKey

	return &JWT{
		Secret:     secret,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

// Create создает JWT токен для заданной роли с временем истечения 24 часа
func (j *JWT) Create(role string) (string, error) {
	// Текущее время
	now := time.Now()

	// Установка времени истечения через 24 часа
	expirationTime := now.Add(24 * time.Hour)

	// Создание токена с claims, включая exp
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"role": role,
		"exp":  expirationTime.Unix(),
	})

	// Подпись токена с использованием приватного ключа
	signedString, err := token.SignedString(j.PrivateKey)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

// ParseToken проверяет срок действия токена и роль
func (j *JWT) ParseToken(tokenString string) (string, error) {
	// 1) Декодирование и проверка подписи токена с использованием публичного ключа
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Мы должны проверить, что токен использует правильный алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Возвращаем публичный ключ для проверки подписи
		return j.PublicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("error parsing token: %v", err)
	}

	// 2) Проверка срока действия токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	experationTime, ok := claims["exp"].(float64)
	if !ok {
		return "", fmt.Errorf("token does not have an exp claim")
	}

	// Преобразуем exp в time.Time
	expirationTime := time.Unix(int64(experationTime), 0)
	if expirationTime.Before(time.Now()) {
		return "", fmt.Errorf("token has expired")
	}

	// Проверка роли
	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("token does not have a role claim")
	}

	// Возвращаем роль, если все проверки пройдены
	return role, nil
}
