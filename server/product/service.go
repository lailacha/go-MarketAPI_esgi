package product

import "fmt"

// on définit notre interface
type Service interface {
	// GetAll() ([]Product, error)
	Get(id int) (Product, error)
	Create(inputProduct InputProduct) (Product, error)
	Update(id int, inputProduct InputProduct) (Product, error)
	Delete(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) Create(inputProduct InputProduct) (Product, error) {

	fmt.Println("service create")

	// we verify if the name already exists in the db, otherwise we create the product
	_, err := s.repo.GetByName(inputProduct.Name)

	if err == nil {
		return Product{}, fmt.Errorf("product already exists")
	}

	product := Product{
		Name:  inputProduct.Name,
		Price: inputProduct.Price,
	}

	newProduct := s.repo.Create(product)

	return newProduct, nil
}

func (s *service) Update(id int, inputProduct InputProduct) (Product, error) {


	// we verify if the pdouct already exists

	oldP, err := s.repo.GetById(id);

	if(err != nil) {
		return Product{}, fmt.Errorf("product doesn't exists")
	}


	if(inputProduct.Name != oldP.Name) || err != nil {

		// we verify if the name already exists in the db, otherwise we create the product
		_, err := s.repo.GetByName(inputProduct.Name)

		if err == nil {
			return Product{}, fmt.Errorf("product already exists")
		}

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

func (s *service) GetAll() ([]Product, error) {

	products, err := s.repo.GetAll()

	if err != nil {
		return products, err
	}

	return products, nil
}