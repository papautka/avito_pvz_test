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

func (db *Db) CreateTableUser() error {
	// Включаем расширение pgcrypto для генерации UUID
	_, err := db.MyDb.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		log.Fatalf("Ошибка при создании расширения pgcrypto: %v", err)
		return err
	}

	// Создаем enum тип для роли, если он еще не создан
	_, err = db.MyDb.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_enum') THEN
			CREATE TYPE role_enum AS ENUM ('employee', 'moderator');
		END IF;
	END$$;
	`)
	if err != nil {
		log.Fatalf("Ошибка при создании типа ENUM: %v", err)
		return err
	}

	// Создаем таблицу users
	_, err = db.MyDb.Exec(`CREATE TABLE IF NOT EXISTS users (
		id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255),
		role role_enum NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
		return err
	}

	fmt.Println("✅ Таблица 'users' успешно создана")
	return nil
}
