package receptions

import (
	"github.com/google/uuid"
	"time"
)

type Reception struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"date_time"`
	PvzID    uuid.UUID `json:"pvz_id"`
	Status   string    `json:"status"` // enum: [in_progress, close]
}

//required: [dateTime, pvzId, status]

func NewReception(dateTime time.Time, pvzID uuid.UUID, status string) *Reception {
	return &Reception{
		DateTime: dateTime,
		PvzID:    pvzID,
		Status:   status,
	}
}
