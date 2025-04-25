package users

import (
	"avito_pvz_test/internal/dto/errorDto"
	"avito_pvz_test/internal/dto/payload"
	"avito_pvz_test/pkg/req"
	"fmt"
	"log"
	"net/http"
)

// HandlerUser Интерфейс handler
type HandlerUser interface {
	GetTokenByRole() http.HandlerFunc
	CreateUser() http.HandlerFunc
	AuthenticateUser() http.HandlerFunc
}

// HandUser Структура, реализующая интерфейс
type HandUser struct {
	userService ServiceUser
}

// NewHandUser Фабрика, возвращающая интерфейс
func NewHandUser(userService ServiceUser) HandlerUser {
	return &HandUser{
		userService,
	}
}

func (userHandler *HandUser) GetTokenByRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.TokenRequestRole](&w, r)
		if err != nil {
			log.Println("GetTokenByRole функция HandleBody вернула nil", err)
			return
		}
		tokenStr, err := userHandler.userService.GetToken((*body).Role)
		if err != nil {
			log.Println("GetTokenByRole", err)
			return
		}
		tokenFormatJson := &payload.TokenResponse{
			Token: tokenStr,
		}
		req.JsonResponse(&w, tokenFormatJson)
	}
}

func (userHandler *HandUser) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.UserCreateRequest](&w, r)
		if err != nil {
			log.Println("CreateUser: функция HandleBody вернула nil", err)
			return
		}
		createdUser, err := userHandler.userService.Register(body.Email, body.Password, body.Role)
		if err != nil {
			log.Println("CreateUser: не удалось создать пользователя")
			return
		}
		req.JsonResponse(&w, &createdUser)
	}
}

func (userHandler *HandUser) AuthenticateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.UserAuthRequest](&w, r)
		if err != nil {
			log.Println("AuthenticateUser: функция HandleBody вернула nil", err)
			msgError := fmt.Sprintf("%v", err)
			errorDto.ShowResponseError(&w, msgError, http.StatusBadRequest)
			return
		}
		jwtPoint, err := userHandler.userService.Login(body.Email, body.Password)
		if err != nil {
			log.Println("AuthenticateUser", err)
			msgError := fmt.Sprintf("%v", err)
			errorDto.ShowResponseError(&w, msgError, http.StatusBadRequest)
			return
		}
		req.JsonResponse(&w, &jwtPoint)
	}
}
