package cmd

import (
	"avito_pvz_test/internal/products"
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/internal/receptions"
	"avito_pvz_test/internal/users"
	"avito_pvz_test/pkg/database"
	"net/http"
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

type AllData struct {
	Db       *database.Db
	Repo     *AllRepo
	Services *AllService
	Handlers *AllHandler
	Router   *http.ServeMux
}
