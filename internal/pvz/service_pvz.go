package pvz

import (
	"avito_pvz_test/config"
	"avito_pvz_test/pkg/req"
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
