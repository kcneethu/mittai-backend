package models

// ProductWeight represents the weight options for a product
type ProductWeight struct {
	ProductWeightID int    `json:"product_weight_id"`
	ProductID       int    `json:"product_id"`
	Weight          string `json:"weight"`
}
