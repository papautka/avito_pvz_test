package main

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/server"
	"avito_pvz_test/internal/users"
	"avito_pvz_test/pkg/database"
	"log"
)

func StartApp() {
	/* 1) Подгружаем файл config */
	conf := config.NewConfig()

	/* 2) Подключаемся к базе данных */
	db := database.NewDb(conf)

	/* 2.1) Создаем таблицу в бд если она не создана */
	err := db.CreateTableUser()
	if err != nil {
		log.Fatal("Не удалось создать таблицу users:", err)
		return
	}

	/* 3) репозиторий для User */
	userRepository := users.NewUserRepo(db)

	/*4) Запускаем сервер передавая туда репозиторий */
	server.ServerStart(conf, userRepository)
}
