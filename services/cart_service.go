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
	_ "github.com/mattn/go-sqlite3"
)

// CartService handles the cart related operations
type CartService struct {
	DB *db.Repository
}

func NewCartService(db *db.Repository) *CartService {
	return &CartService{
		DB: db,
	}
}

// AddToCartRequest represents the request payload for adding a product to the cart
type AddToCartRequest struct {
	UserID          int `json:"user_id"`
	ProductWeightID int `json:"product_weight_id"`
	Quantity        int `json:"quantity"`
}

func (cs *CartService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/cart", cs.AddToCart).Methods(http.MethodPost)
	r.HandleFunc("/cart/{userID}", cs.GetCartByUserID).Methods(http.MethodGet)
	r.HandleFunc("/cart", cs.UpdateCartItem).Methods(http.MethodPut)
	r.HandleFunc("/cart", cs.RemoveCartItem).Methods(http.MethodDelete)
	r.HandleFunc("/cart/clear", cs.ClearCart).Methods(http.MethodDelete)
}

// AddToCart adds a product to the user's cart
// @Summary Add a product to the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body AddToCartRequest true "Add to cart request payload"
// @Success 200 {object} ErrorResponse "Product added to cart successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 500 {object} ErrorResponse "Failed to add product to cart"
// @Router /cart [post]
func (cs *CartService) AddToCart(w http.ResponseWriter, r *http.Request) {
	var request AddToCartRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set the default quantity to 1 if not provided
	if request.Quantity <= 0 {
		request.Quantity = 1
	}

	tx, err := cs.DB.Begin()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to add product to cart", http.StatusInternalServerError)
		return
	}

	// Insert the cart item into the database
	_, err = tx.Exec("INSERT INTO cart (user_id, product_weight_id, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		request.UserID, request.ProductWeightID, request.Quantity, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		tx.Rollback()
		http.Error(w, "Failed to add product to cart", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to add product to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: "Product added to cart successfully",
	})
}

// GetCartResponse represents the response payload for retrieving the user's cart
type GetCartResponse struct {
	CartItems  []*models.CartItem `json:"cart_items"`
	TotalPrice float64            `json:"total_price"`
}

// GetCartItem represents a cart item
type GetCartItem struct {
	Product   *models.ProductWeight `json:"product"`
	Quantity  int                   `json:"quantity"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// GetCartByUserID retrieves the user's cart by user ID
// @Summary Get the user's cart by user ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} GetCartResponse "User's cart retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 500 {object} ErrorResponse "Failed to fetch cart"
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

	rows, err := cs.DB.Query(`SELECT c.id, c.quantity,
			p.id, p.name, p.description, p.category, p.ingredients, p.nutritional_info, p.image_urls,
			w.id, w.weight, w.price, w.stock_availability, w.created_at, w.updated_at
		FROM cart AS c
		INNER JOIN product_weights AS w ON c.product_weight_id = w.id
		INNER JOIN products AS p ON w.product_id = p.id
		WHERE c.user_id = ?`, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch cart", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cartItems []*models.CartItem
	totalPrice := 0.0

	for rows.Next() {
		var (
			cartID           int
			cartItemQuantity int
			productID        int
			productName      string
			productDesc      string
			productCat       string
			productIng       string
			productNutr      string
			productImgURLs   string
			weightID         int
			weightVal        int
			weightPrice      float64
			weightStock      int
			weightCreatedAt  time.Time
			weightUpdatedAt  time.Time
		)

		err := rows.Scan(&cartID, &cartItemQuantity,
			&productID, &productName, &productDesc, &productCat, &productIng, &productNutr, &productImgURLs,
			&weightID, &weightVal, &weightPrice, &weightStock, &weightCreatedAt, &weightUpdatedAt)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to fetch cart", http.StatusInternalServerError)
			return
		}

		// Create a new product weight
		weight := &models.ProductWeight{
			ID:                weightID,
			ProductID:         productID,
			Weight:            weightVal,
			Price:             weightPrice,
			StockAvailability: weightStock,
			CreatedAt:         weightCreatedAt,
			UpdatedAt:         weightUpdatedAt,
		}

		// Create a new cart item
		cartItem := &models.CartItem{
			Product:   weight,
			Quantity:  cartItemQuantity,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Add the cart item to the list
		cartItems = append(cartItems, cartItem)

		// Calculate the total price
		totalPrice += float64(cartItem.Quantity) * weight.Price
	}

	response := GetCartResponse{
		CartItems:  cartItems,
		TotalPrice: totalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateCartItemRequest represents the request payload for updating a cart item
type UpdateCartItemRequest struct {
	UserID          int `json:"user_id"`
	ProductWeightID int `json:"product_weight_id"`
	Quantity        int `json:"quantity"`
}

// UpdateCartItem updates the quantity of a cart item
// @Summary Update the quantity of a cart item
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body UpdateCartItemRequest true "Update cart item request payload"
// @Success 200 {object} ErrorResponse "Cart item updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 500 {object} ErrorResponse "Failed to update cart item"
// @Router /cart [put]
func (cs *CartService) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	var request UpdateCartItemRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := cs.DB.Begin()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update cart item", http.StatusInternalServerError)
		return
	}

	// Update the cart item quantity in the database
	_, err = tx.Exec("UPDATE cart SET quantity = ? WHERE user_id = ? AND product_weight_id = ?",
		request.Quantity, request.UserID, request.ProductWeightID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		http.Error(w, "Failed to update cart item", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: "Cart item updated successfully",
	})
}

// RemoveCartItemRequest represents the request payload for removing a cart item
type RemoveCartItemRequest struct {
	UserID          int `json:"user_id"`
	ProductWeightID int `json:"product_weight_id"`
}

// RemoveCartItem removes a cart item
// @Summary Remove a cart item
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body RemoveCartItemRequest true "Remove cart item request payload"
// @Success 200 {object} ErrorResponse "Cart item removed successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 500 {object} ErrorResponse "Failed to remove cart item"
// @Router /cart [delete]
func (cs *CartService) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	var request RemoveCartItemRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := cs.DB.Begin()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to remove cart item", http.StatusInternalServerError)
		return
	}

	// Delete the cart item from the database
	_, err = tx.Exec("DELETE FROM cart WHERE user_id = ? AND product_weight_id = ?",
		request.UserID, request.ProductWeightID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		http.Error(w, "Failed to remove cart item", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to remove cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: "Cart item removed successfully",
	})
}

// ClearCartRequest represents the request payload for clearing the cart
type ClearCartRequest struct {
	UserID int `json:"user_id"`
}

// ClearCart removes all items from the cart
// @Summary Clear the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body ClearCartRequest true "Clear cart request payload"
// @Success 200 {object} ErrorResponse "Cart cleared successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 500 {object} ErrorResponse "Failed to clear cart"
// @Router /cart/clear [delete]
func (cs *CartService) ClearCart(w http.ResponseWriter, r *http.Request) {
	var request ClearCartRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := cs.DB.Begin()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	// Delete all cart items for the user from the database
	_, err = tx.Exec("DELETE FROM cart WHERE user_id = ?", request.UserID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: "Cart cleared successfully",
	})
}
