package database

import (
	"fmt"
	"log"
)

func (db *Db) CreateTableReception() error {
	// Включаем расширение pgcrypto для генерации UUID
	_, err := db.MyDb.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		log.Fatalf("Ошибка при создании расширения pgcrypto: %v", err)
		return err
	}

	// Создаем enum тип для status, если он еще не создан
	_, err = db.MyDb.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_enum') THEN
			CREATE TYPE status_enum AS ENUM ('in_progress', 'close');
		END IF;
	END$$;`)
	if err != nil {
		log.Fatalf("Ошибка при создании типа ENUM: %v", err)
		return err
	}

	// Создаем таблицу receptions
	_, err = db.MyDb.Exec(`CREATE TABLE IF NOT EXISTS receptions (
		id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		date_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		pvzId UUID NOT NULL,
		status status_enum NOT NULL,
		CONSTRAINT fk_pvz FOREIGN KEY (pvzId) REFERENCES pvz(id)
	);`)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы receptions: %v", err)
		return err
	}

	// 4. Уникальный индекс на одну активную приёмку на pvzId
	_, err = db.MyDb.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS only_one_active_reception_per_pvz
		ON receptions(pvzId)
		WHERE status = 'in_progress';
	`)
	if err != nil {
		log.Fatalf("Ошибка при создании уникального индекса: %v", err)
		return err
	}

	fmt.Println("✅ Таблица 'receptions' успешно создана")
	return nil
}
