package models

type OrderStatus struct {
	ID         int    `json:"id"`
	PurchaseID int    `json:"purchase_id"`
	Status     string `json:"status"`
}
