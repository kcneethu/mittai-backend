package models

// Order represents an order in the system
type Order struct {
	OrderID         int         `json:"order_id"`
	UserID          int         `json:"user_id"`
	OrderDate       string      `json:"order_date"`
	TotalAmount     float64     `json:"total_amount"`
	Status          string      `json:"status"`
	PaymentMethod   string      `json:"payment_method"`
	DeliveryAddress string      `json:"delivery_address"`
	Items           []OrderItem `json:"items"`
}

// OrderItem represents an item within an order
type OrderItem struct {
	OrderItemID          int     `json:"order_item_id"`
	OrderID              int     `json:"order_id"`
	ProductWeightID      int     `json:"product_weight_id"`
	Quantity             int     `json:"quantity"`
	Price                float64 `json:"price"`
	CustomizationOptions string  `json:"customization_options"`
}
