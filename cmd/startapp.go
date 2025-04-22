package main

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/products"
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/internal/receptions"
	"avito_pvz_test/internal/users"
	"avito_pvz_test/pkg/database"
	"avito_pvz_test/pkg/midware"
	"net/http"
)

// StartApp Старт приложения и инициализация зависимостей по слоям:
// ┌────────────────────────────────────────────┐
// │ Application Setup (StartApp)              │
// ├────────────────────────────────────────────┤
// │ 1. Загрузка конфигурации                   │
// │ 2. Инициализация подключения к БД          │
// │ 3. Создание репозиториев (data access)     │
// │ 4. Создание сервисов (бизнес-логика)       │
// │ 5. Инициализация хендлеров (transport)     │
// │ 6. Настройка маршрутов и запуск сервера    │
// └────────────────────────────────────────────┘
//
// Архитектура: config → db → repo → service → handler
func StartApp() {
	/*------------------------------------------------------------------------------*/

	/* 1) Подгружаем файл config */
	conf := config.NewConfig()

	/*------------------------------------------------------------------------------*/

	/* 2) работаем с базой данных */
	db := database.CreateDb(conf)

	/*------------------------------------------------------------------------------*/

	/* 3) создаем репозитории на основе бд */
	/* 3.1) репозиторий для User */
	userRepo := users.NewRepoUser(db)

	/* 3.2) репозиторий для PVZ */
	pvzRepo := pvz.NewRepoPVZ(db)

	/* 3.3) репозиторий для Reception */
	receptionRepo := receptions.NewRepoRecep(db)

	/* 3.4) репозиторий для Products */
	productRepo := products.NewRepoProduct(db)

	/*------------------------------------------------------------------------------*/

	/* 4) создаем сервисы на основе репозиториев */
	/* 4.1) сервис для User */
	userService := users.NewServUser(userRepo, conf)

	/* 4.2) сервис для PVZ */
	pvzService := pvz.NewServPvz(pvzRepo)

	/* 4.3) сервис для Reception */
	receptionService := receptions.NewServReception(receptionRepo)

	/* 4.4) сервис для Products */
	productsService := products.NewServProduct(productRepo, receptionRepo)

	/*------------------------------------------------------------------------------*/

	/* 5) создаем Handlers на основе сервисов */

	/* 5.1) handler для User */
	userHandler := users.NewHandUser(userService)

	/* 5.2) handler для PVZ */
	pvzHandler := pvz.NewHandPvz(pvzService)

	/* 5.3) handler для Reception */
	receptionHandler := receptions.NewReceptionHandler(receptionService)

	/* 5.4) handler для Products */
	productHandler := products.NewHandProduct(productsService)

	/*------------------------------------------------------------------------------*/
	/* 6) создаем router в который будем класть все Handler */

	router := http.NewServeMux()

	/*------------------------------------------------------------------------------*/

	/* 7) подключаем ручки к роутеру */

	/* 7.1) ручки для User */
	router.HandleFunc("POST /register", userHandler.CreateUser())
	router.HandleFunc("POST /login", userHandler.AuthenticateUser())
	router.HandleFunc("POST /dummyLogin", userHandler.GetTokenByRole())

	/* 7.2) ручки для PVZ */
	router.Handle("POST /pvz", midware.CheckRoleByToken(pvzHandler.CreatePVZ(), "moderator"))
	router.Handle("POST /pvz/{pvzId}/close_last_reception", midware.CheckRoleByToken(pvzHandler.CloseLastReceptionByPvz(), "client"))

	/* 7.3) ручки для Reception */
	router.Handle("POST /receptions", midware.CheckRoleByToken(receptionHandler.CreateReception(), "client"))

	/* 7.4) ручки для Products */
	router.Handle("POST /products", midware.CheckRoleByToken(productHandler.Create(), "client"))

	// 8. Передаем в server наши ручки
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 9. Запускаем server
	server.ListenAndServe()
}
