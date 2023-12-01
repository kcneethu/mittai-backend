package services

import (
	"encoding/json"
	"fmt"
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

func (os *OrderStatusService) UpdateOrderStatus(purchaseID int, status string) error {

	query := "UPDATE orderstatus SET status = ? WHERE purchase_id = ?"
	//print the entire query
	log.Println(query, status, purchaseID)
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

	fmt.Println("purchaseIDStr", purchaseIDStr)
	fmt.Println("purchaseID", purchaseID)
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
