package models

// PaymentMode represents a payment mode
type PaymentMode struct {
	ID       int    `json:"id"`
	Mode     string `json:"mode"`
	IsActive bool   `json:"is_active"`
}
