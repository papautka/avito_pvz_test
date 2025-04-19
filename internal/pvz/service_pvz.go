package pvz

import (
	"avito_pvz_test/config"
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

	uuidVal, err = parseUUIDOrGenerate(id)
	if err != nil {
		return nil, err
	}

	// обработка даты
	var regTime time.Time
	regTime, err = parseTimeOrNow(registrationDate)
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

func parseUUIDOrGenerate(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.New(), nil
	}
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("неверный формат UUID: %w", err)
	}
	return parsed, nil
}

func parseTimeOrNow(value string) (time.Time, error) {
	if value == "" {
		return time.Now(), nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("неверный формат даты: %w", err)
	}
	return parsed, nil
}
