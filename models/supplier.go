package models

// Supplier represents a supplier in the system
type Supplier struct {
	SupplierID    int    `json:"supplier_id"`
	Name          string `json:"name"`
	ContactNumber string `json:"contact_number"`
	Email         string `json:"email"`
}
