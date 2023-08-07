package models

import "time"

// Wishlist represents an item in a user's wishlist
type Wishlist struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WishlistRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}
