package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gorilla/mux"
)

// UserService represents a service for user-related operations
type UserService struct {
	DB *db.Repository
}

func NewUserService(db *db.Repository) *UserService {
	return &UserService{
		DB: db,
	}
}

// RegisterRoutes registers the user service routes
func (us *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", us.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", us.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id}", us.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", us.DeleteUser).Methods("DELETE")
}

// CreateUser creates a new user
// @Summary Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 200 {string} string "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 409 {object} ErrorResponse "Contact number already exists"
// @Failure 500 {object} ErrorResponse "Failed to create user"
// @Router /users [post]
func (us *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the contact number is unique
	if us.isContactNumberExists(user.ContactNumber) {
		http.Error(w, "Contact number already exists", http.StatusConflict)
		return
	}

	// Save the user to the database
	err = us.saveUser(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User created successfully"))
}

func (us *UserService) saveUser(user *models.User) error {
	// Save the user to the 'users' table
	query := `INSERT INTO users (first_name, last_name, email, contact_number, verified_account) VALUES (?, ?, ?, ?, ?)`
	result, err := us.DB.Exec(query, user.FirstName, user.LastName, user.Email, user.ContactNumber, user.VerifiedAccount)
	if err != nil {
		return err
	}

	// Get the auto-generated user ID
	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.UserID = int(userID)

	// Save the user's addresses to the 'addresses' table if provided
	if user.Address != nil {
		addresses := *user.Address
		for _, address := range addresses {
			query = `INSERT INTO addresses (user_id, address_line1, address_line2, city, state, zip_code) VALUES (?, ?, ?, ?, ?, ?)`
			_, err := us.DB.Exec(query, user.UserID, address.AddressLine1, address.AddressLine2, address.City, address.State, address.ZipCode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetUserByID retrieves a user by ID
// @Summary Retrieve a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User "User retrieved successfully"
// @Failure 500 {object} ErrorResponse "Failed to retrieve user"
// @Router /users/{id} [get]
func (us *UserService) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := us.getUserByID(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (us *UserService) getUserByID(userID string) (*models.User, error) {
	// Get the user from the 'users' table
	query := `SELECT * FROM users WHERE user_id = ?`
	row := us.DB.QueryRow(query, userID)

	user := &models.User{}
	err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.ContactNumber, &user.VerifiedAccount)
	if err != nil {
		return nil, err
	}

	// Get the user's addresses from the 'addresses' table
	query = `SELECT * FROM addresses WHERE user_id = ?`
	rows, err := us.DB.Query(query, user.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*models.Address
	for rows.Next() {
		address := &models.Address{}
		err := rows.Scan(&address.AddressID, &address.UserID, &address.AddressLine1, &address.AddressLine2, &address.City, &address.State, &address.ZipCode)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	// Assign addresses to the user's Addresses field
	if len(addresses) > 0 {
		user.Address = &addresses
	}
	return user, nil
}

// UpdateUser updates a user
// @Summary Update a user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "User object"
// @Success 200 {string} string "User updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to update user"
// @Router /users/{id} [put]
func (us *UserService) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the user in the database
	err = us.updateUser(userID, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}

// UpdateUser updates a user
func (us *UserService) updateUser(userID string, user *models.User) error {
	// Update the user in the 'users' table
	query := `UPDATE users SET first_name = ?, last_name = ?, email = ?, contact_number = ?, verified_account = ? WHERE user_id = ?`
	_, err := us.DB.Exec(query, user.FirstName, user.LastName, user.Email, user.ContactNumber, user.VerifiedAccount, userID)
	if err != nil {
		return err
	}

	// Delete the user's existing addresses from the 'addresses' table
	query = `DELETE FROM addresses WHERE user_id = ?`
	_, err = us.DB.Exec(query, userID)
	if err != nil {
		return err
	}

	// Insert the updated addresses into the 'addresses' table
	if user.Address != nil {
		addresses := *user.Address
		for _, address := range addresses {
			query = `INSERT INTO addresses (user_id, address_line1, address_line2, city, state, zip_code) VALUES (?, ?, ?, ?, ?, ?)`
			_, err := us.DB.Exec(query, userID, address.AddressLine1, address.AddressLine2, address.City, address.State, address.ZipCode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteUser deletes a user
// @Summary Delete a user
// @Tags Users
// @Param id path string true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 500 {object} ErrorResponse "Failed to delete user"
// @Router /users/{id} [delete]
func (us *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	err := us.deleteUser(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}

func (us *UserService) deleteUser(userID string) error {
	// Delete the user from the 'users' table
	query := `DELETE FROM users WHERE user_id = ?`
	_, err := us.DB.Exec(query, userID)
	if err != nil {
		return err
	}

	// Delete the user's addresses from the 'addresses' table
	query = `DELETE FROM addresses WHERE user_id = ?`
	_, err = us.DB.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}

// isContactNumberExists checks if a contact number already exists in the database
func (us *UserService) isContactNumberExists(contactNumber string) bool {
	query := `SELECT COUNT(*) FROM users WHERE contact_number = ?`
	row := us.DB.QueryRow(query, contactNumber)

	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}

	return count > 0
}
