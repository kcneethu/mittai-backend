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

// CreatePurchase creates a new purchase
// @Summary Create a purchase
// @Tags Purchase
// @Accept json
// @Produce json
// @Param request body models.PurchaseRequest true "Purchase request payload"
// @Success 200 {object} models.PurchaseResponse "Purchase created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 409 {object} ErrorResponse "Insufficient stock"
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

	if request.UserID <= 0 || request.AddressID <= 0 || request.PaymentID <= 0 {
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

	// Check product availability and update stock
	for _, item := range request.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 || item.ProductWeightID <= 0 {
			tx.Rollback()
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Check if the product exists
		product, err := ps.ProductService.GetProductByID(strconv.Itoa(item.ProductID))
		if err != nil {
			log.Println(err)
			tx.Rollback()
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// Find the specific weight variant of the product
		var weight *models.ProductWeight
		for _, w := range product.Weights {
			if w.ID == item.ProductWeightID {
				weight = w
				break
			}
		}

		if weight == nil {
			tx.Rollback()
			http.Error(w, "Invalid product weight ID", http.StatusBadRequest)
			return
		}

		// Check if there is enough stock
		if item.Quantity > weight.StockAvailability {
			tx.Rollback()
			http.Error(w, "Insufficient stock", http.StatusConflict)
			return
		}

		// Reduce the stock quantity
		updatedStock := weight.StockAvailability - item.Quantity
		_, err = tx.Exec("UPDATE product_weights SET stock_availability = ? WHERE id = ?", updatedStock, item.ProductWeightID)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			http.Error(w, "Failed to update stock", http.StatusInternalServerError)
			return
		}

		// Insert purchase items into the database
		_, err = tx.Exec("INSERT INTO purchase_items (purchase_id, product_id, product_name, product_price, quantity, total_price, product_weight_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
			purchaseID, item.ProductID, item.ProductName, item.ProductPrice, item.Quantity, item.TotalPrice, item.ProductWeightID)
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
