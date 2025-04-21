package pvz

import (
	"avito_pvz_test/config"
	"avito_pvz_test/pkg/req"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type PvzService struct {
	PVZRepo *PVZRepo
	Config  *config.Config
}

func NewPvzService(repo *PVZRepo, config *config.Config) *PvzService {
	return &PvzService{
		PVZRepo: repo,
		Config:  config,
	}
}

func (pvz *PvzService) Register(id string, registrationDate string, city string) (*PVZ, error) {
	// проверяем UUID если не передан
	var uuidVal uuid.UUID
	var err error

	uuidVal, err = req.ParseUUIDOrGenerate(id)
	if err != nil {
		return nil, err
	}

	// обработка даты
	var regTime time.Time
	regTime, err = req.ParseTimeOrNow(registrationDate)
	if err != nil {
		return nil, err
	}

	newPvz := NewPVZ(uuidVal, regTime, city)
	createdPVZ, err := pvz.PVZRepo.Create(newPvz)
	if err != nil {
		log.Printf("Error creating PVZ: %v", err)
		return nil, err
	}
	return createdPVZ, nil
}

func (pvz *PvzService) ChangeStatusReceptionByPvzOnClose(id string) (*ReceptionForPvz, error) {
	// проверяем UUID если не передан
	var uuidVal uuid.UUID
	var err error

	uuidVal, err = req.ParseUUIDPvz(id)
	if err != nil {
		return nil, fmt.Errorf("Некорректное значение id")
	}
	pvzStruct, err := pvz.PVZRepo.FindPVZById(uuidVal)
	if err != nil {
		return nil, fmt.Errorf("Нет pvz c таким значением")
	}
	fmt.Println("ID нашего PVZ", pvzStruct.ID)
	recepForPvz, err := pvz.PVZRepo.UpdateStatus(pvzStruct.ID)
	fmt.Println("uuidReception", recepForPvz, "err", err)
	if err != nil {
		return nil, fmt.Errorf("У данного pvzId: ", pvzStruct.ID, " нет приемок или она закрыта. Код ошибки: ", err)
	}
	return recepForPvz, nil
}
