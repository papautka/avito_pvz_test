package users

import (
	"avito_pvz_test/pkg/database"
	"fmt"
	"log"
)

type UserRepoDeps struct {
	Database *database.Db
}

type UserRepo struct {
	Database *database.Db
}

func NewUserRepo(database *database.Db) *UserRepo {
	return &UserRepo{
		Database: database,
	}
}

func (repo *UserRepo) Create(user *User) (*User, error) {
	query := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id`
	err := repo.Database.MyDb.QueryRow(query, user.Email, user.Password, user.Role).Scan(&user.Id)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return user, nil
}
