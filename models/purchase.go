package models

import "time"

// Purchase represents a purchase made by a user
type Purchase struct {
	ID         int             `json:"id"`
	UserID     int             `json:"user_id"`
	AddressID  int             `json:"address_id"`
	PaymentID  int             `json:"payment_id"`
	TotalPrice float64         `json:"total_price"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	Items      []*PurchaseItem `json:"items"`
}

type CreatePurchase struct {
	AddressID     int            `json:"address_id"`
	PaymentID     int            `json:"payment_id"`
	PurchaseItems []PurchaseItem `json:"purchase_items"`
	UserID        int            `json:"user_id"`
}

// PurchaseItem represents a purchased item
type PurchaseItem struct {
	ProductID       int     `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductWeightID int     `json:"product_weight_id"`
	ProductPrice    float64 `json:"product_price"`
	Quantity        int     `json:"quantity"`
	TotalPrice      float64 `json:"total_price"`
}

// PurchaseRequest represents the request payload for creating a purchase
type PurchaseRequest struct {
	UserID              int                   `json:"user_id"`
	AddressID           int                   `json:"address_id"`
	PaymentID           int                   `json:"payment_id"`
	PurchaseItemRequest []PurchaseItemRequest `json:"purchase_items"`
}

type PurchaseItemRequest struct {
	ProductID       int `json:"product_id"`
	ProductWeightID int `json:"product_weight_id"`
	Quantity        int `json:"quantity"`
}

// PurchaseResponse represents the response payload for a created purchase
type PurchaseResponse struct {
	PurchaseID int `json:"purchase_id"`
}
