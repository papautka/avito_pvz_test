package pvz

import (
	"avito_pvz_test/pkg/req"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type ServicePvz interface {
	Register(id string, registrationDate string, city string) (*PVZ, error)
	ChangeStatusReceptionByPvzOnClose(id string) (*ReceptionForPvz, error)
	GetArrayPvz(filter *req.FilterWithPagination) (*PvzListResponse, error)
	DeleteProduct(pvzId string) (*Product, error)
}

type ServPvz struct {
	pvzRepoInterface RepositoryPvz
}

func NewServPvz(repo RepositoryPvz) ServicePvz {
	return &ServPvz{
		pvzRepoInterface: repo,
	}
}

func (pvz *ServPvz) Register(id string, registrationDate string, city string) (*PVZ, error) {
	// проверяем UUID если не передан
	var uuidVal uuid.UUID
	var err error

	uuidVal, err = req.ParseUUIDOrGenerate(id)
	if err != nil {
		return nil, err
	}

	// обработка даты
	var regTime time.Time
	regTime, err = req.ParseTimeOrNow(registrationDate)
	if err != nil {
		return nil, err
	}

	newPvz := NewPVZ(uuidVal, regTime, city)
	createdPVZ, err := pvz.pvzRepoInterface.Create(newPvz)
	if err != nil {
		log.Printf("Error creating PVZ: %v", err)
		return nil, err
	}
	return createdPVZ, nil
}

func (pvz *ServPvz) ChangeStatusReceptionByPvzOnClose(id string) (*ReceptionForPvz, error) {
	// проверяем UUID если не передан
	var uuidVal uuid.UUID
	var err error

	uuidVal, err = req.ParseUUIDPvz(id)
	if err != nil {
		return nil, fmt.Errorf("некорректное значение id")
	}
	pvzStruct, err := pvz.pvzRepoInterface.FindPVZById(uuidVal)
	if err != nil {
		return nil, fmt.Errorf("нет pvz c таким значением")
	}
	recepForPvz, err := pvz.pvzRepoInterface.UpdateStatus(pvzStruct.ID)
	if err != nil {
		return nil, fmt.Errorf("у данного pvzId нет приемок или она закрыта")
	}
	return recepForPvz, nil
}

func (pvz *ServPvz) GetArrayPvz(filter *req.FilterWithPagination) (*PvzListResponse, error) {
	slicePvz, err := pvz.pvzRepoInterface.GetPVZPageAndLimit(filter)
	if err != nil {
		return nil, err
	}
	return slicePvz, nil
}

func (pvz *ServPvz) DeleteProduct(pvzId string) (*Product, error) {
	// проверяем UUID если не передан или если не корректен
	var uuidVal uuid.UUID
	var err error
	uuidVal, err = req.ParseUUIDPvz(pvzId)
	if err != nil {
		return nil, fmt.Errorf("некорректное значение id")
	}
	productStruct, err := pvz.pvzRepoInterface.DeleteLastProductInOpenReception(uuidVal)
	if err != nil {
		return nil, err
	}
	return productStruct, nil
}
