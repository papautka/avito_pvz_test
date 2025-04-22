package pvz

import (
	"avito_pvz_test/pkg/database"
	"fmt"
	"github.com/google/uuid"
)

type RepositoryPvz interface {
	Create(pvz *PVZ) (*PVZ, error)
	FindPVZById(UUIDpvz uuid.UUID) (*PVZ, error)
	UpdateStatus(UUIDpvz uuid.UUID) (*ReceptionForPvz, error)
}

type RepoPVZ struct {
	Database *database.Db
}

func NewRepoPVZ(database *database.Db) RepositoryPvz {
	return &RepoPVZ{
		Database: database,
	}
}

// Create создание PVZ только для модераторов
func (repo *RepoPVZ) Create(pvz *PVZ) (*PVZ, error) {
	query := `INSERT INTO pvz (id, registration_date, city) VALUES ($1, $2, $3) RETURNING id`
	err := repo.Database.MyDb.QueryRow(query, pvz.ID, pvz.RegistrationDate, pvz.City).Scan(&pvz.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return pvz, nil
}

// FindPVZById поиск PVZ по id (для ручки /receptions - создания новой приемки)
func (repo *RepoPVZ) FindPVZById(UUIDpvz uuid.UUID) (*PVZ, error) {
	pvz := &PVZ{}
	query := `SELECT id, registration_date, city FROM pvz WHERE id=$1`
	result, err := repo.Database.MyDb.Query(query, UUIDpvz)
	if err != nil {
		return nil, fmt.Errorf("нет такого значения UUID в базе данных: %w", err)
	}
	defer result.Close()

	if result.Next() {
		err = result.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения данных из результата: %w", err)
		}
		return pvz, nil
	}
	return nil, fmt.Errorf("пункт выдачи с UUID %s не найден", UUIDpvz)
}

func (repo *RepoPVZ) UpdateStatus(UUIDpvz uuid.UUID) (*ReceptionForPvz, error) {
	query := `
		UPDATE receptions 
		SET status = 'close'
		WHERE id = (
			SELECT id FROM receptions
			WHERE pvzId = $1
				AND status != 'close'
			ORDER BY date_time DESC
			LIMIT 1
		)
		RETURNING id, date_time, pvzId, status
	`

	rows, err := repo.Database.MyDb.Query(query, UUIDpvz)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		receptForPvz := &ReceptionForPvz{}
		err = rows.Scan(&receptForPvz.ID, &receptForPvz.DateTime, &receptForPvz.PvzID, &receptForPvz.Status)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения результата: %w", err)
		}
		return receptForPvz, nil
	}

	// Если rows.Next() не сработал, значит, RETURNING ничего не вернуло (не было подходящей приемки)
	return nil, fmt.Errorf("нет приемок, которые нужно закрывать")
}
