package payload

import "github.com/google/uuid"

type ProductCreateRequest struct {
	Type  string    `json:"type" validate:"required,oneof=электроника одежда обувь"` // [электроника, одежда, обувь]
	PvzId uuid.UUID `json:"pvzId" validate:"required,uuid4"`                         // facd8332-1a63-4bba-88e0-08bc88f7a30e
}
