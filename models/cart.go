package models

import "time"

// Cart represents a user's cart
type Cart struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Items     []*CartItem `json:"items"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// CartItem represents an item in the user's cart
type CartItem struct {
	Product   *ProductWeight `json:"product"`
	Quantity  int            `json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
