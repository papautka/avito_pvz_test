package pvz

import (
	"github.com/google/uuid"
	"time"
)

type PVZ struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

func NewPVZ(RegistrationDate time.Time, City string) *PVZ {
	return &PVZ{
		RegistrationDate: RegistrationDate,
		City:             City,
	}
}
