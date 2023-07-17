package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gorilla/mux"
)

// CartService represents a service for cart-related operations
type CartService struct {
	Repository *db.Repository
	CartDB     map[int]*models.Cart
}

// NewCartService creates a new CartService instance
func NewCartService(repository *db.Repository) *CartService {
	return &CartService{
		Repository: repository,
		CartDB:     make(map[int]*models.Cart),
	}
}

// RegisterRoutes registers the cart service routes
func (cs *CartService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/cart/{userID}/{productWeightID}", cs.AddToCart).Methods("POST")
	r.HandleFunc("/cart/{userID}", cs.GetCartByUserID).Methods("GET")
	r.HandleFunc("/cart/{userID}/{productWeightID}", cs.UpdateCartItem).Methods("PUT")
	r.HandleFunc("/cart/{userID}/{productWeightID}", cs.RemoveCartItem).Methods("DELETE")
	r.HandleFunc("/cart/{userID}/clear", cs.ClearCart).Methods("POST")
}

// AddToCart adds an item to the user's cart
// @Summary Add an item to the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param productWeightID path string true "Product Weight ID"
// @Success 200 {string} string "Item added to cart successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID or product weight ID"
// @Failure 404 {object} ErrorResponse "Cart not found or Product weight not found"
// @Router /cart/{userID}/{productWeightID} [post]
func (cs *CartService) AddToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["userID"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	productWeightIDStr := vars["productWeightID"]
	productWeightID, err := strconv.Atoi(productWeightIDStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid product weight ID", http.StatusBadRequest)
		return
	}

	productWeight, err := cs.getProductWeightByID(productWeightID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch product weight details", http.StatusInternalServerError)
		return
	}

	newCartItem := &models.CartItem{
		Product:   productWeight,
		Quantity:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cart, exists := cs.CartDB[userID]
	if !exists {
		cart, err = cs.getCartByUserID(userID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to fetch cart", http.StatusInternalServerError)
			return
		}
		if cart == nil {
			cart = &models.Cart{
				UserID:    userID,
				Items:     make([]*models.CartItem, 0),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
		cs.CartDB[userID] = cart
	}

	for _, item := range cart.Items {
		if item.Product != nil && item.Product.ID == newCartItem.Product.ID {
			item.Quantity += newCartItem.Quantity

			err = cs.updateCart(cart)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to update cart", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Item added to cart successfully"))
			return
		}
	}

	cart.Items = append(cart.Items, newCartItem)

	cart.TotalPrice = cart.GetTotalPrice()

	stmt, err := cs.Repository.DB.Prepare("INSERT INTO cart_items (cart_id, product_weight_id, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(cart.ID, newCartItem.Product.ID, newCartItem.Quantity, newCartItem.CreatedAt, newCartItem.UpdatedAt)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to insert cart item", http.StatusInternalServerError)
		return
	}

	err = cs.updateCart(cart)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item added to cart successfully"))
}

// GetCartByUserID retrieves the user's cart by user ID
// @Summary Get the user's cart by user ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} GetCartResponse "User's cart retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Router /cart/{userID} [get]
func (cs *CartService) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["userID"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cart, exists := cs.CartDB[userID]
	if !exists {
		cart, err = cs.getCartByUserID(userID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to fetch cart", http.StatusInternalServerError)
			return
		}
		if cart == nil {
			cart = &models.Cart{
				UserID:    userID,
				Items:     make([]*models.CartItem, 0),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
		cs.CartDB[userID] = cart
	}

	totalPrice := cart.GetTotalPrice()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetCartResponse{
		Cart:       cart,
		TotalPrice: totalPrice,
	})
}

// UpdateCartItem updates the quantity of a cart item in the user's cart
// @Summary Update the quantity of a cart item
// @Tags Cart
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param productWeightID path string true "Product Weight ID"
// @Param quantity query int true "Quantity"
// @Success 200 {string} string "Cart item updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID or product weight ID or quantity"
// @Failure 404 {object} ErrorResponse "Cart item not found"
// @Router /cart/{userID}/{productWeightID} [put]
func (cs *CartService) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productWeightIDStr := vars["productWeightID"]
	userIDstr := vars["userID"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get the user's cart from the database
	cart, exists := cs.CartDB[userID]
	if !exists {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}
	cart.TotalPrice = cart.GetTotalPrice()

	productWeightID, err := strconv.Atoi(productWeightIDStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid product weight ID", http.StatusBadRequest)
		return
	}

	// Find the cart item by matching the product's ID
	for _, item := range cart.Items {
		if item.Product != nil && item.Product.ID == productWeightID {
			// Update the cart item quantity
			quantityStr := r.URL.Query().Get("quantity")
			quantity, err := strconv.Atoi(quantityStr)
			if err != nil {
				log.Println(err)
				http.Error(w, "Invalid quantity", http.StatusBadRequest)
				return
			}

			item.Quantity = quantity

			// Update the cart item quantity in the database
			err = cs.updateCartItemQuantity(userID, productWeightID, quantity)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to update cart item quantity", http.StatusInternalServerError)
				return
			}

			// Update the cart in the database
			err = cs.updateCart(cart)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to update cart", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Cart item updated successfully"))
			return
		}
	}

	http.Error(w, "Cart item not found", http.StatusNotFound)
}

