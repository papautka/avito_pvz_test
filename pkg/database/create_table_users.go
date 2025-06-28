package database

import (
	"fmt"
	"log"
)

func (db *Db) CreateTableUser() error {
	// Создаем таблицу users
	_, err := db.MyDb.Exec(`CREATE TABLE IF NOT EXISTS users (
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
