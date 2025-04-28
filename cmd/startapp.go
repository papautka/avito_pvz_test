package cmd

import (
	"avito_pvz_test/config"
	"avito_pvz_test/internal/products"
	"avito_pvz_test/internal/pvz"
	"avito_pvz_test/internal/receptions"
	"avito_pvz_test/internal/users"
	"avito_pvz_test/pkg/database"
	"avito_pvz_test/pkg/midware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func CreateRouterTest() *AllData {
	conf := loadConfig()
	db := loadDB(conf)
	repo := loadRepository(db)
	services := loadService(repo, conf)
	hs := loadHandlers(services)
	router := connectHandlers(hs)
	return &AllData{
		Db:       db,
		Repo:     repo,
		Services: services,
		Handlers: hs,
		Router:   router,
	}
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

func StartApp() {
	app := CreateRouter()
	server := http.Server{
		Addr:    ":8080",
		Handler: app,
	}
	// добавил graceful shutdown (зачем он нужен?)
	// чтобы сервер не вырубался без завершения активных соединений

	// канал для получения сигнала завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// запуск Сервера в отдельной горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	log.Printf("server listening at %s", server.Addr)
	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}
	log.Println("Server exiting gracefully")
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
