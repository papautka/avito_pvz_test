package database

import (
	"fmt"
	"log"
)

func (db *Db) CreateTablePVZ() error {
	// Создаем таблицу pvz
	_, err := db.MyDb.Exec(`CREATE TABLE IF NOT EXISTS pvz (
		id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		city city_enum NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
		return err
	}

	fmt.Println("✅ Таблица 'pvz' успешно создана")
	return nil
}
