package models

// Promotion represents a product promotion in the system
type Promotion struct {
	PromotionID        int     `json:"promotion_id"`
	ProductID          int     `json:"product_id"`
	DiscountCode       string  `json:"discount_code"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
	DiscountPercentage float64 `json:"discount_percentage"`
}
