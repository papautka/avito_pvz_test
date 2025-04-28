package users

import (
	"avito_pvz_test/pkg/database"
	"fmt"
	"log"
)

type RepositoryUser interface {
	CreateUser(user *User) (*User, error)
	DropUser(email string) error
	FindUserByEmailPass(email, password string) (*User, error)
}

type RepoUser struct {
	Database *database.Db
}

func NewRepoUser(database *database.Db) RepositoryUser {
	return &RepoUser{
		Database: database,
	}
}

func (repo *RepoUser) CreateUser(user *User) (*User, error) {
	query := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id`
	err := repo.Database.MyDb.QueryRow(query, user.Email, user.Password, user.Role).Scan(&user.Id)
	if err != nil {
		log.Println("Create ", err)
		return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return user, nil
}

func (repo *RepoUser) FindUserByEmailPass(email, password string) (*User, error) {

	// запрос
	query := `SELECT id, email, password, role FROM users where email = $1 and password = $2`
	// из бд достаем данные и инициализируем структуру
	var user User
	err := repo.Database.MyDb.QueryRow(query, email, password).Scan(&user.Id, &user.Email, &user.Password, &user.Role)
	if err != nil {
		log.Println("FindUserByEmailPass: ", err)
		return nil, fmt.Errorf("нет пользователя с логином и паролем %s %s", email, password)
	}
	return &user, nil
}

func (repo *RepoUser) DropUser(email string) error {
	query := `DELETE FROM users WHERE email = $1`
	_, err := repo.Database.MyDb.Exec(query, email)
	if err != nil {
		log.Println("Drop ", err)
	}
	return err
}
