package payement

import (
	"gorm.io/gorm"
)

type Repository interface {
    GetAll() ([]Payement, error)
    GetById(id int) (Payement, error)
    Create(payement Payement) (Payement, error)
    Update(id int, payement InputPayement) (Payement, error)
    Delete(id int) error
}

type repository struct {
    db *gorm.DB
}

func NewPayementRepository(db *gorm.DB) *repository {
    return &repository{db}
}

func (repo *repository) Create(inputPayement Payement) (Payement, error) {

    err := repo.db.Create(&inputPayement).Error

    if err != nil {
        return inputPayement, err
    }

    return inputPayement, nil
}


func (repo *repository) GetById(id int) (Payement, error) {
    var payement Payement
    err := repo.db.Where(&Payement{Id: id}).First(&payement).Error

    if err != nil {
        return Payement{}, err
    }

    return payement, nil
}

func (repo *repository) GetAll() ([]Payement, error) {
    var payements []Payement
    err := repo.db.Find(&payements).Error

    if err != nil {
        return payements, err
    }

    return payements, nil
}

func (repo *repository) Update(id int, inputPayement InputPayement) (Payement, error) {

    
    payement, err := repo.GetById(id)

    if err != nil {
        return payement, err
    }

    
    if(inputPayement.ProductID != 0){
        payement.ProductID = inputPayement.ProductID
    }

    if(inputPayement.PricePaid != 0){
        payement.PricePaid = inputPayement.PricePaid
    }

    err = repo.db.Save(&payement).Error

    if err != nil {
        return payement, err
    }

    return payement, nil
    
}

func (repo *repository) Delete(id int) error {
    
    payement := Payement{Id: id}

    transac := repo.db.Delete(&payement)

    if transac.Error != nil {
        return transac.Error
    }

    if transac.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }

    return nil
}