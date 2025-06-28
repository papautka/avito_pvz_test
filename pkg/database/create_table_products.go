package database

import (
	"fmt"
	"log"
)

func (db *Db) CreateTableProducts() error {
	// Создаем таблицу products
	_, err := db.MyDb.Exec(`CREATE TABLE IF NOT EXISTS products (
		id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		type product_type_enum NOT NULL,
		reception_id UUID REFERENCES reception(id) ON DELETE CASCADE
	);`)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
		return err
	}

	fmt.Println("✅ Таблица 'products' успешно создана")
	return nil
}
