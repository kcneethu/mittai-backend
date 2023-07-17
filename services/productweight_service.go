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

type ProductWeightService struct {
	DB *db.Repository
}

func NewProductWeightService(db *db.Repository) *ProductWeightService {
	return &ProductWeightService{
		DB: db,
	}
}

// AddProductWeightRequest represents the request body for adding a new weight variant
type AddProductWeightRequest struct {
	Weight            int     `json:"weight" form:"weight"`
	Price             float64 `json:"price" form:"price"`
	StockAvailability int     `json:"stock_availability" form:"stock_availability"`
}

// UpdateProductWeightRequest represents the request body for updating an existing weight variant
type UpdateProductWeightRequest struct {
	Weight            int     `json:"weight" form:"weight"`
	Price             float64 `json:"price" form:"price"`
	StockAvailability int     `json:"stock_availability" form:"stock_availability"`
}

// AddProductWeight adds a new weight variant for a product
// @Summary Add a new weight variant for a product
// @Tags Product Weights
// @Param productID path string true "Product ID"
// @Param weight body AddProductWeightRequest true "Product weight details"
// @Produce json
// @Success 200 {object} SuccessResponse "Weight added successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or product ID"
// @Failure 500 {object} ErrorResponse "Failed to add weight"
// @Router /productweight/{productID}/weights [post]
func (ps *ProductWeightService) AddProductWeight(w http.ResponseWriter, r *http.Request) {
	// Retrieve the product ID from the path parameters
	vars := mux.Vars(r)
	productIDStr := vars["productID"]

	// Convert productID from string to int
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var weight AddProductWeightRequest
	err = json.NewDecoder(r.Body).Decode(&weight)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create a new ProductWeight instance
	newWeight := &models.ProductWeight{
		ProductID:         productID,
		Weight:            weight.Weight,
		Price:             weight.Price,
		StockAvailability: weight.StockAvailability,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Save the weight to the database
	err = ps.saveProductWeight(newWeight)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to add weight", http.StatusInternalServerError)
		return
	}

	// Send the response
	response := SuccessResponse{
		Message:   "Weight added successfully",
		ProductID: newWeight.ProductID,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateProductWeight updates an existing weight variant for a product
// @Summary Update an existing weight variant for a product
// @Tags Product Weights
// @Param productID path string true "Product ID"
// @Param weightID path string true "Weight ID"
// @Param weight body UpdateProductWeightRequest true "Product weight details"
// @Success 200 "Weight updated successfully"
// @Failure 400 "Invalid request body or product/weight ID"
// @Failure 500 "Failed to update weight"
// @Router /productweight/{productID}/weights/{weightID} [put]
func (ps *ProductWeightService) UpdateProductWeight(w http.ResponseWriter, r *http.Request) {
	// Retrieve the product ID and weight ID from the path parameters
	vars := mux.Vars(r)
	productID := vars["productID"]
	weightID := vars["weightID"]

	// Parse the request body
	var weight UpdateProductWeightRequest
	err := json.NewDecoder(r.Body).Decode(&weight)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the product weight in the database
	err = ps.updateProductWeight(productID, weightID, weight)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update weight", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Weight updated successfully"))
}

// saveProductWeight saves a weight variant for a product to the database
func (ps *ProductWeightService) saveProductWeight(weight *models.ProductWeight) error {
	query := `INSERT INTO product_weights (product_id, weight, price, stock, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := ps.DB.Exec(query, weight.ProductID, weight.Weight, weight.Price, weight.StockAvailability, weight.CreatedAt, weight.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// updateProductWeight updates an existing weight variant for a product in the database
func (ps *ProductWeightService) updateProductWeight(productID string, weightID string, weight UpdateProductWeightRequest) error {
	query := `UPDATE product_weights SET weight = ?, price = ?, stock = ?, updated_at = ? WHERE product_id = ? AND id = ?`
	_, err := ps.DB.Exec(query, weight.Weight, weight.Price, weight.StockAvailability, time.Now(), productID, weightID)
	if err != nil {
		return err
	}
	return nil
}

// DefineRoutes sets up the routes for the ProductWeightService
func (ps *ProductWeightService) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/productweight/{productID}/weights", ps.AddProductWeight).Methods("POST")
	router.HandleFunc("/productweight/{productID}/weights/{weightID}", ps.UpdateProductWeight).Methods("PUT")
}
