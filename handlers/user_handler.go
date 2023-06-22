package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gorilla/mux"
)

// UsersHandler handles requests related to users
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listUsers(w, r)
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listUsers retrieves all users from the database and writes the response
// @Summary Retrieve all users
// @Description Retrieves all users from the database
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.ListUsers()
	if err != nil {
		log.Printf("failed to get users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		log.Printf("failed to marshal users response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// createUser creates a new user in the database
// @Summary Create a new user
// @Description Creates a new user in the database
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("failed to decode user data: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := db.CreateUser(user)
	if err != nil {
		log.Printf("failed to create user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the created user ID in the response
	resp := CreateUserResponse{
		UserID: userID,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Printf("failed to marshal user ID response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respJSON)
}

// UpdateUserPhone updates the phone number for a user
// @Summary Update user phone number
// @Description Updates the phone number for a user
// @Tags users
// @Accept json
// @Param request body UpdateUserPhoneRequest true "Update user phone request object"
// @Success 200 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/phone [put]
func UpdateUserPhone(w http.ResponseWriter, r *http.Request) {
	var request UpdateUserPhoneRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("failed to decode update user phone request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.UpdateUserPhone(request.UserID, request.PhoneNumber)
	if err != nil {
		log.Printf("failed to update user phone: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateUserAddress updates an existing address for a user
// @Summary Update user address
// @Description Updates an existing address for a user
// @Tags users
// @Accept json
// @Param request body UpdateUserAddressRequest true "Update user address request object"
// @Success 200 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/address [put]
func UpdateUserAddress(w http.ResponseWriter, r *http.Request) {
	var request UpdateUserAddressRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("failed to decode update user address request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.UpdateUserAddress(request.UserID, request.Address)
	if err != nil {
		log.Printf("failed to update user address: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AddUserAddress adds a new address for a user
// @Summary Add user address
// @Description Adds a new address for a user
// @Tags users
// @Accept json
// @Param request body AddUserAddressRequest true "Add user address request object"
// @Success 200 {object} EmptyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/address [post]
func AddUserAddress(w http.ResponseWriter, r *http.Request) {
	var request AddUserAddressRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("failed to decode add user address request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.AddUserAddress(request.UserID, request.Address)
	if err != nil {
		log.Printf("failed to add user address: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetUserAddresses retrieves all addresses for a user
// @Summary Get user addresses
// @Description Retrieves all addresses for a user
// @Tags users
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} models.Address
// @Failure 500 {object} ErrorResponse
// @Router /users/{userID}/addresses [get]
func GetUserAddresses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIDStr, ok := params["userID"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addresses, err := db.GetUserAddresses(userID)
	if err != nil {
		log.Printf("failed to get user addresses: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(addresses)
	if err != nil {
		log.Printf("failed to marshal user addresses response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// FlushAllTables deletes all records from all tables
// @Summary Flush all tables
// @Description Deletes all records from all tables
// @Tags administration
// @Success 200 {object} EmptyResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/flush [delete]
func FlushAllTables(w http.ResponseWriter, r *http.Request) {
	err := db.FlushAllTables()
	if err != nil {
		log.Printf("failed to flush all tables: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CreateUserResponse represents the response for creating a new user
type CreateUserResponse struct {
	UserID int `json:"userID"`
}

// UpdateUserPhoneRequest represents the request to update user phone number
type UpdateUserPhoneRequest struct {
	UserID      int    `json:"userID"`
	PhoneNumber string `json:"phoneNumber"`
}

// UpdateUserAddressRequest represents the request to update user address
type UpdateUserAddressRequest struct {
	UserID  int            `json:"userID"`
	Address models.Address `json:"address"`
}

// AddUserAddressRequest represents the request to add user address
type AddUserAddressRequest struct {
	UserID  int            `json:"userID"`
	Address models.Address `json:"address"`
}

// EmptyResponse represents an empty response
type EmptyResponse struct{}
