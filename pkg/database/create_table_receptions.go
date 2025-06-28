package database

import (
	"fmt"
	"log"
)

func (db *Db) CreateTableReception() error {
	// Создаем таблицу reception
	_, err := db.MyDb.Exec(`CREATE TABLE IF NOT EXISTS reception (
		id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		pvz_id UUID REFERENCES pvz(id) ON DELETE CASCADE,
		status status_enum DEFAULT 'in_progress'
	);`)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
		return err
	}

	fmt.Println("✅ Таблица 'reception' успешно создана")
	return nil
}
