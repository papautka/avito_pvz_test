package repos

import (
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/internal/users"
)

type AllRepository struct {
	UserRepo *users.UserRepo
	PvzRepo  *pvz.PVZRepo
}

func NewAllRepository(userRepo *users.UserRepo, pvzRepo *pvz.PVZRepo) *AllRepository {
	return &AllRepository{
		UserRepo: userRepo,
		PvzRepo:  pvzRepo,
	}
}
