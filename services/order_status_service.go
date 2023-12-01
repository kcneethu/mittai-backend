package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gklps/mittai-backend/db"
	"github.com/gorilla/mux"
)

type OrderStatusService struct {
	DB *db.Repository
}

// StatusResponse represents the response structure for retrieving order status.
type StatusResponse struct {
	Status string `json:"status"`
}

func NewOrderStatusService(db *db.Repository) *OrderStatusService {
	return &OrderStatusService{
		DB: db,
	}
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
// @Router /purchase/{purchaseID}/status [put]
func (os *OrderStatusService) UpdateOrderStatus(purchaseID int, status string) error {
	query := "UPDATE orderstatus SET status = ? WHERE purchase_id = ?"
	_, err := os.DB.Exec(query, status, purchaseID)
	return err
}

// @Summary Get the status of an order by purchase ID
// @Tags OrderStatus
// @Produce json
// @Param purchaseID path int true "Purchase ID"
// @Success 200 {object} StatusResponse "Order status retrieved successfully"
// @Failure 400 "Bad request"
// @Failure 500 "Failed to fetch order status"
// @Router /purchase/{purchaseID}/status [get]
func (ps *PurchaseService) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	purchaseIDStr := vars["purchaseID"]

	purchaseID, err := strconv.Atoi(purchaseIDStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid purchase ID", http.StatusBadRequest)
		return
	}

	// Retrieve the order status from the database
	status, err := ps.getOrderStatusFromDB(purchaseID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch order status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (ps *PurchaseService) getOrderStatusFromDB(purchaseID int) (string, error) {
	query := "SELECT status FROM orderstatus WHERE purchase_id = ?"
	row := ps.DB.QueryRow(query, purchaseID)

	var status string
	err := row.Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}
