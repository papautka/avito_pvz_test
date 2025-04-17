package payload

type UserCreateRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Role     string `json:"role" validate:"required"`
	Password string `json:"password" validate:"min=8"`
}
