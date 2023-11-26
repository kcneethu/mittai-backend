package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (us *UserService) saveOTP(userID int, otp string) error {
	query := `INSERT INTO otp (user_id, otp_value) VALUES (?, ?)`
	_, err := us.DB.Exec(query, userID, otp)
	return err
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (us *UserService) sendOTPEmail(email, otp string, userid int) error {
	client := &http.Client{}

	// Construct the verification link with user ID and OTP
	verifyLink := fmt.Sprintf("http://swagger.mittaitheruvu.com/verify-otp/%d", userid)
	// Construct the HTML email content
	htmlContent := fmt.Sprintf(
		"<!DOCTYPE html><html><head><title>OTP Email</title></head><body>"+
			"<h1>Your OTP is: %s</h1>"+
			"<p>Click <a href=\"%s\">here</a> to confirm your account.</p>"+
			"</body></html>", otp, verifyLink)

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  "Mittaitheruvu",
			"email": "otp-services@mittaitheruvu.com",
		},
		"params": map[string]string{
			"OTP": otp,
		},
		"to": []map[string]string{
			{
				"email": email,
			},
		},
		"subject":     "OTP Email",
		"htmlContent": htmlContent,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}

	req.Header.Add("api-key", "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")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var emailResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&emailResponse)
	if err != nil {
		return err
	}

	log.Printf("Email response: %v", emailResponse)

	return nil
}

// VerifyOTP verifies the provided OTP and updates verified_account
// @Summary Verify OTP and update verified_account
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param otp body models.VerifyOTPRequest true "OTP object"
// @Success 200 {string} string "OTP verified and account updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {string} string "Invalid OTP"
// @Failure 500 {object} ErrorResponse "Failed to update user"
// @Router /verify-otp/{id} [post]
func (us *UserService) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var input struct {
		OTP string `json:"otp"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if us.verifyOTP(userID, input.OTP) {
		// Update verified_account to true
		err := us.updateVerifiedAccount(userID, true)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (us *UserService) verifyOTP(userID, otp string) bool {
	query := `SELECT otp_value FROM otp WHERE user_id = ? AND otp_value = ? AND generated_at >= ?`
	row := us.DB.QueryRow(query, userID, otp, time.Now().Add(-12*time.Hour))

	var storedOTP string
	err := row.Scan(&storedOTP)
	if err != nil {
		return false
	}

	return storedOTP == otp
}

func (us *UserService) updateVerifiedAccount(userID string, verified bool) error {
	query := `UPDATE users SET verified_account = ? WHERE user_id = ?`
	_, err := us.DB.Exec(query, verified, userID)
	return err
}
