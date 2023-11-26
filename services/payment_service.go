package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// PaymentService handles the payment mode related operations
type PaymentService struct {
	DB *db.Repository
}

// NewPaymentService creates a new instance of PaymentService
func NewPaymentService(db *db.Repository) *PaymentService {
	return &PaymentService{
		DB: db,
	}
}

// RegisterRoutes registers the payment mode routes
func (ps *PaymentService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/payment/modes", ps.GetPaymentModes).Methods(http.MethodGet)
}

// GetPaymentModes retrieves all payment modes
func (ps *PaymentService) GetPaymentModes(w http.ResponseWriter, r *http.Request) {
	rows, err := ps.DB.Query("SELECT id, mode, is_active FROM payment_mode")
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch payment modes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var paymentModes []*models.PaymentMode

	for rows.Next() {
		var paymentMode models.PaymentMode

		err := rows.Scan(&paymentMode.ID, &paymentMode.Mode, &paymentMode.IsActive)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to fetch payment modes", http.StatusInternalServerError)
			return
		}

		paymentModes = append(paymentModes, &paymentMode)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paymentModes)
}
