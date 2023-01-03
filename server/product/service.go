package product

import "fmt"

// on définit notre interface
type Service interface {
	FindAll() ([]Product, error)
	Get(id int) (Product, error)
	Create(name string, price string) (Product, error)
	Update(id int, product Product) (Product, error)
	Delete(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) Create(name string, price string) (Product, error) {

	fmt.Println("service create")


	// we verify if the name already exists in the db, otherwise we create the product
	_, err := s.repo.GetByName(name)

	if err == nil {
		return Product{}, fmt.Errorf("product already exists")
	}


	productObject := Product{
		Name: name,
		Price: price,
	}

	 s.repo.Create(productObject)

	return productObject, nil
}


func (s *service) Update(id int, inputProduct Product) (Product, error) {

	// we verify if the name already exists in the db, otherwise we create the product
	_, err := s.repo.GetByName(inputProduct.Name)

	if err == nil {
		return Product{}, fmt.Errorf("product already exists")
	}

	updatedProduct, err := s.repo.Update(id, inputProduct)

	if err != nil {
		return Product{}, err
	}

	return updatedProduct, nil

}	


func (s *service) Get(id int) (Product, error) {
	
	product, err := s.repo.GetById(id)

	if err != nil {
		return Product{}, err
	}

	return product, nil

}



func (s *service) Delete(id int) error {

	err := s.repo.Delete(id)

	if err != nil {
		return err
	}

	return nil

}

func (s *service) FindAll() ([]Product, error) {

	products, err := s.repo.GetAll()

	if err != nil {
		return products, err
	}

	return products, nil
}