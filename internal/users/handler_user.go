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

func NewUserHandler(router *http.ServeMux, userHandlerDep *UserHandlerDependency) *UserHandler {
	userHandler := &UserHandler{
		userHandlerDep.UserService,
		userHandlerDep.Config,
	}
	router.HandleFunc("POST /register", userHandler.CreateUser())
	return userHandler
}
