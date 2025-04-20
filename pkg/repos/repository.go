package repos

import (
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/internal/receptions"
	"avito_pvz_test/internal/users"
)

type AllRepository struct {
	UserRepo      *users.UserRepo
	PvzRepo       *pvz.PVZRepo
	ReceptionRepo *receptions.ReceptionRepo
}

func NewAllRepository(userRepo *users.UserRepo, pvzRepo *pvz.PVZRepo, recRepo *receptions.ReceptionRepo) *AllRepository {
	return &AllRepository{
		UserRepo:      userRepo,
		PvzRepo:       pvzRepo,
		ReceptionRepo: recRepo,
	}
}
