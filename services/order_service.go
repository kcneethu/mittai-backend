package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gklps/mittai-backend/db"
	"github.com/gorilla/mux"
)

type OrderStatus struct {
	DB *db.Repository
}

func NewOrderStatusService(db *db.Repository) *OrderStatus {
	return &OrderStatus{
		DB: db,
	}
}

func (os *OrderStatus) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/orderstatus/{purchaseID}", os.UpdateOrderStatusHandler).Methods(http.MethodPut)
	r.HandleFunc("/orderstatus/{purchaseID}", os.GetOrderStatusHandler).Methods(http.MethodGet)
}

func (os *OrderStatus) UpdateOrderStatus(purchaseID int, status string) error {
	query := "INSERT OR REPLACE INTO orderstatus (purchase_id, status) VALUES (?, ?)"
	//print query with values
	log.Println(query, purchaseID, status)
	_, err := os.DB.Exec(query, status, purchaseID)
	return err
}

func (os *OrderStatus) GetOrderStatus(purchaseID int) (string, error) {
	query := "SELECT status FROM orderstatus WHERE purchase_id = ?"
	row := os.DB.QueryRow(query, purchaseID)

	var status string
	err := row.Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

// @Summary Update the status of an order by purchase ID
// @Tags OrderStatus
// @Accept json
// @Produce json
// @Param purchaseID path int true "Purchase ID"
// @Param status query string true "New order status"
// @Success 200 {string} string "Order status updated successfully"
// @Failure 400 "Bad request"
// @Failure 500 "Failed to update order status"
// @Router /orderstatus/{purchaseID} [put]
func (os *OrderStatus) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	purchaseIDStr := vars["purchaseID"]
	status := r.URL.Query().Get("status")

	purchaseID, err := strconv.Atoi(purchaseIDStr)
	if err != nil {
		http.Error(w, "Invalid purchase ID", http.StatusBadRequest)
		return
	}

	// Update the order status
	err = os.UpdateOrderStatus(purchaseID, status)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order status updated successfully"))
}

// @Summary Get the status of an order by purchase ID
// @Tags OrderStatus
// @Produce json
// @Param purchaseID path int true "Purchase ID"
// @Success 200 {string} string "Order status retrieved successfully"
// @Failure 400 "Bad request"
// @Failure 500 "Failed to fetch order status"
// @Router /orderstatus/{purchaseID} [get]
func (os *OrderStatus) GetOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	purchaseIDStr := vars["purchaseID"]

	purchaseID, err := strconv.Atoi(purchaseIDStr)
	if err != nil {
		http.Error(w, "Invalid purchase ID", http.StatusBadRequest)
		return
	}

	// Get the order status
	status, err := os.GetOrderStatus(purchaseID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch order status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}
