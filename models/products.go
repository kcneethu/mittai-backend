package models

// Product represents a product in the system
type Product struct {
	ProductID              int             `json:"product_id"`
	Name                   string          `json:"name"`
	Description            string          `json:"description"`
	Category               string          `json:"category"`
	Price                  float64         `json:"price"`
	Availability           int             `json:"availability"`
	Ingredients            string          `json:"ingredients"`
	NutritionalInformation string          `json:"nutritional_information"`
	ImageURL               string          `json:"image_url"`
	ProductWeights         []ProductWeight `json:"product_weights"`
}
