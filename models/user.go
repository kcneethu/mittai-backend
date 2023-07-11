package models

type User struct {
	UserID        int      `json:"userID"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	ContactNumber string   `json:"contactNumber"`
	Address       Address  `json:"address"`
}
