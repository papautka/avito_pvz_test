package cmd

import (
	"avito_pvz_test/internal/products"
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/internal/receptions"
	"avito_pvz_test/internal/users"
)

type AllRepo struct {
	UsersRepo     users.RepositoryUser
	PvzRepo       pvz.RepositoryPvz
	ReceptionRepo receptions.RepositoryReception
	ProductRepo   products.RepositoryProduct
}

type AllService struct {
	UsersService     users.ServiceUser
	PvzService       pvz.ServicePvz
	ReceptionService receptions.ServiceReception
	ProductService   products.ServiceProduct
}

type AllHandler struct {
	usersHandler     users.HandlerUser
	pvzHandler       pvz.HandlerPvz
	receptionHandler receptions.HandlerReception
	productHandler   products.HandlerProduct
}
