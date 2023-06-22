package handlers

/*

import (
	"encoding/json"
	"net/http"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
)

// UsersHandler handles requests related to users
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listUsers(w, r)
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ProductsHandler handles requests related to products
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listProducts(w, r)
	case http.MethodPost:
		createProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// OrdersHandler handles requests related to orders
func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listOrders(w, r)
	case http.MethodPost:
		createOrder(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listUsers retrieves all users from the database and writes the response
func listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, users)
}

// listProducts retrieves all products from the database and writes the response
func listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := db.ListProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, products)
}

// createProduct creates a new product in the database and writes the response
func createProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	productID, err := db.CreateProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int{"productID": productID}
	jsonResponse(w, response)
}

// listOrders retrieves all orders from the database and writes the response
func listOrders(w http.ResponseWriter, r *http.Request) {
	// Retrieve orders from the database
	orders, err := db.ListOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, orders)
}

// createOrder creates a new order in the database and writes the response
func createOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	orderID, err := db.CreateOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int{"orderID": orderID}
	jsonResponse(w, response)
}

// jsonResponse writes a JSON response with the specified data
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
*/
