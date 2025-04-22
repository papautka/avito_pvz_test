package database

import (
	"fmt"
	"log"
)

func (db *Db) CreateTableProducts() error {
	// вкл расширение pgcrypto для генерации UUID
	_, err := db.MyDb.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		log.Fatalf("Ошибка при создании расширения pgcrypto: %v", err)
		return err
	}
	// Создаем enum тип для type, если он ещё не создан
	_, err = db.MyDb.Exec(`DO $$
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'type_enum') THEN
		   CREATE TYPE type_enum AS ENUM ('электроника','одежда','обувь');
		END IF;
	END$$;	
	`)
	if err != nil {
		log.Fatalf("Ошибка при создании типа ENUM: %v", err)
		return err
	}
	// Создаем таблицу products
	_, err = db.MyDb.Exec(`
	CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		dateTime TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		type type_enum NOT NULL,
	    receptionId UUID NOT NULL,
	    CONSTRAINT fk_products FOREIGN KEY (receptionId) REFERENCES receptions(id)
	);
	`)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы products: %v", err)
		return err
	}

	fmt.Println("✅ Таблица products успешно создана.")
	return nil
}
