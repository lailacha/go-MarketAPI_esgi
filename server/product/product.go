package product

import "time"

// on définit notre classe
type Product struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Price string `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
 
