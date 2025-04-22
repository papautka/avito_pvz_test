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

func CreateDb(conf *config.Config) *Db {
	/* 2) Подключаемся к базе данных */
	db := NewDb(conf)

	/* 2.1) Создаем таблицу в бд для user если она не создана */
	err := db.CreateTableUser()
	if err != nil {
		log.Fatal("Не удалось создать таблицу users:", err)
		return nil
	}
	/* 2.2) Создаем таблицу в бд для PVZ если она не создана */
	err = db.CreateTablePVZ()
	if err != nil {
		log.Fatal("Не удалось создать таблицу pvz:", err)
		return nil
	}
	/* 2.3) Создаем таблицу в бд для Reception(приемки) если она не создана*/
	err = db.CreateTableReception()
	if err != nil {
		log.Fatal("Не удалось создать таблицу reception:", err)
		return nil
	}

	/* 2.4) Создаем таблицу в бд для Products если она не создана */
	err = db.CreateTableProducts()
	if err != nil {
		log.Fatal("Не удалось создать таблицу products:", err)
		return nil
	}
	return db
}
