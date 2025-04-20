package receptions

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/pkg/req"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ReceptionService struct {
	ReceptionRepository *ReceptionRepo
	Config              *config.Config
}

func NewReceptionService(repo *ReceptionRepo, conf *config.Config) *ReceptionService {
	return &ReceptionService{
		ReceptionRepository: repo,
		Config:              conf,
	}
}

// функция создания приемки
func (service *ReceptionService) CreateReception(pvzId string) (*Reception, error) {

	// 1. Проверка на корректность данных Request
	var uuidValPvz uuid.UUID
	var err error
	uuidValPvz, err = req.ParseUUIDPvz(pvzId)
	if err != nil {
		return nil, fmt.Errorf("Некорректный формат pvzId : %w", err)
	}

	// 4. Убедиться что данный pvzId вообще есть в базе данных в таблице PVZ
	pvzObj, err := service.FindPVZById(uuidValPvz)
	if err != nil {
		return nil, fmt.Errorf("Нет pvzId в базе данных pvz : %w", err)
	}
	// 5. в таблице Приемка найти последнюю приемку привязанную к PVZ:
	// 5.1. Если её нет создать её со статусом close
	reception, err := service.ReturnLastReceptionOrEmpty(pvzObj.ID)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	// 5.2. Если она есть то проверить статус у последней приемки
	if reception.Status == "in_progress" {
		return nil, fmt.Errorf("в данном PVZ уже есть незакрытая приемка")
	}
	createReception, err := service.ReceptionRepository.Create(reception)
	if err != nil {
		return nil, err
	}
	return createReception, nil
}

// Функция которая проверяет есть ли вообще приемки в таблице приемка с указанным pvz
// если есть то вернуть последнюю приемку
//
//	если у последней приемки статус in_progress то вернуть ошибку
//
// если нет то создать новую приемку со статусом CLOSE
func (service *ReceptionService) ReturnLastReceptionOrEmpty(UUIDPVZ uuid.UUID) (*Reception, error) {
	// 1. Создаем пустую приемку
	reception := NewReception(time.Now(), UUIDPVZ, "close")
	query := `SELECT id, date_time, pvzId, status FROM receptions WHERE pvzId = $1 ORDER BY date_time DESC LIMIT 1`
	result, err := service.ReceptionRepository.Database.MyDb.Query(query, UUIDPVZ)
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
func (service *ReceptionService) FindPVZById(UUIDpvz uuid.UUID) (*pvz.PVZ, error) {
	pvz := &pvz.PVZ{}
	query := `SELECT id, registration_date, city FROM pvz WHERE id=$1`
	result, err := service.ReceptionRepository.Database.MyDb.Query(query, UUIDpvz)
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
