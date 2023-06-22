package services

import (
	"fmt"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
)

// CreateUser creates a new user
func CreateUser(user models.User) (int, error) {
	userID, err := db.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

// GetUserByID retrieves a user by its ID
func GetUserByID(userID int) (models.User, error) {
	user, err := db.GetUserByID(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUserPhone updates the phone number of a user
func UpdateUserPhone(userID int, phoneNumber string) error {
	err := db.UpdateUserPhone(userID, phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to update user phone: %w", err)
	}

	return nil
}

// AddUserAddress adds a new address for a user
func AddUserAddress(userID int, address models.Address) error {
	err := db.AddUserAddress(userID, address)
	if err != nil {
		return fmt.Errorf("failed to add user address: %w", err)
	}

	return nil
}

// UpdateUserAddress updates an existing address for a user
func UpdateUserAddress(userID int, address models.Address) error {
	err := db.UpdateUserAddress(userID, address)
	if err != nil {
		return fmt.Errorf("failed to update user address: %w", err)
	}

	return nil
}

// GetUserAddresses retrieves all addresses for a user
func GetUserAddresses(userID int) ([]models.Address, error) {
	addresses, err := db.GetUserAddresses(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user addresses: %w", err)
	}

	return addresses, nil
}
