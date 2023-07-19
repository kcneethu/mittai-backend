package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gorilla/mux"
)

// PurchaseService handles the purchase related operations
type PurchaseService struct {
	DB *db.Repository
}

// NewPurchaseService creates a new instance of PurchaseService
func NewPurchaseService(db *db.Repository) *PurchaseService {
	return &PurchaseService{
		DB: db,
	}
}

// RegisterRoutes registers the purchase routes
func (ps *PurchaseService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/purchase", ps.CreatePurchase).Methods(http.MethodPost)
	r.HandleFunc("/purchase/{userID}", ps.GetPurchasesByUserID).Methods(http.MethodGet)
}

// CreatePurchase creates a new purchase
// @Summary Create a purchase
// @Tags Purchase
// @Accept json
// @Produce json
// @Param request body models.PurchaseRequest true "Purchase request payload"
// @Success 200 {object} models.PurchaseResponse "Purchase created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 500 {object} ErrorResponse "Failed to create purchase"
// @Router /purchase [post]
func (ps *PurchaseService) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	var request models.PurchaseRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create purchase", http.StatusInternalServerError)
		return
	}

	// Calculate total price
	totalPrice := 0.0
	for _, item := range request.Items {
		totalPrice += item.TotalPrice
	}

	// Insert the purchase into the database
	result, err := tx.Exec("INSERT INTO purchases (user_id, total_price, address_id, payment_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		request.UserID, totalPrice, request.AddressID, request.PaymentID, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		tx.Rollback()
		http.Error(w, "Failed to create purchase", http.StatusInternalServerError)
		return
	}

	// Get the ID of the created purchase
	purchaseID, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		http.Error(w, "Failed to create purchase", http.StatusInternalServerError)
		return
	}

	// Insert purchase items into the database
	for _, item := range request.Items {
		_, err := tx.Exec("INSERT INTO purchase_items (purchase_id, product_id, product_name, product_price, quantity, total_price) VALUES (?, ?, ?, ?, ?, ?)",
			purchaseID, item.ProductID, item.ProductName, item.ProductPrice, item.Quantity, item.TotalPrice)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			http.Error(w, "Failed to create purchase", http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create purchase", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.PurchaseResponse{
		PurchaseID: int(purchaseID),
	})
}

// GetPurchasesByUserID retrieves purchases made by a specific user
// @Summary Get purchases by user ID
// @Tags Purchase
// @Param userID path int true "User ID"
// @Produce json
// @Success 200 {array} models.Purchase "Purchases retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 500 {object} ErrorResponse "Failed to fetch purchases"
// @Router /purchase/{userID} [get]
func (ps *PurchaseService) GetPurchasesByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["userID"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	rows, err := ps.DB.Query("SELECT id, total_price, address_id, payment_id, created_at, updated_at FROM purchases WHERE user_id = ?", userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch purchases", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var purchases []*models.Purchase

	for rows.Next() {
		var purchase models.Purchase

		err := rows.Scan(&purchase.ID, &purchase.TotalPrice, &purchase.AddressID, &purchase.PaymentID, &purchase.CreatedAt, &purchase.UpdatedAt)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to fetch purchases", http.StatusInternalServerError)
			return
		}

		// Retrieve purchase items for each purchase
		items, err := ps.getPurchaseItemsByPurchaseID(purchase.ID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to fetch purchases", http.StatusInternalServerError)
			return
		}

		purchase.Items = items
		purchases = append(purchases, &purchase)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(purchases)
}

// getPurchaseItemsByPurchaseID retrieves purchase items for a specific purchase ID
func (ps *PurchaseService) getPurchaseItemsByPurchaseID(purchaseID int) ([]*models.PurchaseItem, error) {
	rows, err := ps.DB.Query("SELECT product_id, product_name, product_price, quantity, total_price FROM purchase_items WHERE purchase_id = ?", purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.PurchaseItem

	for rows.Next() {
		var item models.PurchaseItem

		err := rows.Scan(&item.ProductID, &item.ProductName, &item.ProductPrice, &item.Quantity, &item.TotalPrice)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}
