package payload

type PvzCreateRequest struct {
	Id               string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city" validate:"required"`
}
