package products

import (
	"avito_pvz_test/internal/receptions"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ServiceProduct interface {
	Create(pvzID uuid.UUID, typeProducts string) (*Product, error)
}

type ServProduct struct {
	repoProducts   RepositoryProduct
	repoReceptions receptions.RepositoryReception
}

func NewServProduct(repoProduct RepositoryProduct, repoReception receptions.RepositoryReception) ServiceProduct {
	return &ServProduct{
		repoProducts:   repoProduct,
		repoReceptions: repoReception,
	}
}

func (s *ServProduct) Create(pvzID uuid.UUID, typeProducts string) (*Product, error) {
	// 1. TODO убедиться что данный pvz вообще существует

	// 2. TODO если PVZ есть то убедиться что у него есть последняя приемка и она не закрыта
	// 3. TODO если 1-ое или 2-ое условие не выполняется то вернуть ошибку (_Неверный запрос или нет активной приемки_)
	reception, err := s.repoReceptions.ReturnLastReceptionOrEmpty(pvzID)
	if err != nil {
		return nil, err
	}
	// 4. TODO если и PVZ такой есть и у него приемка в статусе 'in_progress' то надо создать product
	if reception.Status != "in_progress" {
		return nil, fmt.Errorf("статус приемки должен быть in_progress")
	}
	// 5. TODO достать ReceptionId у pvzID

	product := &Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        typeProducts,
		ReceptionId: reception.ID,
	}

	resProduct, err := s.repoProducts.Create(product)
	if err != nil {
		return nil, err
	}
	return resProduct, nil
}
