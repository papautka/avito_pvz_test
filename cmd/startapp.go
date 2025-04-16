package main

import (
	"avito_pvz_test/config"
	"fmt"
)

func StartApp() {
	/* Подгружаем файл config */
	conf := config.NewConfig()
	fmt.Println(conf)
}
