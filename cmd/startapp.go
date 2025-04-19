package main

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/internal/server"
	"avito_pvz_test/internal/users"
	"avito_pvz_test/pkg/database"
	"avito_pvz_test/pkg/repos"
	"log"
)

func StartApp() {
	/* 1) Подгружаем файл config */
	conf := config.NewConfig()

	/* 2) работаем с базой данных */
	db := CreateDb(conf)

	/* 3) создаем репозитории на основе бд */
	allRepos := CreateRepository(db)

	/*4) Запускаем сервер передавая туда репозиторий */
	server.ServerStart(conf, allRepos)
}

func CreateDb(conf *config.Config) *database.Db {
	/* 2) Подключаемся к базе данных */
	db := database.NewDb(conf)

	/* 2.1) Создаем таблицу в бд для user если она не создана */
	err := db.CreateTableUser()
	if err != nil {
		log.Fatal("Не удалось создать таблицу users:", err)
		return nil
	}

	/* 2.2) Создаем таблицу в бд для PVZ если она не создана */
	err = db.CreateTablePVZ()
	if err != nil {
		log.Fatal("Не удалось создать таблицу users:", err)
		return nil
	}
	return db
}

func CreateRepository(db *database.Db) *repos.AllRepository {
	/* 3.1) репозиторий для User */
	userRepository := users.NewUserRepo(db)

	/* 3.2) репозиторий для PVZ */
	pvzRepository := pvz.NewPVZRepo(db)

	return repos.NewAllRepository(userRepository, pvzRepository)
}
