package services

import (
	"encoding/json"
	"fmt"
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
	DB             *db.Repository
	ProductService *ProductService
}

// NewPurchaseService creates a new instance of PurchaseService
func NewPurchaseService(db *db.Repository, prodService *ProductService) *PurchaseService {
	return &PurchaseService{
		DB:             db,
		ProductService: prodService,
	}
}

// RegisterRoutes registers the purchase routes
func (ps *PurchaseService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/purchase", ps.CreatePurchase).Methods(http.MethodPost)
	r.HandleFunc("/purchase/{userID}", ps.GetPurchasesByUserID).Methods(http.MethodGet)
}

// @Summary Create a new purchase
// @Tags Purchases
// @Accept json
// @Produce json
// @Param purchase body models.PurchaseRequest true "Purchase payload"
// @Success 200 {string} string "Purchase created successfully"
// @Failure 400 "Bad request"
// @Failure 500 "Failed to create purchase"
// @Router /purchase [post]
func (ps *PurchaseService) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into the struct
	var purchase models.CreatePurchase
	err := json.NewDecoder(r.Body).Decode(&purchase)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Store the purchase in the database (this is a simplified example and may require more detailed DB operations)
	err = ps.storePurchaseInDB(purchase)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create purchase", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Purchase created successfully"))
}

func (ps *PurchaseService) storePurchaseInDB(purchase models.CreatePurchase) error {
	// Begin a transaction
	tx, err := ps.DB.Begin()
	if err != nil {
		return err
	}

	// Insert into purchase table
	query := `INSERT INTO purchases (address_id, payment_id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, purchase.AddressID, purchase.PaymentID, purchase.UserID, time.Now(), time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get the last inserted ID of the purchase
	purchaseID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert purchase items
	for _, item := range purchase.PurchaseItems {
		product, err := ps.ProductService.GetProductByID(strconv.Itoa(item.ProductID))
		if err != nil {
			tx.Rollback()
			return err
		}

		var weight *models.ProductWeight
		for _, w := range product.Weights {
			if w.ID == item.ProductWeightID {
				weight = w
				break
			}
		}

		if weight == nil {
			tx.Rollback()
			return fmt.Errorf("Weight not found for Product ID: %d, Weight ID: %d", item.ProductID, item.ProductWeightID)
		}

		itemTotalPrice := weight.Price * float64(item.Quantity)
		query := `INSERT INTO purchase_items (purchase_id, product_id, product_name, product_weight_id, product_price, quantity, total_price) VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(query, purchaseID, item.ProductID, product.Name, item.ProductWeightID, weight.Price, item.Quantity, itemTotalPrice)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetPurchasesByUserID retrieves purchases made by a specific user
// @Summary Get purchases by user ID
// @Tags Purchases
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
