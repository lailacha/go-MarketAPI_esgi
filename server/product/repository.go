package product

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]Product, error)
	GetById(id int) (Product, error)
	GetByName(name string) (Product, error)
	Create(product Product) Product
	Update(id int, inputProduct InputProduct) (Product, error)
	Delete(id int) error
}

type repository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (pr *repository) GetAll() ([]Product, error) {
	var products []Product

	err := pr.db.Find(&products).Error

	if err != nil {
		return products, err
	}

	return products, nil
}

func (pr *repository) GetById(id int) (Product, error) {
	var product Product
	err := pr.db.Where(&Product{Id: id}).First(&product).Error

	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (pr *repository) GetByName(name string) (Product, error) {
	var product Product
	err := pr.db.Where(&Product{Name: name}).First(&product).Error

	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (pr *repository) Create(product Product) Product {
	pr.db.Create(&product)
	return product
}

func (pr *repository) Update(id int, inputProduct InputProduct) (Product, error) {

	product, err := pr.GetById(id)

	if err != nil {
		return Product{}, err
	}

	if inputProduct.Name != "" {
		product.Name = inputProduct.Name
	}

	if inputProduct.Price != "" {
		product.Price = inputProduct.Price
	}

	pr.db.Save(&product)

	return product, nil
}

func (pr *repository) Delete(id int) error {

	product := Product{Id: id}

	transaction := pr.db.Delete(&product)

	if transaction.Error != nil {
		return transaction.Error
	}

	if transaction.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
