package models

// Address represents a user address
type Address struct {
	AddressID int    `json:"addressID"`
	UserID    int    `json:"userID"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zipCode"`
}
