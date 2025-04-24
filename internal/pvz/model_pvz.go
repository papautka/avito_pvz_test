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

type Reception struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"date_time"`
	PvzID    uuid.UUID `json:"pvz_id"`
	Status   string    `json:"status"` // enum: [in_progress, close]
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	DateTime    time.Time `json:"date_time"`
	Type        string    `json:"type"`
	ReceptionId uuid.UUID `json:"reception_id"`
}

type ReceptionResponse struct {
	Reception    Reception `json:"reception"`
	ArrayProduct []Product `json:"array_product"`
}

type PvzResponse struct {
	Pvz            PVZ                 `json:"pvz"`
	ArrayReception []ReceptionResponse `json:"array_reception"`
}

type PvzListResponse struct {
	ArrayPvzResponse []PvzResponse `json:"array_pvz"`
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
