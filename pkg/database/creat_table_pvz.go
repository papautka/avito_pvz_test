package database

import (
	"fmt"
	"log"
)

func (db *Db) CreateTablePVZ() error {
	// вкл расширение pgcrypto для генерации UUID
	_, err := db.MyDb.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		log.Fatalf("Ошибка при создании расширения pgcrypto: %v", err)
		return err
	}
	// Создаем enum тип для city, если он ещё не создан
	_, err = db.MyDb.Exec(`DO $$
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'city_enum') THEN
		   CREATE TYPE city_enum AS ENUM ('Москва','Санкт-Петербург','Казань');
		END IF;
	END$$;	
	`)
	if err != nil {
		log.Fatalf("Ошибка при создании типа ENUM: %v", err)
		return err
	}
	// Создаем таблицу pvz
	_, err = db.MyDb.Exec(`
	CREATE TABLE IF NOT EXISTS pvz (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		registration_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		city city_enum NOT NULL
	);
	`)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы pvz: %v", err)
		return err
	}

	fmt.Println("✅ Таблица pvz успешно создана.")
	return nil
}
