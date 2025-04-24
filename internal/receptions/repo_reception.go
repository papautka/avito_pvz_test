package receptions

import (
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/pkg/database"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type RepositoryReception interface {
	Create(reception *Reception) (*Reception, error)
	ReturnLastReceptionOrEmpty(UUIDPVZ uuid.UUID) (*Reception, error)
	FindPVZById(UUIDpvz uuid.UUID) (*pvz.PVZ, error)
}

type RepoRecep struct {
	Database *database.Db
}

func NewRepoRecep(database *database.Db) RepositoryReception {
	return &RepoRecep{
		Database: database,
	}
}

// Create Создание Приемки (только для клиентов)
func (repo *RepoRecep) Create(reception *Reception) (*Reception, error) {
	query := `INSERT INTO receptions (date_time, pvzId, status) VALUES ($1, $2, $3) RETURNING id`
	err := repo.Database.MyDb.QueryRow(query, reception.DateTime, reception.PvzID, reception.Status).Scan(&reception.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании приемки: %w", err)
	}
	return reception, nil
}

// ReturnLastReceptionOrEmpty Функция, которая проверяет есть ли вообще приемки в таблице приемка с указанным pvzID
func (repo *RepoRecep) ReturnLastReceptionOrEmpty(UUIDPVZ uuid.UUID) (*Reception, error) {
	// 1. Создаем пустую приемку
	reception := NewReception(time.Now(), UUIDPVZ, "close")
	query := `SELECT id, date_time, pvzId, status FROM receptions WHERE pvzId = $1 ORDER BY date_time DESC LIMIT 1`
	result, err := repo.Database.MyDb.Query(query, UUIDPVZ)
	if err != nil {
		fmt.Println("У данного PVZ не было приемки следовательно мы вернули ошибку")
		return nil, err
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

// FindPVZById поиск PVZ по id (для ручки /receptions - создания новой приемки)
func (repo *RepoRecep) FindPVZById(UUIDpvz uuid.UUID) (*pvz.PVZ, error) {
	pvz := &pvz.PVZ{}
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
