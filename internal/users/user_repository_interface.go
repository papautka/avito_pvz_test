package users

import (
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	FindUserByEmailPass(ctx context.Context, email, password string) (*User, error)
}