// RemoveCartItem removes a cart item from the user's cart
// @Summary Remove a cart item from the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param productWeightID path string true "Product Weight ID"
// @Success 200 {string} string "Cart item removed successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID or product weight ID"
// @Failure 404 {object} ErrorResponse "Cart not found or Cart item not found"
// @Router /cart/{userID}/{productWeightID} [delete]
func (cs *CartService) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productWeightIDStr := vars["productWeightID"]
	userIDstr := vars["userID"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get the user's cart from the database
	cart, exists := cs.CartDB[userID]
	if !exists {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	productWeightID, err := strconv.Atoi(productWeightIDStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid product weight ID", http.StatusBadRequest)
		return
	}

	// Find the index of the cart item by matching the product's ID
	index := -1
	for i, item := range cart.Items {
		if item.Product.ID == productWeightID {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Cart item not found", http.StatusNotFound)
		return
	}
	cart.TotalPrice = cart.GetTotalPrice()

	// Remove the cart item from the cart
	cart.Items = append(cart.Items[:index], cart.Items[index+1:]...)

	// Update the cart in the database
	err = cs.updateCart(cart)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cart item removed successfully"))
}

// ClearCart clears the cart for a specific user
// @Summary Clear the cart for a specific user
// @Tags Cart
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {string} string "Cart cleared successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 404 {object} ErrorResponse "Cart not found"
// @Failure 500 {object} ErrorResponse "Failed to clear cart"
// @Router /cart/{userID}/clear [post]
func (cs *CartService) ClearCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["userID"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cart, exists := cs.CartDB[userID]
	if !exists {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	cart.Items = make([]*models.CartItem, 0)
	cart.TotalPrice = 0

	err = cs.updateCart(cart)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cart cleared successfully"))
}

// getProductWeightByID fetches the product weight details from the database based on the given productWeightID
func (cs *CartService) getProductWeightByID(productWeightID int) (*models.ProductWeight, error) {
	query := "SELECT id, product_id, weight, price, stock_availability, created_at, updated_at FROM product_weights WHERE id = ?"

	// Prepare the query
	stmt, err := cs.Repository.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	row := stmt.QueryRow(productWeightID)

	// Create variables to store the retrieved data
	var (
		id                int
		productID         int
		weight            int
		price             float64
		stockAvailability int
		createdAt         time.Time
		updatedAt         time.Time
	)

	// Scan the row data into the variables
	err = row.Scan(&id, &productID, &weight, &price, &stockAvailability, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	// Create and return the product weight instance
	productWeight := &models.ProductWeight{
		ID:                id,
		ProductID:         productID,
		Weight:            weight,
		Price:             price,
		StockAvailability: stockAvailability,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	return productWeight, nil
}

// updateCart updates the cart in the database
func (cs *CartService) updateCart(cart *models.Cart) error {
	var err error
	if cart.ID == 0 {
		// Create a new cart in the database
		stmt, err := cs.Repository.DB.Prepare("INSERT INTO carts (user_id, created_at, updated_at) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(cart.UserID, cart.CreatedAt, time.Now())
		if err != nil {
			return err
		}

		// Get the generated cart ID
		cartID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		cart.ID = int(cartID)
	} else {
		// Update the existing cart in the database
		stmt, err := cs.Repository.DB.Prepare("UPDATE carts SET updated_at = ? WHERE id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(time.Now(), cart.ID)
		if err != nil {
			return err
		}
	}

	// Delete the existing cart items for the user from the database
	_, err = cs.Repository.DB.Exec("DELETE FROM cart_items WHERE cart_id = ?", cart.ID)
	if err != nil {
		return err
	}

	// Insert the new cart items into the database
	stmt, err := cs.Repository.DB.Prepare("INSERT INTO cart_items (cart_id, product_weight_id, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range cart.Items {
		_, err := stmt.Exec(cart.ID, item.Product.ID, item.Quantity, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateCartItemQuantity updates the quantity of a cart item in the database
func (cs *CartService) updateCartItemQuantity(userID, productWeightID, quantity int) error {
	// Prepare the SQL statement
	stmt, err := cs.Repository.DB.Prepare("UPDATE cart_items SET quantity = ? WHERE cart_id = ? AND product_weight_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(quantity, userID, productWeightID)
	if err != nil {
		return err
	}

	return nil
}

// getCartByUserID retrieves the user's cart from the database
func (cs *CartService) getCartByUserID(userID int) (*models.Cart, error) {
	query := "SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = ?"

	// Prepare the query
	stmt, err := cs.Repository.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	row := stmt.QueryRow(userID)

	// Create variables to store the retrieved data
	var (
		id        int
		createdAt time.Time
		updatedAt time.Time
	)

	// Scan the row into the variables
	err = row.Scan(&id, &userID, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if no cart found
		}
		return nil, err
	}

	// Create a new cart instance
	cart := &models.Cart{
		ID:        id,
		UserID:    userID,
		Items:     make([]*models.CartItem, 0),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	// Retrieve the cart items from the database
	rows, err := cs.Repository.DB.Query("SELECT product_weight_id, quantity FROM cart_items WHERE cart_id = ?", cart.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productWeightID, quantity int
		err := rows.Scan(&productWeightID, &quantity)
		if err != nil {
			return nil, err
		}

		// Retrieve the product weight details from the database
		productWeight, err := cs.getProductWeightByID(productWeightID)
		if err != nil {
			return nil, err
		}

		// Create a new cart item and append it to the cart's items
		cartItem := &models.CartItem{
			Product:   productWeight,
			Quantity:  quantity,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		cart.Items = append(cart.Items, cartItem)
	}

	return cart, nil
}

// GetCartResponse represents a response containing the user's cart and total price
type GetCartResponse struct {
	Cart       *models.Cart `json:"cart"`
	TotalPrice float64      `json:"totalPrice"`
}
