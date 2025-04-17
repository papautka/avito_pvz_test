package server

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/users"
	"net/http"
)

func ServerStart(conf *config.Config, reps *users.UserRepo) {

	// 1. Создаем пустой router
	router := http.NewServeMux()

	// 2. Подключаем сторонние service
	userServcie := users.NewUserService(reps)

	// 3. Подключаем userHandDependency для USER
	userHandDepend := users.UserHandlerDependency{
		userServcie,
		conf,
	}

	// 4. Подключаем ручки USERS к router
	users.NewUserHandler(router, &userHandDepend)

	// 5. Передаем server наши ручки
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 6. Запускаем server
	server.ListenAndServe()
}
