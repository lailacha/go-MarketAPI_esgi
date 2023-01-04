package payement

import (
	"time"
)

type Payement struct {
	Id   int    `json:"id"`
	ProductID int `json:"product_id"`
	PricePaid int `json:"price_paid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
