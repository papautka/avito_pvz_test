package database

import (
	"avito_pvz_test/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Db struct {
	MyDb *sql.DB
}

func NewDb(conf *config.Config) *Db {
	// открываем соеденение
	mainDb, err := sql.Open("postgres", conf.Db.DsnDb)
	if err != nil {
		log.Fatal("Ошибка при открытии соединения с БД: %v", err)
		return nil
	}
	// Проверим подключение
	if err = mainDb.Ping(); err != nil {
		log.Fatalf("Ошибка при пинге БД: %v", err)
	}
	fmt.Println("✅ Успешное подключение к базе данных")
	return &Db{
		MyDb: mainDb,
	}
}
