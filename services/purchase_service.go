package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gorilla/mux"
)

// PurchaseService handles the purchase related operations
type PurchaseService struct {
	DB              *db.Repository
	ProductService  *ProductService
	RecentPurchases map[string]time.Time
	CartService     *CartService
	Mutex           sync.Mutex
	OrderStatus     *OrderStatus // Add OrderStatusService for handling order status updates
}

// NewPurchaseService creates a new instance of PurchaseService
func NewPurchaseService(db *db.Repository, prodService *ProductService, cartService *CartService) *PurchaseService {
	return &PurchaseService{
		DB:              db,
		ProductService:  prodService,
		RecentPurchases: make(map[string]time.Time),
		CartService:     cartService,
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

	// Convert the purchase request to a string to use as a key for the map
	purchaseKey, err := json.Marshal(purchase)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	ps.Mutex.Lock()

	// Check if the same request has been processed recently
	if timestamp, exists := ps.RecentPurchases[string(purchaseKey)]; exists {
		if time.Since(timestamp) <= time.Duration(1.5*float64(time.Second)) {
			http.Error(w, "Duplicate request", http.StatusTooManyRequests)
			ps.Mutex.Unlock() // Don't forget to unlock before returning
			return
		}
	}

	// Store the purchase in the database
	purchaseID, err := ps.storePurchaseInDB(purchase)
	fmt.Println("latest purchase id", purchaseID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create purchase", http.StatusInternalServerError)
		return
	}

	// Insert 'accepted' status in the orderstatus table using the retrieved purchaseID
	err = ps.OrderStatus.AddOrderStatus(int(purchaseID))
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	// Update the map with the current request and timestamp
	ps.RecentPurchases[string(purchaseKey)] = time.Now()
	ps.Mutex.Unlock()

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Purchase created successfully"))
}

func (ps *PurchaseService) storePurchaseInDB(purchase models.CreatePurchase) (int64, error) {
	// Begin a transaction
	tx, err := ps.DB.Begin()
	if err != nil {
		return 0, err
	}

	// Insert into purchase table
	query := `INSERT INTO purchases (address_id, payment_id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, purchase.AddressID, purchase.PaymentID, purchase.UserID, time.Now(), time.Now())
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get the last inserted ID of the purchase
	purchaseID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert purchase items
	for _, item := range purchase.PurchaseItems {
		product, err := ps.ProductService.GetProductByID(strconv.Itoa(item.ProductID))
		if err != nil {
			tx.Rollback()
			return 0, err
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
			return 0, fmt.Errorf("Weight not found for Product ID: %d, Weight ID: %d", item.ProductID, item.ProductWeightID)
		}

		itemTotalPrice := weight.Price * float64(item.Quantity)
		query := `INSERT INTO purchase_items (purchase_id, product_id, product_name, product_weight_id, product_price, quantity, total_price) VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(query, purchaseID, item.ProductID, product.Name, item.ProductWeightID, weight.Price, item.Quantity, itemTotalPrice)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	_, err = tx.Exec("DELETE FROM cart WHERE user_id = ?", purchase.UserID)

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return purchaseID, nil
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

	rows, err := ps.DB.Query("SELECT id, address_id, payment_id, created_at, updated_at FROM purchases WHERE user_id = ? ORDER BY id DESC", userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch purchases", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var purchases []*models.Purchase

	for rows.Next() {
		var purchase models.Purchase

		err := rows.Scan(&purchase.ID, &purchase.AddressID, &purchase.PaymentID, &purchase.CreatedAt, &purchase.UpdatedAt)
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
	// Adjusting the SQL query to include a JOIN with product_weights and fetch weight.
	query := `SELECT 
                pi.product_id, 
                pi.product_name, 
                pi.product_price, 
                pw.weight, 
                pi.quantity, 
                pi.total_price,
				pi.product_weight_id,
				pw.measurement
              FROM 
                purchase_items AS pi 
              JOIN 
                product_weights AS pw ON pi.product_weight_id = pw.id 
              WHERE 
                pi.purchase_id = ?
				ORDER BY 
    			pi.purchase_id DESC;`

	rows, err := ps.DB.Query(query, purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.PurchaseItem

	for rows.Next() {
		var item models.PurchaseItem

		// Scan now includes weight from product_weights table.
		err := rows.Scan(&item.ProductID, &item.ProductName, &item.ProductPrice, &item.Weight, &item.Quantity, &item.TotalPrice, &item.ProductWeightID, &item.Measurement)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}
