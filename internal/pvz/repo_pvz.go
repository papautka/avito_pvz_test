package pvz

import (
	"avito_pvz_test/pkg/database"
	"avito_pvz_test/pkg/req"
	"fmt"
	"github.com/google/uuid"
	"strings"
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
		return nil, err
	}
	defer rows.Close()
	var pvzListResponse PvzListResponse
	var pvzListId []uuid.UUID
	for rows.Next() {
		curPvz := PvzResponse{}
		err = rows.Scan(&curPvz.Pvz.ID, &curPvz.Pvz.RegistrationDate, &curPvz.Pvz.City)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения результата: %w", err)
		}
		if err != nil {
			return nil, err
		}
		pvzListId = append(pvzListId, curPvz.Pvz.ID)
		pvzListResponse.ArrayPvzResponse = append(pvzListResponse.ArrayPvzResponse, curPvz)
	}
	// присваиваем Reception соответствующим PVZ
	receptionMap, err := repo.findReceptionsByPvzIDs(pvzListId)
	if err != nil {
		return nil, err
	}
	for index, elemPvzResponse := range pvzListResponse.ArrayPvzResponse {
		pvzId := elemPvzResponse.Pvz.ID
		if receptions, ok := receptionMap[pvzId]; ok {
			pvzListResponse.ArrayPvzResponse[index].ArrayReception = receptions
		}
	}

	return &pvzListResponse, nil
}

func (repo *RepoPVZ) findReceptionsByPvzIDs(pvzIDs []uuid.UUID) (map[uuid.UUID][]ReceptionResponse, error) {
	if len(pvzIDs) == 0 {
		return nil, nil
	}
	// генерируем $1, $2, $3, ..., $n
	placeholders := make([]string, len(pvzIDs))
	args := make([]interface{}, len(pvzIDs))
	for i, id := range pvzIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	// более оптимально пишем Query запрос
	query := fmt.Sprintf(`SELECT id, date_time, pvzId, status FROM receptions WHERE pvzId IN (%s)`,
		strings.Join(placeholders, ","))

	rows, err := repo.Database.MyDb.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultRecResponse := make(map[uuid.UUID][]ReceptionResponse, len(pvzIDs))
	resultReceptionId := make([]uuid.UUID, 0, len(pvzIDs))
	for rows.Next() {
		var r ReceptionResponse
		err = rows.Scan(&r.Reception.ID, &r.Reception.DateTime, &r.Reception.PvzID, &r.Reception.Status)
		if err != nil {
			return nil, err
		}
		resultRecResponse[r.Reception.PvzID] = append(resultRecResponse[r.Reception.PvzID], r)
		resultReceptionId = append(resultReceptionId, r.Reception.ID)
	}
	// Получаем продукты по receptionId
	productValue, err := repo.findAllProductReceptionId(resultReceptionId)
	if err != nil {
		return nil, err
	}
	// присваиваем продукты соответствующим Reception
	for _, elemReception := range resultRecResponse {
		for i := range elemReception {
			receptionsId := elemReception[i].Reception.ID
			elemReception[i].ArrayProduct = productValue[receptionsId]
		}
	}
	return resultRecResponse, nil
}

func (repo *RepoPVZ) findAllProductReceptionId(recIds []uuid.UUID) (map[uuid.UUID]([]Product), error) {
	if len(recIds) == 0 {
		return nil, nil
	}
	// генерируем $1, $2, ..., $n
	placeholders := make([]string, len(recIds))
	args := make([]interface{}, len(recIds))
	for i, id := range recIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT id, datetime, type, receptionId FROM products WHERE receptionid IN (%s)`, strings.Join(placeholders, ","))
	rows, err := repo.Database.MyDb.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	resultRecResponse := make(map[uuid.UUID][]Product, len(recIds))
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.DateTime, &p.Type, &p.ReceptionId)
		if err != nil {
			return nil, err
		}
		resultRecResponse[p.ReceptionId] = append(resultRecResponse[p.ReceptionId], p)
	}
	return resultRecResponse, nil
}
