package pvz

import (
	"avito_pvz_test/pkg/database"
	"avito_pvz_test/pkg/req"
	"fmt"
	"github.com/google/uuid"
)

type RepositoryPvz interface {
	Create(pvz *PVZ) (*PVZ, error)
	FindPVZById(UUIDpvz uuid.UUID) (*PVZ, error)
	UpdateStatus(UUIDpvz uuid.UUID) (*ReceptionForPvz, error)
	GetPVZPageAndLimit(filter *req.FilterWithPagination) (*PvzListResponse, error)
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

func (repo *RepoPVZ) GetPVZPageAndLimit(filter *req.FilterWithPagination) (*PvzListResponse, error) {
	query := `SELECT id, registration_date, city FROM pvz
				WHERE registration_date BETWEEN $1 AND $2
				ORDER BY registration_date desc
				LIMIT $3 OFFSET $4`
	rows, err := repo.Database.MyDb.Query(query, filter.StartDate, filter.EndDate, filter.Limit, filter.Offset)
	if err != nil {
		fmt.Println("GetPVZPageAndLimit 1 ", err)
		return nil, err
	}
	defer rows.Close()
	var pvzListResponse PvzListResponse
	for rows.Next() {
		curPvz := PvzResponse{}
		err = rows.Scan(&curPvz.Pvz.ID, &curPvz.Pvz.RegistrationDate, &curPvz.Pvz.City)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения результата: %w", err)
		}
		curPvz.ArrayReception, err = repo.findReceptionPvzId(&curPvz.Pvz.ID)
		if err != nil {
			return nil, err
		}
		pvzListResponse.ArrayPvzResponse = append(pvzListResponse.ArrayPvzResponse, curPvz)
	}
	return &pvzListResponse, nil
}

func (repo *RepoPVZ) findReceptionPvzId(curPvzId *uuid.UUID) ([]ReceptionResponse, error) {
	query := `SELECT id, date_time, pvzId, status FROM receptions WHERE pvzId = $1`
	rows, err := repo.Database.MyDb.Query(query, curPvzId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var slRecep []ReceptionResponse
	for rows.Next() {
		curReception := ReceptionResponse{}
		err = rows.Scan(&curReception.Reception.ID, &curReception.Reception.DateTime, &curReception.Reception.PvzID, &curReception.Reception.Status)
		if err != nil {
			return nil, err
		}
		curReception.ArrayProduct, err = repo.findAllProductReceptionId(&curReception.Reception.ID)
		if err != nil {
			return nil, err
		}
		slRecep = append(slRecep, curReception)
	}
	return slRecep, nil
}

func (repo *RepoPVZ) findAllProductReceptionId(curReceptionId *uuid.UUID) ([]Product, error) {
	query := `SELECT id, datetime, type, receptionId FROM products WHERE receptionId = $1`
	rows, err := repo.Database.MyDb.Query(query, curReceptionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sliceProduct []Product
	for rows.Next() {
		product := Product{}
		err = rows.Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionId)
		if err != nil {
			return nil, err
		}
		sliceProduct = append(sliceProduct, product)
	}
	return sliceProduct, nil
}
