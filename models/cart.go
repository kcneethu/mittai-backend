package models

import "time"

// Cart represents a user's cart
type Cart struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	Items      []*CartItem `json:"items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	TotalPrice float64     `json:"total_price"` // add this line
}

// CartItem represents an item in the user's cart
type CartItem struct {
	Product   *ProductWeight `json:"product"`
	Quantity  int            `json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// GetTotalPrice calculates the total price of the cart
func (c *Cart) GetTotalPrice() float64 {
	totalPrice := 0.0
	for _, item := range c.Items {
		totalPrice += float64(item.Quantity) * item.Product.Price
	}
	c.TotalPrice = totalPrice
	return c.TotalPrice
}
