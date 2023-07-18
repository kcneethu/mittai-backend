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
		stmt, err := cs.Repository.DB.Prepare("INSERT INTO carts (user_id, items, created_at, updated_at) VALUES (?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		// Create an array of cart item IDs
		cartItemIDs := make([]int, len(cart.Items))
		for i, item := range cart.Items {
			cartItemIDs[i] = item.Product.ID
		}

		// Convert the array of cart item IDs to JSON
		itemsJSON, err := json.Marshal(cartItemIDs)
		if err != nil {
			return err
		}

		result, err := stmt.Exec(cart.UserID, itemsJSON, cart.CreatedAt, cart.UpdatedAt)
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
		// Fetch the current items from the database
		currentItems, err := cs.getCartItemsByCartID(cart.ID)
		if err != nil {
			return err
		}

		// Append the new items to the current items
		currentItems = append(currentItems, cart.Items...)

		// Create an array of cart item IDs
		cartItemIDs := make([]int, len(currentItems))
		for i, item := range currentItems {
			cartItemIDs[i] = item.Product.ID
		}

		// Convert the array of cart item IDs to JSON
		itemsJSON, err := json.Marshal(cartItemIDs)
		if err != nil {
			return err
		}

		stmt, err := cs.Repository.DB.Prepare("UPDATE carts SET items = ?, created_at = ?, updated_at = ? WHERE id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(itemsJSON, cart.CreatedAt, cart.UpdatedAt, cart.ID)
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
	stmt, err := cs.Repository.DB.Prepare("UPDATE cart_items SET quantity = ?, updated_at = ? WHERE cart_id = ? AND product_weight_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(quantity, time.Now(), userID, productWeightID)
	if err != nil {
		return err
	}

	return nil
}

// getCartByUserID retrieves the user's cart from the database
func (cs *CartService) getCartByUserID(userID int) (*models.Cart, error) {
	query := `
		SELECT c.id, c.user_id, c.created_at, c.updated_at, ci.id, ci.product_weight_id, ci.quantity, ci.created_at, ci.updated_at
		FROM carts AS c
		LEFT JOIN cart_items AS ci ON c.id = ci.cart_id
		WHERE c.user_id = ?
	`

	// Prepare the query
	stmt, err := cs.Repository.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart := &models.Cart{
		UserID: userID,
		Items:  make([]*models.CartItem, 0),
	}

	for rows.Next() {
		var (
			cartID            int
			userID            int
			cartCreatedAt     time.Time
			cartUpdatedAt     time.Time
			cartItemID        sql.NullInt64
			productWeightID   int
			quantity          sql.NullInt64
			cartItemCreatedAt sql.NullTime
			cartItemUpdatedAt sql.NullTime
		)

		// Scan the row into the variables
		err := rows.Scan(
			&cartID, &userID, &cartCreatedAt, &cartUpdatedAt,
			&cartItemID, &productWeightID, &quantity, &cartItemCreatedAt, &cartItemUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Create a new cart if it doesn't exist
		if cart.ID == 0 {
			cart.ID = cartID
			cart.CreatedAt = cartCreatedAt
			cart.UpdatedAt = cartUpdatedAt
		}

		// Add the cart item to the cart
		if cartItemID.Valid {
			cartItem := &models.CartItem{
				Product: &models.ProductWeight{
					ID: productWeightID,
				},
				Quantity:  int(quantity.Int64),
				CreatedAt: cartItemCreatedAt.Time,
				UpdatedAt: cartItemUpdatedAt.Time,
			}
			cart.Items = append(cart.Items, cartItem)
		}
	}

	return cart, nil
}

// getCartItemsByCartID retrieves the cart items for a specific cart ID from the database
func (cs *CartService) getCartItemsByCartID(cartID int) ([]*models.CartItem, error) {
	query := `
		SELECT ci.id, ci.product_weight_id, ci.quantity, ci.created_at, ci.updated_at
		FROM cart_items AS ci
		WHERE ci.cart_id = ?
	`

	// Prepare the query
	stmt, err := cs.Repository.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*models.CartItem, 0)

	for rows.Next() {
		var (
			itemID          int
			productWeightID int
			quantity        int
			itemCreatedAt   time.Time
			itemUpdatedAt   time.Time
		)

		// Scan the row into the variables
		err := rows.Scan(
			&itemID, &productWeightID, &quantity, &itemCreatedAt, &itemUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Create a new cart item and add it to the list
		item := &models.CartItem{
			Product:   &models.ProductWeight{ID: productWeightID},
			Quantity:  quantity,
			CreatedAt: itemCreatedAt,
			UpdatedAt: itemUpdatedAt,
		}

		items = append(items, item)
	}

	return items, nil
}

// GetCartResponse represents a response containing the user's cart and total price
type GetCartResponse struct {
	Cart       *models.Cart `json:"cart"`
	TotalPrice float64      `json:"totalPrice"`
}
