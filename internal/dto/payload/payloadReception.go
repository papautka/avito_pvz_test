package payload

type ReceptionCreateRequest struct {
	PvzId string `json:"pvzId" validate:"required"`
}
