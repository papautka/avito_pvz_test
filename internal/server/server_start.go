package server

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/internal/users"
	"avito_pvz_test/pkg/repos"
	"net/http"
)

func ServerStart(conf *config.Config, reps *repos.AllRepository) {

	// 1. Создаем пустой router
	router := http.NewServeMux()

	// 2. Подключаем к router ручки для User
	ConnectHandlerForUser(router, conf, reps.UserRepo)

	// 2.1. Подключаем к router ручки для PVZ
	ConnectHandlerForPvz(router, conf, reps.PvzRepo)

	// 5. Передаем server наши ручки
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 6. Запускаем server
	server.ListenAndServe()
}

func ConnectHandlerForPvz(router *http.ServeMux, conf *config.Config, reps *pvz.PVZRepo) {
	// 2. Подключаем сторонние service
	pvzService := pvz.NewPvzService(reps, conf)

	// 3. Подключаем pvzHandDependency для PVZ
	pvzHandDepend := pvz.PvzHandlerDependency{
		pvzService,
		conf,
	}

	// 4. Подключаем ручки PVZ к router
	pvz.NewPvzHandler(router, &pvzHandDepend)
}

func ConnectHandlerForUser(router *http.ServeMux, conf *config.Config, reps *users.UserRepo) {
	// 2. Подключаем сторонние service
	userServcie := users.NewUserService(reps, conf)

	// 3. Подключаем userHandDependency для USER
	userHandDepend := users.UserHandlerDependency{
		userServcie,
		conf,
	}

	// 4. Подключаем ручки USERS к router
	users.NewUserHandler(router, &userHandDepend)
}
