package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gklps/mittai-backend/services"
	"github.com/go-chi/chi"
)

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

// ListProductsHandler retrieves a list of all products
// @Summary Retrieve all products
// @Description Retrieves a list of all products
// @Tags products
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := services.ListProducts()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get products: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

// listProducts retrieves all products from the database and writes the response
func listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := db.ListProducts()
	if err != nil {
		log.Printf("failed to get products: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(products)
	if err != nil {
		log.Printf("failed to marshal products response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// GetProductByIDHandler retrieves a product by its ID
// @Summary Retrieve a product by ID
// @Description Retrieves a product by its ID
// @Tags products
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{product_id} [get]
func GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := services.GetProductByID(productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get product: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// createProduct creates a new product in the database
func createProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Printf("failed to decode product data: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productID, err := db.CreateProduct(product)
	if err != nil {
		log.Printf("failed to create product: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]int{"productID": productID}
	jsonResponse(w, response)
}

// createProduct creates a new product in the database
// @Summary Create a new product
// @Description Creates a new product in the database
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product object"
// @Success 201 {object} map[string]int
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Printf("failed to decode product data: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productID, err := db.CreateProduct(product)
	if err != nil {
		log.Printf("failed to create product: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]int{"productID": productID}
	jsonResponse(w, response)
}
