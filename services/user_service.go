package services

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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
	r.HandleFunc("/login", us.Login).Methods("POST") // Add this line for the login route
	r.HandleFunc("/verify-otp/{id}", us.VerifyOTP).Methods("POST")
	r.HandleFunc("/users/{id}/name", us.GetUserNameByID).Methods("GET") // Add this new route
	// The new route for checking if a mobile number exists
	r.HandleFunc("/users/check-mobile", us.CheckMobileNumberExists).Methods("POST")
	// The new route for creating a user with minimal details
	r.HandleFunc("/users/create", us.CreateUserWithDetails).Methods("POST")
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

	// Send the response with the generated userID
	response := map[string]interface{}{
		"message": "User created successfully",
		"user_id": user.UserID,
	}

	otp := generateOTP()
	// Save OTP to the database
	err = us.saveOTP(user.UserID, otp)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Send OTP email

	if user.Email.Valid {
		err = us.sendOTPEmail(user.Email.String, otp, user.UserID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to send OTP email", http.StatusInternalServerError)
			return
		}
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func (us *UserService) saveUser(user *models.User) error {
	// Save the user to the 'users' table
	query := `INSERT INTO users (first_name, last_name, email, contact_number, verified_account, hashed_password) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := us.DB.Exec(query, user.FirstName, user.LastName, user.Email, user.ContactNumber, user.VerifiedAccount, user.Password)
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
	err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.ContactNumber, &user.VerifiedAccount, &user.Password)
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

// hashPassword hashes the user's password using bcrypt
func (us *UserService) hashPassword(password string) (string, error) {
	// Hash the password using bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Base64 encode the hashed password before returning it as a string
	hashedPassword := base64.StdEncoding.EncodeToString(hashedBytes)
	return hashedPassword, nil
}

// GetUserNameByID retrieves the first and last name of a user by ID
// @Summary Retrieve the first and last name of a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User "User name retrieved successfully"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Failed to retrieve user name"
// @Router /users/{id}/name [get]
func (us *UserService) GetUserNameByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := us.getUserByID(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve user name", http.StatusInternalServerError)
		return
	}

	// Create a response with the user's first and last name
	response := map[string]interface{}{
		"userID":    user.UserID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CheckMobileNumberExists checks if a mobile number exists in the users table.
// @Summary Check if a mobile number exists
// @Description Checks if a given mobile number exists in the database and returns a response indicating the existence.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.ContactNumberRequest true "Request body containing the contact number to check"
// @Success 200 {object} models.MobileCheckResponse "Mobile number check response indicating whether the mobile number exists along with the user ID if it does."
// @Failure 400 {object} ErrorResponse "Invalid request - JSON body is required and must contain a 'contact_number'."
// @Failure 500 {object} ErrorResponse "Internal server error - Failed to check mobile number due to a server error."
// @Router /users/check-mobile [post]
func (us *UserService) CheckMobileNumberExists(w http.ResponseWriter, r *http.Request) {
	var req models.ContactNumberRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.ContactNumber == "" {
		http.Error(w, "Invalid request: contact_number is required", http.StatusBadRequest)
		return
	}

	userID, err := us.getUserIDByContactNumber(req.ContactNumber)
	if err != nil {
		log.Printf("Error retrieving user by contact number: %v", err)
		http.Error(w, "Failed to check mobile number", http.StatusInternalServerError)
		return
	}

	response := models.MobileCheckResponse{}
	if userID > 0 {
		response.Exists = true
		response.UserID = fmt.Sprintf("%d", userID)
	} else {
		response.Exists = false
		response.Message = "Mobile number not found"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (us *UserService) getUserIDByContactNumber(contactNumber string) (int, error) {
	var userID int
	query := `SELECT user_id FROM users WHERE contact_number = ?`
	err := us.DB.QueryRow(query, contactNumber).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No results found is an expected outcome, not an error. Return 0 for userID and nil for error.
			return 0, nil
		}
		// An actual error occurred, return 0 and the error.
		return 0, err
	}

	return userID, nil // User found, return the userID and nil for error.
}

// CreateUserWithDetails creates a new user with details
// @Summary Create a new user with first name, last name, and contact number
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserCreationRequest true "User creation request"
// @Success 200 {object} models.UserCreationResponse "User creation response"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 409 {object} ErrorResponse "Contact number already exists"
// @Router /users/create [post]
func (us *UserService) CreateUserWithDetails(w http.ResponseWriter, r *http.Request) {
	var request models.UserCreationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if us.isContactNumberExists(request.ContactNumber) {
		response := models.UserCreationResponse{UserID: "Contact Number Already Exists"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	userID, err := us.createUserWithDetails(request)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response := models.UserCreationResponse{UserID: fmt.Sprintf("%d", userID)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (us *UserService) createUserWithDetails(request models.UserCreationRequest) (int, error) {
	// Assume the UserCreationRequest struct and UserCreationResponse struct are defined elsewhere
	query := `INSERT INTO users (first_name, last_name, contact_number) VALUES (?, ?, ?)`
	result, err := us.DB.Exec(query, request.FirstName, request.LastName, request.ContactNumber)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}
