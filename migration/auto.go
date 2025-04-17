package migration

import (
	"avito_pvz_test/internal/handler/auth"
	"avito_pvz_test/variable"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(variable.Env_file)
	if err != nil {
		log.Fatal(variable.Msg_err_env)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv(variable.KEY_DNS_DB)), &gorm.Config{})
	if err != nil {
		log.Fatal(variable.Msg_err_open_db)
	}
	err = db.AutoMigrate(&auth.User{})
	if err != nil {
		log.Fatal(variable.Msg_err_auto_migrate)
	}
	fmt.Println(variable.Msg_suc_auto_migrate)

}
