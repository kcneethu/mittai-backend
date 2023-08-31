package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gorilla/mux"
)

type ProductService struct {
	DB *db.Repository
}

func NewProductService(db *db.Repository) *ProductService {
	return &ProductService{
		DB: db,
	}
}

// AddProductResponse represents the response for the AddProduct method

type AddProductResponse struct {
	Message   string `json:"message"`
	ProductID int    `json:"product_id"`
}
type SuccessResponse struct {
	Message   string `json:"message"`
	ProductID int    `json:"product_id"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// AddProductRequest represents the request body for the AddProduct endpoint.
type AddProductRequest struct {
	Name            string   `json:"name" form:"name"`
	Description     string   `json:"description" form:"description"`
	Category        string   `json:"category" form:"category"`
	Ingredients     string   `json:"ingredients" form:"ingredients"`
	NutritionalInfo string   `json:"nutritional_info" form:"nutritional_info"`
	ImageURLs       []string `json:"image_urls" form:"image_urls"`
}

// RegisterRoutes registers the product management routes
func (ps *ProductService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/products", ps.ListProducts).Methods("GET")
	r.HandleFunc("/products/{id}", ps.GetProductDetails).Methods("GET")
	r.HandleFunc("/products", ps.AddProduct).Methods("POST")
	r.HandleFunc("/products/{id}", ps.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", ps.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/products/{productID}/weights/{weightID}", ps.UpdateProductWeightByID).Methods("PUT")
}

// AddProduct adds a new product to the inventory
// @Summary Add a new product to the inventory
// @Tags Products
// @Accept json
// @Param product body AddProductRequest true "Product details"
// @Produce json
// @Success 200 {object} AddProductResponse "Product added successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to add product"
// @Router /products [post]
func (ps *ProductService) AddProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var product AddProductRequest
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if image_urls is empty and assign an empty array if necessary
	if product.ImageURLs == nil {
		product.ImageURLs = []string{}
	}

	// Create a new Product instance
	newProduct := &models.Product{
		Name:            product.Name,
		Description:     product.Description,
		Category:        product.Category,
		Ingredients:     product.Ingredients,
		NutritionalInfo: product.NutritionalInfo,
		ImageURLs:       product.ImageURLs,
		Weights:         []*models.ProductWeight{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save the product to the database
	err = ps.saveProduct(newProduct)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to add product", http.StatusInternalServerError)
		return
	}

	// Send the response
	response := AddProductResponse{
		Message:   "Product added successfully",
		ProductID: newProduct.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// saveProduct saves a product to the database and assigns it a generated ID
func (ps *ProductService) saveProduct(product *models.Product) error {
	tx, err := ps.DB.Begin()
	if err != nil {
		return err
	}

	// Insert product into 'products' table
	query := `INSERT INTO products (name, description, category, ingredients, nutritional_info, image_urls, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, product.Name, product.Description, product.Category, product.Ingredients, product.NutritionalInfo, strings.Join(product.ImageURLs, ","), product.CreatedAt, product.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get the generated product ID
	productID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	product.ID = int(productID)

	// Insert product weights into 'product_weights' table
	for _, weight := range product.Weights {
		weight.ProductID = int(productID)
		query := `INSERT INTO product_weights (product_id, weight, price, stock, created_at, updated_at, measurement ) VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err := tx.Exec(query, weight.ProductID, weight.Weight, weight.Price, weight.StockAvailability, weight.CreatedAt, weight.UpdatedAt, weight.Measurement)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UpdateProduct updates an existing product in the inventory
// @Summary Update an existing product in the inventory
// @Tags Products
// @Accept multipart/form-data
// @Param id path string true "Product ID"
// @Param name formData string true "Product name"
// @Param description formData string true "Product description"
// @Param category formData string true "Product category"
// @Param ingredients formData string true "Product ingredients"
// @Param nutritional_info formData string true "Product nutritional information"
// @Param image_urls formData string true "Product image URLs (comma-separated)"
// @Success 200 "Product updated successfully"
// @Failure 400 "Invalid form data"
// @Failure 500 "Failed to update product"
// @Router /products/{id} [put]
func (ps *ProductService) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Retrieve the product ID from the path parameters
	vars := mux.Vars(r)
	productID := vars["id"]

	// Convert productID from string to int
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Parse the form data
	err = r.ParseMultipartForm(32 << 20) // 32MB limit for file upload
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Retrieve the product details from the form
	name := r.FormValue("name")
	description := r.FormValue("description")
	category := r.FormValue("category")
	ingredients := r.FormValue("ingredients")
	nutritionalInfo := r.FormValue("nutritional_info")
	imageURLs := r.FormValue("image_urls")

	// Split the comma-separated image URLs into an array
	imageURLsArray := strings.Split(imageURLs, ",")

	// Create a new Product instance
	product := &models.Product{
		ID:              productIDInt,
		Name:            name,
		Description:     description,
		Category:        category,
		Ingredients:     ingredients,
		NutritionalInfo: nutritionalInfo,
		ImageURLs:       imageURLsArray,
		Weights:         []*models.ProductWeight{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Update the product in the database
	err = ps.updateProduct(product)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

// updateProduct updates a product in the database
func (ps *ProductService) updateProduct(product *models.Product) error {
	tx, err := ps.DB.Begin()
	if err != nil {
		return err
	}

	// Update product in 'products' table
	query := `UPDATE products SET name = ?, description = ?, category = ?, ingredients = ?, nutritional_info = ?, image_urls = ?, updated_at = ? WHERE id = ?`
	_, err = tx.Exec(query, product.Name, product.Description, product.Category, product.Ingredients, product.NutritionalInfo, strings.Join(product.ImageURLs, ","), product.UpdatedAt, product.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert or update product weights in 'product_weights' table
	for _, weight := range product.Weights {
		weight.ProductID = product.ID

		// Check if the weight already exists
		query := `SELECT COUNT(*) FROM product_weights WHERE product_id = ? AND weight = ?`
		var count int
		err := tx.QueryRow(query, weight.ProductID, weight.Weight).Scan(&count)
		if err != nil {
			tx.Rollback()
			return err
		}

		if count > 0 {
			// Update the existing weight
			query = `UPDATE product_weights SET price = ?, stock = ?, updated_at = ? WHERE product_id = ? AND weight = ?`
			_, err = tx.Exec(query, weight.Price, weight.StockAvailability, weight.UpdatedAt, weight.ProductID, weight.Weight)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// Insert a new weight
			query = `INSERT INTO product_weights (product_id, weight, price, stock, created_at, updated_at, measurement) VALUES (?, ?, ?, ?, ?, ?)`
			_, err = tx.Exec(query, weight.ProductID, weight.Weight, weight.Price, weight.StockAvailability, weight.CreatedAt, weight.UpdatedAt, weight.Measurement)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetProductDetails retrieves a product from the inventory by its ID
// @Summary Get product details by ID
// @Tags Products
// @Param id path string true "Product ID"
// @Produce json
// @Success 200 {object} models.Product "Product details"
// @Failure 500 "Failed to retrieve product details"
// @Router /products/{id} [get]
func (ps *ProductService) GetProductDetails(w http.ResponseWriter, r *http.Request) {
	// Retrieve the product ID from the path parameters
	vars := mux.Vars(r)
	productID := vars["id"]

	// Get the product from the database
	product, err := ps.getProductByID(productID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve product details", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (ps *ProductService) GetProductByID(productID string) (*models.Product, error) {
	return ps.getProductByID(productID)
}

// getProductByID retrieves a product from the database by its ID
func (ps *ProductService) getProductByID(productID string) (*models.Product, error) {
	query := `SELECT id, name, description, category, ingredients, nutritional_info, image_urls, created_at, updated_at FROM products WHERE id = ?`
	row := ps.DB.QueryRow(query, productID)

	product := &models.Product{}
	var imageUrls string
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Ingredients, &product.NutritionalInfo, &imageUrls, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Split the comma-separated image URLs into an array
	imageURLs := strings.Split(imageUrls, ",")
	product.ImageURLs = make([]string, len(imageURLs))
	copy(product.ImageURLs, imageURLs)

	// Retrieve the product weights from the database
	weights, err := ps.getProductWeights(product.ID)
	if err != nil {
		return nil, err
	}
	product.Weights = weights

	return product, nil
}

// getProductWeights retrieves the product weights from the database
func (ps *ProductService) getProductWeights(productID int) ([]*models.ProductWeight, error) {
	query := `SELECT * FROM product_weights WHERE product_id = ?`
	rows, err := ps.DB.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weights []*models.ProductWeight
	for rows.Next() {
		weight := &models.ProductWeight{}
		err := rows.Scan(&weight.ID, &weight.ProductID, &weight.Weight, &weight.Price, &weight.Measurement, &weight.StockAvailability, &weight.CreatedAt, &weight.UpdatedAt)
		if err != nil {
			return nil, err
		}
		weights = append(weights, weight)
	}

	return weights, nil
}

// ListProducts returns a list of all products in the inventory
// @Summary List all products
// @Tags Products
// @Produce json
// @Success 200 {array} models.Product "List of products"
// @Failure 500 "Failed to retrieve products"
// @Router /products [get]
func (ps *ProductService) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Get the list of products from the database
	products, err := ps.getAllProducts()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}

	// Get the weights for each product
	for _, product := range products {
		weights, err := ps.getProductWeights(product.ID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to retrieve product weights", http.StatusInternalServerError)
			return
		}
		product.Weights = weights
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// getAllProducts retrieves all products from the database
func (ps *ProductService) getAllProducts() ([]*models.Product, error) {
	query := `SELECT id, name, description, category, ingredients, nutritional_info, image_urls, created_at, updated_at FROM products`
	rows, err := ps.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		var imageUrls string
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Ingredients, &product.NutritionalInfo, &imageUrls, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Split the comma-separated image URLs into an array
		product.ImageURLs = strings.Split(imageUrls, ",")

		products = append(products, product)
	}

	return products, nil
}

// DeleteProduct deletes a product from the inventory
// @Summary Delete a product from the inventory
// @Tags Products
// @Param id path string true "Product ID"
// @Success 200 "Product deleted successfully"
// @Failure 500 "Failed to delete product"
// @Router /products/{id} [delete]
func (ps *ProductService) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Retrieve the product ID from the path parameters
	vars := mux.Vars(r)
	productID := vars["id"]

	// Delete the product from the database
	err := ps.deleteProduct(productID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))
}

// deleteProduct deletes a product from the database
func (ps *ProductService) deleteProduct(productID string) error {
	tx, err := ps.DB.Begin()
	if err != nil {
		return err
	}

	// Delete product from 'products' table
	query := `DELETE FROM products WHERE id = ?`
	_, err = tx.Exec(query, productID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete product weights from 'product_weights' table
	query = `DELETE FROM product_weights WHERE product_id = ?`
	_, err = tx.Exec(query, productID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UpdateProductWeightByID updates the price of a product variant based on weight ID, removes weight information, and creates a new weight and price
// @Summary Update product price by weight ID, remove weight, and create new weight and price
// @Tags Products
// @Param productID path string true "Product ID"
// @Param weightID path string true "Weight ID"
// @Param price formData float64 true "Product price"
// @Success 200 "Product price updated successfully"
// @Failure 400 "Invalid form data"
// @Failure 500 "Failed to update product price"
// @Router /products/{productID}/weights/{weightID} [put]

func (ps *ProductService) UpdateProductWeightByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve the product ID and weight ID from the path parameters
	vars := mux.Vars(r)
	productID := vars["productID"]
	weightID := vars["weightID"]
	measurement := vars["measurement"]

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Retrieve the price from the form
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	// Update the product weight price in the database
	err = ps.updateProductWeightByID(productID, weightID, price, measurement)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update product price", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product price updated successfully"))
}

// updateProductWeightByID updates the price of a product variant based on weight ID, removes weight information, and creates a new weight and price
func (ps *ProductService) updateProductWeightByID(productID string, weightID string, price float64, measurement string) error {
	// Get the current weight information for the product
	weight, err := ps.getProductWeightByID(productID, weightID)
	if err != nil {
		return err
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return err
	}

	// Update the price of the current weight
	query := `UPDATE product_weights SET price = ? WHERE id = ?`
	_, err = tx.Exec(query, price, weight.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Remove the weight information from the product
	query = `DELETE FROM product_weights WHERE id = ?`
	_, err = tx.Exec(query, weight.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create a new weight and price for the product
	newWeight := &models.ProductWeight{
		ProductID:         weight.ProductID,
		Weight:            weight.Weight,
		Price:             price,
		StockAvailability: weight.StockAvailability,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Measurement:       weight.Measurement,
	}

	query = `INSERT INTO product_weights (product_id, weight, price, stock, created_at, updated_at, measurement) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(query, newWeight.ProductID, newWeight.Weight, newWeight.Price, newWeight.StockAvailability, newWeight.CreatedAt, newWeight.UpdatedAt, newWeight.Measurement)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// getProductWeightByID retrieves a product weight by ID
func (ps *ProductService) getProductWeightByID(productID string, weightID string) (*models.ProductWeight, error) {
	query := `SELECT * FROM product_weights WHERE product_id = ? AND id = ?`
	row := ps.DB.QueryRow(query, productID, weightID)

	weight := &models.ProductWeight{}
	err := row.Scan(&weight.ID, &weight.ProductID, &weight.Weight, &weight.Price, &weight.StockAvailability, &weight.CreatedAt, &weight.UpdatedAt, &weight.Measurement)
	if err != nil {
		return nil, err
	}

	return weight, nil
}

// GetProductPriceByID retrieves the price of a product variant based on its weight ID
func (ps *ProductService) GetProductPriceByID(weightID int) float64 {
	query := `SELECT price FROM product_weights WHERE id = ?`
	row := ps.DB.QueryRow(query, weightID)

	var price float64
	err := row.Scan(&price)
	if err != nil {
		log.Println(err)
		return 0
	}

	return price
}

// GetProductNameByID retrieves the name of a product based on its ID
func (ps *ProductService) GetProductNameByID(productID int) (string, error) {
	query := `SELECT name FROM products WHERE id = ?`
	row := ps.DB.QueryRow(query, productID)

	var name string
	err := row.Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}
