package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mittai-backend/db"
	"github.com/myapp/utils"
)

// ListProducts handler for fetching all products
func ListProducts(w http.ResponseWriter, r *http.Request) {
	// Fetch all products from the database
	products, err := db.ListProducts()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch products: %v", err))
		return
	}

	// Return the list of products as JSON response
	utils.SendJSONResponse(w, http.StatusOK, products)
}

// GetProductDetails handler for fetching details of a specific product
func GetProductDetails(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from request path parameters
	vars := mux.Vars(r)
	productID := vars["productID"]

	// Fetch the product details from the database using the ID
	product, err := db.GetProductByID(productID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch product details: %v", err))
		return
	}

	// Return the product details as JSON response
	utils.SendJSONResponse(w, http.StatusOK, product)
}

// Implement other product-related handlers (AddProduct, UpdateProduct, DeleteProduct) similarly
