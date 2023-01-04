package payement

type InputPayement struct {
	ProductID int `json:"product_id"`
	PricePaid string `json:"price_paid"`
}