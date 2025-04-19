package users

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/dto/payload"
	"avito_pvz_test/pkg/req"
	"log"
	"net/http"
)

type UserHandlerDependency struct {
	*UserService
	*config.Config
}

type UserHandler struct {
	*UserService
	*config.Config
}

func (userHandler *UserHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.UserCreateRequest](&w, r)
		if err != nil {
			log.Println("CreateUser: функция HandleBody вернула nil", err)
			return
		}
		createdUser, err := userHandler.Register(body.Email, body.Password, body.Role)
		if err != nil {
			log.Println("CreateUser: не удалось создать пользователя")
			return
		}
		req.JsonResponse(&w, &createdUser)
	}
}

func (userHandler *UserHandler) AuthenticateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.UserAuthRequest](&w, r)
		if err != nil {
			log.Println("AuthenticateUser: функция HandleBody вернула nil", err)
			return
		}
		jwtPoint, err := userHandler.Login(body.Email, body.Password)
		if err != nil {
			log.Println("AuthenticateUser", err)
			return
		}
		req.JsonResponse(&w, &jwtPoint)
	}
}

func (userHandler *UserHandler) GetTokenByRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.TokenRequestRole](&w, r)
		if err != nil {
			log.Println("GetTokenByRole функция HandleBody вернула nil", err)
			return
		}
		tokenStr, err := userHandler.GetToken((*body).Role)
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

func NewUserHandler(router *http.ServeMux, userHandlerDep *UserHandlerDependency) *UserHandler {
	userHandler := &UserHandler{
		userHandlerDep.UserService,
		userHandlerDep.Config,
	}
	router.HandleFunc("POST /register", userHandler.CreateUser())
	router.HandleFunc("POST /login", userHandler.AuthenticateUser())
	router.HandleFunc("POST /dummyLogin", userHandler.GetTokenByRole())

	return userHandler
}
