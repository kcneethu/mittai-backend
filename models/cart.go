package models

// Cart represents a user's cart in the system
type Cart struct {
	CartID           int    `json:"cart_id"`
	UserID           int    `json:"user_id"`
	ProductWeightID  int    `json:"product_weight_id"`
	Quantity         int    `json:"quantity"`
	AppliedDiscounts string `json:"applied_discounts"`
}
