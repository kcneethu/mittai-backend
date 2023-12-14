package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gorilla/mux"
)

// AddressService represents a service for address-related operations
type AddressService struct {
	DB *db.Repository
}

// NewAddressService creates a new AddressService instance
func NewAddressService(db *db.Repository) *AddressService {
	return &AddressService{
		DB: db,
	}
}

type EmailCheckRequest struct {
	Email string `json:"email"`
}

// RegisterRoutes registers the address service routes
func (as *AddressService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/addresses", as.CreateAddress).Methods("POST")
	r.HandleFunc("/addresses/{id}", as.GetAddressByID).Methods("GET")
	r.HandleFunc("/addresses/{id}", as.UpdateAddress).Methods("PUT")
	r.HandleFunc("/addresses/{id}", as.DeleteAddress).Methods("DELETE")
	r.HandleFunc("/users/{user_id}/addresses", as.GetAddressesByUserID).Methods("GET")
	r.HandleFunc("/users/check-email", as.CheckEmailExists).Methods("POST")

}

// CreateAddress creates a new address
// @Summary Create a new address
// @Tags Addresses
// @Accept json
// @Produce json
// @Param address body models.Address true "Address object"
// @Success 200 {object} CreateAddressResponse "Address created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to create address"
// @Router /addresses [post]
func (as *AddressService) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var address models.Address
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Save the address to the database and get the newly created address ID
	addressID, err := as.saveAddress(&address)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create address", http.StatusInternalServerError)
		return
	}

	// Send the response with the newly created address ID
	w.Header().Set("Content-Type", "application/json")
	response := CreateAddressResponse{
		AddressID: addressID,
	}
	json.NewEncoder(w).Encode(response)
}

func (as *AddressService) saveAddress(address *models.Address) (int, error) {
	// Save the address to the 'addresses' table
	query := `INSERT INTO addresses (user_id, address_line1, address_line2, city, state, zip_code) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := as.DB.Exec(query, address.UserID, address.AddressLine1, address.AddressLine2, address.City, address.State, address.ZipCode)
	if err != nil {
		return 0, err
	}

	// Get the auto-generated address ID
	addressID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(addressID), nil
}

// CreateAddressResponse represents the response model for the CreateAddress endpoint
type CreateAddressResponse struct {
	AddressID int `json:"address_id"`
}

// GetAddressByID retrieves an address by ID
// @Summary Retrieve an address by ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Param id path string true "Address ID"
// @Success 200 {object} models.Address "Address retrieved successfully"
// @Failure 500 {object} ErrorResponse "Failed to retrieve address"
// @Router /addresses/{id} [get]
func (as *AddressService) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addressID := vars["id"]

	address, err := as.getAddressByID(addressID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve address", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

func (as *AddressService) getAddressByID(addressID string) (*models.Address, error) {
	// Get the address from the 'addresses' table
	query := `SELECT * FROM addresses WHERE address_id = ?`
	row := as.DB.QueryRow(query, addressID)

	address := &models.Address{}
	err := row.Scan(&address.AddressID, &address.UserID, &address.AddressLine1, &address.AddressLine2, &address.City, &address.State, &address.ZipCode)
	if err != nil {
		return nil, err
	}

	return address, nil
}

// UpdateAddress updates an address
// @Summary Update an address
// @Tags Addresses
// @Accept json
// @Produce json
// @Param id path string true "Address ID"
// @Param address body models.Address true "Address object"
// @Success 200 {string} string "Address updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to update address"
// @Router /addresses/{id} [put]
func (as *AddressService) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addressID := vars["id"]

	var address models.Address
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the address in the database
	err = as.updateAddress(addressID, &address)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update address", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Address updated successfully"))
}

func (as *AddressService) updateAddress(addressID string, address *models.Address) error {
	// Update the address in the 'addresses' table
	query := `UPDATE addresses SET address_line1 = ?, address_line2 = ?, city = ?, state = ?, zip_code = ? WHERE address_id = ?`
	_, err := as.DB.Exec(query, address.AddressLine1, address.AddressLine2, address.City, address.State, address.ZipCode, addressID)
	return err
}

// DeleteAddress deletes an address
// DeleteAddress deletes an address
// @Summary Delete an address
// @Tags Addresses
// @Param id path string true "Address ID"
// @Success 200 {string} string "Address deleted successfully"
// @Failure 500 {object} ErrorResponse "Failed to delete address"
// @Router /addresses/{id} [delete]
func (as *AddressService) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addressID := vars["id"]

	err := as.deleteAddress(addressID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete address", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Address deleted successfully"))
}

func (as *AddressService) deleteAddress(addressID string) error {
	// Delete the address from the 'addresses' table
	query := `DELETE FROM addresses WHERE address_id = ?`
	_, err := as.DB.Exec(query, addressID)
	return err
}

// GetAddressesByUserID fetches all addresses for a given user ID
// @Summary Fetch all addresses for a given user ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Address "Addresses retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 500 {object} ErrorResponse "Failed to get addresses"
// @Router /users/{user_id}/addresses [get]
func (as *AddressService) GetAddressesByUserID(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user_id from the request URL
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch the addresses from the database for the given user_id
	addresses, err := as.getAddressesByUserID(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get addresses", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addresses)
}

func (as *AddressService) getAddressesByUserID(userID int) ([]*models.Address, error) {
	// Query the 'addresses' table to get all addresses associated with the given user_id
	query := `SELECT address_id, address_line1, address_line2, city, state, zip_code FROM addresses WHERE user_id = ?`
	rows, err := as.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*models.Address
	for rows.Next() {
		address := &models.Address{}
		err := rows.Scan(&address.AddressID, &address.AddressLine1, &address.AddressLine2, &address.City, &address.State, &address.ZipCode)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// @Summary Check if an email ID exists
// @Tags Users
// @Accept json
// @Produce json
// @Param emailCheckRequest body EmailCheckRequest true "Request body containing email"
// @Success 200 {object} models.EmailCheckResponse "Email existence response"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Router /users/check-email [post]
func (as *AddressService) CheckEmailExists(w http.ResponseWriter, r *http.Request) {
	var req EmailCheckRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Error decoding request body: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := req.Email
	log.Println("Received email: ", email)

	if email == "" {
		http.Error(w, "Email field is required in the request body", http.StatusBadRequest)
		return
	}

	exists, err := as.emailExists(email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to check email", http.StatusInternalServerError)
		return
	}

	response := models.EmailCheckResponse{Exists: exists}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (as *AddressService) emailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	err := as.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		log.Println("Error querying the database: ", err)
		return false, err
	}
	return count > 0, nil
}
