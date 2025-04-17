package users

import (
	"log"
)

type UserService struct {
	UserRepository *UserRepo
}

func NewUserService(repo *UserRepo) *UserService {
	return &UserService{UserRepository: repo}
}

// фунцкия регистарции пользователя
func (service *UserService) Register(email, password, role string) (*User, error) {
	// 1. TODO: реализовать валидацию email
	// TODO : если пользователь с таким email уже есть то ошибка
	user := NewUser(email, password, role)
	createdUser, err := service.UserRepository.Create(user)
	if err != nil {
		log.Println("CreateUser: не удалось создать пользователя123")
		return nil, err
	}
	return createdUser, nil
}
