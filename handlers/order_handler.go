package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gklps/mittai-backend/utils"
)

// OrdersHandler handles requests related to orders
func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ListOrdersHandler(w, r)
	case http.MethodPost:
		CreateOrderHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ListOrdersHandler retrieves a list of all orders
// @Summary Retrieve all orders
// @Description Retrieves a list of all orders
// @Tags orders
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} ErrorResponse
// @Router /orders [get]
func ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
	orders, err := db.ListOrders()
	if err != nil {
		log.Printf("failed to get orders: %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get orders")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, orders)
}

// CreateOrderHandler creates a new order in the database
// @Summary Create a new order
// @Description Creates a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Order object"
// @Success 201 {object} map[string]int
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [post]
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Printf("failed to decode order data: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "Failed to decode order data")
		return
	}

	// Validate the order data
	if err := validateOrderData(order); err != nil {
		log.Printf("invalid order data: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid order data")
		return
	}

	// Perform additional operations on the order data if needed

	orderID, err := db.CreateOrder(order)
	if err != nil {
		log.Printf("failed to create order: %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	response := map[string]int{"orderID": orderID}
	utils.SendJSONResponse(w, http.StatusCreated, response)
}

func validateOrderData(order models.Order) error {
	if len(order.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}
	return nil
}
