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
	Password        string      `json:"password"`
	// isactive  		bool        `json:"isactive"`
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

type EmailCheckResponse struct {
	Exists bool `json:"exists"`
}

// UserCreationRequest represents the request body for creating a new user based on phone number.
type UserCreationRequest struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	ContactNumber string `json:"contactNumber"`
}

// UserCreationResponse represents the response body for user creation requests.
type UserCreationResponse struct {
	UserID string `json:"userID"`
}

// MobileCheckResponse represents the JSON response for mobile number existence checks.
type MobileCheckResponse struct {
	Exists  bool   `json:"exists"`
	UserID  string `json:"userID,omitempty"`  // UserID is included only if Exists is true
	Message string `json:"message,omitempty"` // Optional message, e.g., for errors or status updates
}

// ContactNumberRequest represents the JSON structure for checking if a mobile number exists.
type ContactNumberRequest struct {
	ContactNumber string `json:"contact_number"`
}
