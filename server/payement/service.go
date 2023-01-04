package payement

import (
	"time"

	"github.com/lailacha/go-MarketAPI_esgi/server/product"
)

type Service interface {
	GetAll() ([]Payement, error)
	Get(id int) (Payement, error)
	Create(product product.Product) (Payement, error)
	Update(id int, payement InputPayement) (Payement, error)
	Delete(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) Create(product product.Product) (Payement, error) {

	payementObject := Payement{
		ProductID: product.Id,
		PricePaid: product.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newPayement, err := s.repo.Create(payementObject)

    if err != nil {
        return Payement{}, err
    }

	return newPayement, nil
}

func (s *service) Get(id int) (Payement, error) {

	payement, err := s.repo.GetById(id)

	if err != nil {
		return Payement{}, err
	}

	return payement, nil
}

func (s *service) Update(id int, inputPayement InputPayement) (Payement, error) {

	updatedPayement, err := s.repo.Update(id, inputPayement)

	if err != nil {
		return Payement{}, err
	}

	return updatedPayement, nil

}

func (s *service) GetAll() ([]Payement, error) {

	GetAllpayement, err := s.repo.GetAll()

	if err != nil {
		return GetAllpayement, err
	}

	return GetAllpayement, nil
}

func (s *service) Delete(id int) error {

	err := s.repo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
