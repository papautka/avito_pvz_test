package users

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/dto/payload"
	"avito_pvz_test/pkg/jwt"
	"context"
	"fmt"
	"log"
)

type UsersService struct {
	UserRepo UserRepository // Интерфейс, а не конкретная реализация
	Config   *config.Config
}

func NewUserService(repo UserRepository, conf *config.Config) *UsersService {
	return &UsersService{
		UserRepo: repo,
		Config:   conf,
	}
}

// фунцкия регистарции пользователя
func (service *UsersService) Register(ctx context.Context, email, password, role string) (*User, error) {
	user := NewUser(email, password, role)
	createdUser, err := service.UserRepo.CreateUser(ctx, user)
	if err != nil {
		log.Println("Register: не удалось создать пользователя")
		return nil, err
	}
	return createdUser, nil
}

// функция авторизации пользователя
func (service *UsersService) Login(ctx context.Context, email, password string) (*payload.TokenResponse, error) {
	user, err := service.UserRepo.FindUserByEmailPass(ctx, email, password)
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
func (service *UsersService) GetToken(role string) (string, error) {
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
