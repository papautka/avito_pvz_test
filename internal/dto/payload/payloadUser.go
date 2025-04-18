package payload

type UserCreateRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Role     string `json:"role" validate:"required"`
	Password string `json:"password" validate:"min=8"`
}

type UserAuthRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type TokenRequestRole struct {
	Role string `json:"role" validate:"required,oneof=employee moderator"`
}
