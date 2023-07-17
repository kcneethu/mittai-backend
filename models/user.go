package models

// User represents a user in the system
type User struct {
	UserID          int         `json:"userID"`
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	Email           string      `json:"email"`
	ContactNumber   string      `json:"contactNumber"`
	Address         *[]*Address `json:"address"`
	VerifiedAccount bool        `json:"verifiedAccount"`
}

// Address represents a user's address
type Address struct {
	AddressID    int    `json:"addressID"`
	UserID       int    `json:"userID"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zipCode"`
}
