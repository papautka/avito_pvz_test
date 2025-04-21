package pvz

import (
	"avito_pvz_test/pkg/database"
	"fmt"
	"github.com/google/uuid"
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
		return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return pvz, nil
}

// поиск PVZ по id (для ручки /receptions - создания новой приемки)
func (repo *PVZRepo) FindPVZById(UUIDpvz uuid.UUID) (*PVZ, error) {
	pvz := &PVZ{}
	query := `SELECT id, registration_date, city FROM pvz WHERE id=$1`
	result, err := repo.Database.MyDb.Query(query, UUIDpvz)
	if err != nil {
		return nil, fmt.Errorf("Нет такого значения UUID в базе данных: %w", err)
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

func (repo *PVZRepo) UpdateStatus(UUIDpvz uuid.UUID) (*ReceptionForPvz, error) {
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
		fmt.Println("receptForPvz 11111:", receptForPvz, "err 11111", err)
		return receptForPvz, nil
	}

	// Если rows.Next() не сработал, значит, RETURNING ничего не вернуло (не было подходящей приемки)
	return nil, fmt.Errorf("нет приемок, которые нужно закрывать")
}
