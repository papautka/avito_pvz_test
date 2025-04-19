package pvz

import (
	"avito_pvz_test/pkg/database"
	"fmt"
	"log"
)

type PVZRepoDeps struct {
	Database *database.Db
}

type PVZRepo struct {
	Database *database.Db
}

func NewPVZRepo(database *database.Db) *PVZRepo {
	return &PVZRepo{
		Database: database,
	}
}

// создание PVZ только для модераторов
func (repo *PVZRepo) Create(pvz *PVZ) (*PVZ, error) {
	query := `INSERT INTO pvz (id, registration_date, city) VALUES ($1, $2, $3) RETURNING id`
	err := repo.Database.MyDb.QueryRow(query, pvz.ID, pvz.RegistrationDate, pvz.City).Scan(&pvz.ID)
	if err != nil {
		log.Println("Create ", err)
		return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return pvz, nil
}
