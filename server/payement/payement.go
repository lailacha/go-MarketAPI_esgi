package payement

import (
	"github.com/lailacha/go-MarketAPI_esgi/server/product"
	"time"
)

type Payement struct {
	Id   int    `json:"id"`
	ProductID int `json:"product_id"`
	Product product.Product `json:"product"`
	PricePaid string `json:"price_paid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
