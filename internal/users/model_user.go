package users

import (
	"github.com/google/uuid"
	"math/rand/v2"
)

type User struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	Password string    `json:"password"`
}

func NewUser(email, password, role string) *User {
	if password == "" {
		return &User{
			Email:    email,
			Password: generateRandomPassword(8),
			Role:     role,
		}
	}
	return &User{
		Email:    email,
		Password: password,
		Role:     role,
	}
}

func generateRandomPassword(n int) string {
	b := make([]byte, n)
	for i := 0; i < len(b)-1; i++ {
		b[i] = byte(rangeOfRandom(33, 126))
	}
	return string(b)
}

func rangeOfRandom(min, max int) int {
	return rand.IntN(max-min+1) + min
}
