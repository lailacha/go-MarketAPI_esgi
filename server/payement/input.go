package payement

type InputPayement struct {
	ProductID int `json:"product_id"`
	PricePaid int `json:"price_paid"`
}