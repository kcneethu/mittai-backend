package services

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gklps/mittai-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	UserID int `json:"user_id"`
	// Add other necessary fields if needed.
}

// Login logs in a user and returns the user_id if the password is correct
// @Summary User login
// @Tags Users
// @Accept json
// @Produce json
// @Param loginReq body LoginRequest true "Login request object"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "Invalid email or password"
// @Router /login [post]
// Login logs in a user and returns the user_id if the password is correct
func (us *UserService) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Retrieve the user by email from the database
	user, err := us.getUserByEmail(loginReq.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Hash the provided password using the same method as during user registration
	hashedPassword, err := us.hashPassword(loginReq.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Log the hashed passwords for debugging
	log.Println("Stored Hashed Password:", user.Password)
	log.Println("Input Hashed Password:", hashedPassword)

	// Compare the hashed password with the hashed password stored in the database
	if user.Password != hashedPassword {
		log.Println("Invalid password")
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// If the password is correct, respond with the user_id
	response := LoginResponse{
		UserID: user.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getUserByEmail retrieves a user by email from the database
func (us *UserService) getUserByEmail(email string) (*models.User, error) {
	query := `SELECT user_id, first_name, last_name, email, contact_number, verified_account, hashed_password FROM users WHERE email = ?`
	row := us.DB.QueryRow(query, email)

	user := &models.User{}
	err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.ContactNumber, &user.VerifiedAccount, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// comparePasswords compares the provided password with the hashed password
func (us *UserService) comparePasswords(hashedPassword, providedPassword string) error {
	// Compare the provided password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}
