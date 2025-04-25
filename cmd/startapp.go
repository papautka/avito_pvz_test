package main

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/products"
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/internal/receptions"
	"avito_pvz_test/internal/users"
	"avito_pvz_test/pkg/database"
	"avito_pvz_test/pkg/midware"
	"net/http"
)

func StartApp() {
	server := http.Server{
		Addr:    ":8080",
		Handler: CreateRouter(),
	}
	server.ListenAndServe()
}

func loadConfig() *config.Config {
	return config.NewConfig()
}

func loadDB(conf *config.Config) *database.Db {
	return database.CreateDb(conf)
}

func loadRepository(db *database.Db) *AllRepo {
	userRepo := users.NewRepoUser(db)
	pvzRepo := pvz.NewRepoPVZ(db)
	receptionRepo := receptions.NewRepoRecep(db)
	productRepo := products.NewRepoProduct(db)
	return &AllRepo{
		UsersRepo:     userRepo,
		PvzRepo:       pvzRepo,
		ReceptionRepo: receptionRepo,
		ProductRepo:   productRepo,
	}

}

func loadService(a *AllRepo, c *config.Config) *AllService {
	userService := users.NewServUser(a.UsersRepo, c)
	pvzService := pvz.NewServPvz(a.PvzRepo)
	receptionService := receptions.NewServReception(a.ReceptionRepo)
	productsService := products.NewServProduct(a.ProductRepo, a.ReceptionRepo)
	return &AllService{
		userService,
		pvzService,
		receptionService,
		productsService,
	}
}

func loadHandlers(s *AllService) *AllHandler {
	userHandler := users.NewHandUser(s.UsersService)
	pvzHandler := pvz.NewHandPvz(s.PvzService)
	receptionHandler := receptions.NewReceptionHandler(s.ReceptionService)
	productHandler := products.NewHandProduct(s.ProductService)
	return &AllHandler{
		userHandler,
		pvzHandler,
		receptionHandler,
		productHandler,
	}

}

func connectHandlers(hs *AllHandler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", hs.usersHandler.CreateUser())
	router.HandleFunc("POST /login", hs.usersHandler.AuthenticateUser())
	router.HandleFunc("POST /dummyLogin", hs.usersHandler.GetTokenByRole())

	router.Handle("POST /pvz", midware.CheckRoleByToken(hs.pvzHandler.CreatePVZ(), "moderator"))
	router.Handle("POST /pvz/{pvzId}/close_last_reception", midware.CheckRoleByToken(hs.pvzHandler.CloseLastReceptionByPvz(), "client"))
	/* в работе */
	router.HandleFunc("GET /pvz", hs.pvzHandler.GetArrayPvz())
	/* в работе */

	router.Handle("POST /receptions", midware.CheckRoleByToken(hs.receptionHandler.CreateReception(), "client"))

	router.Handle("POST /products", midware.CheckRoleByToken(hs.productHandler.Create(), "client"))
	return router
}

func CreateRouter() *http.ServeMux {
	conf := loadConfig()
	db := loadDB(conf)
	repo := loadRepository(db)
	services := loadService(repo, conf)
	hs := loadHandlers(services)
	router := connectHandlers(hs)
	return router
}
