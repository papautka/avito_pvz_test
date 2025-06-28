package database

import (
	"avito_pvz_test/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
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
	// настройка pull соеденений
	mainDb.SetMaxOpenConns(100)                 // максимум 100 одновременных соединений
	mainDb.SetMaxIdleConns(20)                  // максимум 20 "простаивающих" соединений
	mainDb.SetConnMaxLifetime(time.Minute * 10) // соединения живут до 10 минут

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

	/* 2.0) Применяем миграции */
	migrator := NewMigrator(conf)
	defer migrator.Close()

	if err := migrator.MigrateUp(); err != nil {
		log.Fatal("Не удалось применить миграции:", err)
		return nil
	}

	return db
}

// Close закрывает соединение с базой данных
func (db *Db) Close() error {
	if db.MyDb != nil {
		return db.MyDb.Close()
	}
	return nil
}

// InitExtensions инициализирует необходимые расширения PostgreSQL
func (db *Db) InitExtensions() error {
	// Включаем расширение pgcrypto для генерации UUID
	_, err := db.MyDb.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		return fmt.Errorf("ошибка при создании расширения pgcrypto: %v", err)
	}

	// Создаем enum тип для роли, если он еще не создан
	_, err = db.MyDb.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_enum') THEN
			CREATE TYPE role_enum AS ENUM ('client', 'moderator');
		END IF;
	END$$;`)
	if err != nil {
		return fmt.Errorf("ошибка при создании типа ENUM role_enum: %v", err)
	}

	// Создаем enum тип для статуса приемки
	_, err = db.MyDb.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_enum') THEN
			CREATE TYPE status_enum AS ENUM ('in_progress', 'close');
		END IF;
	END$$;`)
	if err != nil {
		return fmt.Errorf("ошибка при создании типа ENUM status_enum: %v", err)
	}

	// Создаем enum тип для типа товара
	_, err = db.MyDb.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'product_type_enum') THEN
			CREATE TYPE product_type_enum AS ENUM ('электроника', 'одежда', 'обувь');
		END IF;
	END$$;`)
	if err != nil {
		return fmt.Errorf("ошибка при создании типа ENUM product_type_enum: %v", err)
	}

	// Создаем enum тип для города
	_, err = db.MyDb.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'city_enum') THEN
			CREATE TYPE city_enum AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');
		END IF;
	END$$;`)
	if err != nil {
		return fmt.Errorf("ошибка при создании типа ENUM city_enum: %v", err)
	}

	fmt.Println("✅ Расширения и типы PostgreSQL успешно инициализированы")
	return nil
}
