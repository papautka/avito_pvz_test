package server

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/internal/receptions"
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

	// 2.2 Подключаем к router ручки для Reception
	ConnectHandlerForReception(router, conf, reps.ReceptionRepo)

	// 5. Передаем server наши ручки
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 6. Запускаем server
	server.ListenAndServe()
}

func ConnectHandlerForReception(router *http.ServeMux, conf *config.Config, reps *receptions.ReceptionRepo) {
	// 2. Подключаем сторонние service
	receptionService := receptions.NewReceptionService(reps, conf)

	// 3. Подключаем receptionHandDependency для Reception
	receptionHandDepend := receptions.ReceptionHandlerDependency{
		receptionService,
		conf,
	}

	// 4. Подключаем ручки Reception к router
	receptions.NewReceptionHandler(router, &receptionHandDepend)
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

func ConnectHandlerForUser(router *http.ServeMux, conf *config.Config, repo *users.UserRepo) {
	// 1. Инициализация сервиса пользователя
	userService := users.NewUserService(repo, conf)

	// 2. Создание зависимостей хендлера
	handlerDeps := &users.UserHandlerDependency{
		UserService: userService,
		Config:      conf,
	}

	// 3. Подключение маршрутов пользователя
	users.NewUserHandler(router, handlerDeps)
}
