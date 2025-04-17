package main

import (
	"avito_pvz_test/config"
	"avito_pvz_test/pkg/database"
	"fmt"
)

func StartApp() {
	/* Подгружаем файл config */
	conf := config.NewConfig()
	fmt.Println(conf)

	/* Подключаемся к базе данных */
	db := database.NewDb(conf)
	fmt.Println(db)
}
