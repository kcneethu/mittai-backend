package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
	"github.com/gorilla/mux"
)

type WishlistService struct {
	DB *db.Repository
}

func NewWishlistService(db *db.Repository) *WishlistService {
	return &WishlistService{
		DB: db,
	}
}

func (ws *WishlistService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/wishlist/{user_id}", ws.GetWishlistItemsByUserID).Methods(http.MethodGet)
	r.HandleFunc("/wishlist/check/{user_id}/{product_id}", ws.CheckItemInWishlist).Methods(http.MethodGet)
	r.HandleFunc("/wishlist", ws.AddToWishlist).Methods(http.MethodPost)
	r.HandleFunc("/wishlist/{user_id}/{product_id}", ws.RemoveFromWishlist).Methods(http.MethodDelete)
}

// @Summary Add an item to the wishlist
// @Tags Wishlist
// @Accept json
// @Produce json
// @Param item body models.WishlistRequest true "Wishlist item"
// @Success 200 {object} map[string]int "Item added successfully"
// @Failure 400 "Bad request"
// @Failure 500 "Failed to add item to wishlist"
// @Router /wishlist [post]
func (ws *WishlistService) AddToWishlist(w http.ResponseWriter, r *http.Request) {
	var item models.Wishlist
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	insertedID, err := ws.storeItemInWishlist(item)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to add item to wishlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"wishlist_id": insertedID})
}

func (ws *WishlistService) storeItemInWishlist(item models.Wishlist) (int, error) {
	// First, check if the product already exists in the wishlist for the user
	queryCheck := `SELECT COUNT(*) FROM wishlist WHERE user_id = ? AND product_id = ?`
	var count int
	err := ws.DB.QueryRow(queryCheck, item.UserID, item.ProductID).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		log.Println("The product is already in the wishlist for the user")
		return 0, nil
	}

	// Proceed with the insert if it doesn't exist
	query := `INSERT INTO wishlist (user_id, product_id, created_at, updated_at) VALUES (?, ?, ?, ?)`
	result, err := ws.DB.Exec(query, item.UserID, item.ProductID, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return 0, err
	}

	lastInsertedID, err := result.LastInsertId()
	return int(lastInsertedID), err
}

// @Summary Remove an item from the wishlist
// @Tags Wishlist
// @Produce json
// @Param user_id path int true "User ID"
// @Param product_id path int true "Product ID"
// @Success 200 {string} string "Item removed successfully"
// @Failure 400 "Bad request"
// @Failure 500 "Failed to remove item from wishlist"
// @Router /wishlist/{user_id}/{product_id} [delete]
func (ws *WishlistService) RemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	productIDStr := vars["product_id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = ws.removeItemFromWishlist(userID, productID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to remove item from wishlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item removed successfully"))
}

func (ws *WishlistService) removeItemFromWishlist(userID int, productID int) error {
	query := `DELETE FROM wishlist WHERE user_id = ? AND product_id = ?`
	_, err := ws.DB.Exec(query, userID, productID)
	return err
}

// @Summary Get all wishlist items for a user
// @Tags Wishlist
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} models.Wishlist "Wishlist items retrieved successfully"
// @Failure 400 "Bad request"
// @Failure 500 "Failed to fetch wishlist items"
// @Router /wishlist/{user_id} [get]
func (ws *WishlistService) GetWishlistItemsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("Inside GetWishlistItemsByUserID")
	userIDStr := vars["user_id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	wishlistItems, err := ws.fetchWishlistItemsByUserID(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch wishlist items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wishlistItems)
}

func (ws *WishlistService) fetchWishlistItemsByUserID(userID int) ([]*models.Wishlist, error) {
	query := `SELECT id, user_id, product_id, created_at, updated_at FROM wishlist WHERE user_id = ?`
	rows, err := ws.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Wishlist
	for rows.Next() {
		var item models.Wishlist
		err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

// @Summary Check if an item exists in the wishlist
// @Tags Wishlist
// @Produce json
// @Param user_id path int true "User ID"
// @Param product_id path int true "Product ID"
// @Success 200 {object} map[string]bool "Item exists in the wishlist"
// @Failure 400 "Bad request"
// @Failure 500 "Failed to check item in wishlist"
// @Router /wishlist/check/{user_id}/{product_id} [get]
func (ws *WishlistService) CheckItemInWishlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	productIDStr := vars["product_id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	exists, err := ws.itemExistsInWishlist(userID, productID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to check item in wishlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"exists": exists})
}

func (ws *WishlistService) itemExistsInWishlist(userID int, productID int) (bool, error) {
	query := `SELECT COUNT(*) FROM wishlist WHERE user_id = ? AND product_id = ?`
	var count int
	err := ws.DB.QueryRow(query, userID, productID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
