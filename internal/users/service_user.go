package users

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/dto/payload"
	"avito_pvz_test/pkg/jwt"
	"fmt"
	"log"
)

type UserService struct {
	UserRepository *UserRepo
	Config         *config.Config
}

func NewUserService(repo *UserRepo, conf *config.Config) *UserService {
	return &UserService{
		UserRepository: repo,
		Config:         conf,
	}
}

// фунцкия регистарции пользователя
func (service *UserService) Register(email, password, role string) (*User, error) {
	// 1. TODO: реализовать валидацию email
	// TODO : если пользователь с таким email уже есть то ошибка
	user := NewUser(email, password, role)
	createdUser, err := service.UserRepository.Create(user)
	if err != nil {
		log.Println("Register: не удалось создать пользователя")
		return nil, err
	}
	return createdUser, nil
}

// функция авторизации пользователя
func (service *UserService) Login(email, password string) (*payload.TokenResponse, error) {
	user, err := service.UserRepository.FindUserByEmailPass(email, password)
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

// достать токен соответсвующей роли
func (service *UserService) GetToken(role string) (string, error) {
	if role != "client" && role != "moderator" {
		return "", nil
	}
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
