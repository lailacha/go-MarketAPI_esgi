package payement

import (
    "github.com/lailacha/go-MarketAPI_esgi/server/product"
)

type Service interface {
    FindAll() ([]Payement, error)
	FindById(id int) (Payement, error)
    Get(id int) (Payement, error)
    Create(inputProduct product.Product) (Payement, error)
    Update(id int, payement Payement) (Payement, error)
    Delete(id int) error
    // Stream() (Payement, error)
}

type service struct {
    repo Repository
}

func NewService(repo Repository) *service {
    return &service{repo}
}

func (s *service) Create(inputProduct product.Product) (Payement, error) {

    payementObject := Payement{
        ProductID: inputProduct.Id,
        PricePaid: inputProduct.Price,
    }
    
    s.repo.Create(payementObject)
    
    
    return payementObject, nil
}

func (s *service) Get(id int) (Payement, error) {
    
    payement, err := s.repo.FindById(id)

    if err != nil {
        return Payement{}, err
    }

    return payement, nil
}


func (s *service) Update(id int, inputPayement Payement) (Payement, error) {


    updatedPayement, err := s.repo.Update(id, inputPayement)

    if err != nil {
        return Payement{}, err
    }

    return updatedPayement, nil

}

func (s *service) FindAll() ([]Payement, error) {

	findallpayement, err := s.repo.FindAll()

	if err != nil {
		return findallpayement, err
	}

	return findallpayement, nil
}

func (s *service) Delete(id int) error {
    
	err := s.repo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

// func (s *service) Stream() (Payement, error) {
//     //return s.repo.Stream()
// } 