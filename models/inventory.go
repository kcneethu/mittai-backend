package models

// Inventory represents the inventory for a specific product weight
type Inventory struct {
	ProductWeightID   int `json:"product_weight_id"`
	AvailableQuantity int `json:"available_quantity"`
	ReorderThreshold  int `json:"reorder_threshold"`
	SupplierID        int `json:"supplier_id"`
}
