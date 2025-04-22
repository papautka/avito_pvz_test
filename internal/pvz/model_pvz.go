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

type ReceptionForPvz struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"date_time"`
	PvzID    uuid.UUID `json:"pvz_id"`
	Status   string    `json:"status"` // enum: [in_progress, close]
}

func NewPVZ(Id uuid.UUID, RegistrationDate time.Time, City string) *PVZ {
	return &PVZ{
		ID:               Id,
		RegistrationDate: RegistrationDate,
		City:             City,
	}
}
