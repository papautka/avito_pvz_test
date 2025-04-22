package products

import (
	"avito_pvz_test/pkg/database"
	"fmt"
)

type RepositoryProduct interface {
	Create(product *Product) (*Product, error)
}

type RepoProduct struct {
	Database *database.Db
}

func NewRepoProduct(database *database.Db) RepositoryProduct {
	return &RepoProduct{
		Database: database,
	}
}

func (repo *RepoProduct) Create(product *Product) (*Product, error) {
	query := `INSERT INTO products (dateTime, type, receptionId) VALUES ($1, $2, $3) RETURNING id`
	err := repo.Database.MyDb.QueryRow(query, product.DateTime, product.Type, product.ReceptionId).Scan(&product.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании товара: %w", err)
	}
	return product, nil
}
