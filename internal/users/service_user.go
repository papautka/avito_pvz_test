package users

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/dto/payload"
	"avito_pvz_test/pkg/jwt"
	"fmt"
	"log"
)

type ServiceUser interface {
	Register(email, password, role string) (*User, error)
	Login(email, password string) (*payload.TokenResponse, error)
	GetToken(role string) (string, error)
}

type ServUser struct {
	UserRepo RepositoryUser // Интерфейс, а не конкретная реализация
	Config   *config.Config
}

func NewServUser(repo RepositoryUser, conf *config.Config) ServiceUser {
	return &ServUser{
		UserRepo: repo,
		Config:   conf,
	}
}

// Register функция регистрации пользователя
func (service *ServUser) Register(email, password, role string) (*User, error) {
	user := NewUser(email, password, role)
	createdUser, err := service.UserRepo.CreateUser(user)
	if err != nil {
		log.Println("Register: не удалось создать пользователя")
		return nil, err
	}
	return createdUser, nil
}

// Login функция авторизации пользователя
func (service *ServUser) Login(email, password string) (*payload.TokenResponse, error) {
	user, err := service.UserRepo.FindUserByEmailPass(email, password)
	if err != nil {
		log.Println("Login", err)
		return nil, err
	}
	tokenStr, err := service.GetToken(user.Role)
	if err != nil {
		log.Println("Login", err)
		return nil, err
	}
	req := &payload.TokenResponse{
		Token: tokenStr,
	}
	return req, nil
}

// GetToken достать токен соответсвующее роли
func (service *ServUser) GetToken(role string) (string, error) {
	var secret string
	switch role {
	case "client":
		secret = service.Config.Auth.AuthTokenClient
	case "moderator":
		secret = service.Config.Auth.AuthTokenModerator
	}
	jwtToken := jwt.NewJWT(secret)
	tokenStr, err := jwtToken.Create(role)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenStr, nil
}
