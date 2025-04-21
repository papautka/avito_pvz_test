package receptions

import (
	"avito_pvz_test/config"
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
	pvzObj, err := service.ReceptionRepository.FindPVZById(uuidValPvz)
	if err != nil {
		return nil, fmt.Errorf("Нет pvzId в базе данных pvz : %w", err)
	}
	// 5. в таблице Приемка найти последнюю приемку привязанную к PVZ:
	// 5.1. Если её нет создать её со статусом close
	reception, err := service.ReceptionRepository.ReturnLastReceptionOrEmpty(pvzObj.ID)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	// 5.2. Если она есть то проверить статус у последней приемки
	if reception.Status == "in_progress" {
		return nil, fmt.Errorf("в данном PVZ уже есть незакрытая приемка")
	}
	(*reception).Status = "in_progress"
	(*reception).DateTime = time.Now()
	createReception, err := service.ReceptionRepository.Create(reception)
	if err != nil {
		return nil, err
	}
	return createReception, nil
}
