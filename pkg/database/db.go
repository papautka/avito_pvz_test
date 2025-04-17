package database

import (
	"avito_pvz_test/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Db struct {
	MyDb *gorm.DB
}

func NewDb(conf *config.Config) *Db {
	// открываем соеденение
	mainDb, err := gorm.Open(postgres.Open(conf.Db.DsnDb), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось открыть соеденение к базе данных")
		return nil
	}
	return &Db{
		MyDb: mainDb,
	}
}
