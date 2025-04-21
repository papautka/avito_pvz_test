package receptions

import (
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/pkg/database"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ReceptionRepoDeps struct {
	Database *database.Db
}

type ReceptionRepo struct {
	Database *database.Db
}

func NewReceptionRepo(database *database.Db) *ReceptionRepo {
	return &ReceptionRepo{
		Database: database,
	}
}

// Создание Приемки (только для клиентов)
func (repo *ReceptionRepo) Create(reception *Reception) (*Reception, error) {
	query := `INSERT INTO receptions (date_time, pvzId, status) VALUES ($1, $2, $3) RETURNING id`
	err := repo.Database.MyDb.QueryRow(query, reception.DateTime, reception.PvzID, reception.Status).Scan(&reception.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании приемки: %w", err)
	}
	return reception, nil
}

// Функция которая проверяет есть ли вообще приемки в таблице приемка с указанным pvz
func (repo *ReceptionRepo) ReturnLastReceptionOrEmpty(UUIDPVZ uuid.UUID) (*Reception, error) {
	// 1. Создаем пустую приемку
	reception := NewReception(time.Now(), UUIDPVZ, "in_progress")
	query := `SELECT id, date_time, pvzId, status FROM receptions WHERE pvzId = $1 ORDER BY date_time DESC LIMIT 1`
	result, err := repo.Database.MyDb.Query(query, UUIDPVZ)
	if err != nil {
		fmt.Println("У данного PVZ не было приемки следовательно мы вернули пустую")
		return reception, nil
	}
	defer result.Close()
	if result.Next() {
		err = result.Scan(&reception.ID, &reception.DateTime, &reception.PvzID, &reception.Status)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения данных из результата: %w", err)
		}
		return reception, nil
	}
	return reception, nil
}

// поиск PVZ по id (для ручки /receptions - создания новой приемки)
func (repo *ReceptionRepo) FindPVZById(UUIDpvz uuid.UUID) (*pvz.PVZ, error) {
	pvz := &pvz.PVZ{}
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
