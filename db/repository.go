package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/gklps/mittai-backend/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	exePath, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "failed to get executable path")
	}

	dbPath := filepath.Join(filepath.Dir(exePath), "database.db")

	// Check if the database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Create the database file
		file, err := os.Create(dbPath)
		if err != nil {
			return errors.Wrap(err, "failed to create database file")
		}
		file.Close()
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Wrap(err, "failed to open database connection")
	}
	db = conn

	// Create tables if they don't exist
	err = createTables()
	if err != nil {
		return errors.Wrap(err, "failed to create tables")
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	err := db.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close database connection")
	}

	return nil
}

// createTables creates database tables if they don't exist
func createTables() error {
	// Create User table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS User (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		contact_number TEXT,
		address TEXT
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create User table")
	}

	// Create Product table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Product (
		product_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		description TEXT,
		category TEXT,
		price REAL,
		availability INTEGER,
		ingredients TEXT,
		nutritional_information TEXT,
		image_url TEXT
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create Product table")
	}

	// Create Product_Weight table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Product_Weight (
		product_weight_id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER,
		weight TEXT,
		FOREIGN KEY (product_id) REFERENCES Product (product_id)
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create Product_Weight table")
	}

	// Create Order table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "Order" (
		order_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		order_date TEXT,
		total_amount REAL,
		status TEXT,
		payment_method TEXT,
		delivery_address TEXT
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create Order table")
	}

	// Create Order_Items table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Order_Items (
		order_item_id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id INTEGER,
		product_weight_id INTEGER,
		quantity INTEGER,
		price REAL,
		customization_options TEXT,
		FOREIGN KEY (order_id) REFERENCES "Order" (order_id),
		FOREIGN KEY (product_weight_id) REFERENCES Product_Weight (product_weight_id)
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create Order_Items table")
	}

	// Create Inventory table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Inventory (
		product_weight_id INTEGER PRIMARY KEY,
		available_quantity INTEGER,
		reorder_threshold INTEGER,
		supplier_id INTEGER,
		FOREIGN KEY (product_weight_id) REFERENCES Product_Weight (product_weight_id),
		FOREIGN KEY (supplier_id) REFERENCES Supplier (supplier_id)
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create Inventory table")
	}

	// Create Supplier table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Supplier (
		supplier_id INTEGER PRIMARY KEY,
		name TEXT,
		contact_number TEXT,
		email TEXT
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create Supplier table")
	}

	// Create Cart table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Cart (
		cart_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		product_weight_id INTEGER,
		quantity INTEGER,
		applied_discounts TEXT,
		FOREIGN KEY (user_id) REFERENCES User (user_id),
		FOREIGN KEY (product_weight_id) REFERENCES Product_Weight (product_weight_id)
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create Cart table")
	}

	// Create Reviews table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Reviews (
		review_id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER,
		user_id INTEGER,
		rating INTEGER,
		review_text TEXT,
		review_date TEXT,
		FOREIGN KEY (product_id) REFERENCES Product (product_id),
		FOREIGN KEY (user_id) REFERENCES User (user_id)
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create Reviews table")
	}

	// Create Promotions table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Promotions (
		promotion_id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER,
		discount_code TEXT,
		start_date TEXT,
		end_date TEXT,
		discount_percentage REAL,
		FOREIGN KEY (product_id) REFERENCES Product (product_id)
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create Promotions table")
	}

	// Create Analytics table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Analytics (
		analytics_id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		total_sales REAL,
		customer_count INTEGER,
		revenue REAL
	)`)
	if err != nil {
		return errors.Wrap(err, "failed to create Analytics table")
	}

	return nil
}

// FlushAllTables deletes all records from all tables
func FlushAllTables() error {
	_, err := db.Exec("DELETE FROM User")
	if err != nil {
		return errors.Wrap(err, "failed to flush User table")
	}

	_, err = db.Exec("DELETE FROM Product")
	if err != nil {
		return errors.Wrap(err, "failed to flush Product table")
	}

	// Flush records from other tables (Product_Weight, Order, Order_Items, Inventory, Supplier, Cart, Reviews, Promotions, Analytics) similarly

	return nil
}

func CreateUser(user models.User) (int, error) {
	result, err := db.Exec("INSERT INTO User (name, email, contact_number, address) VALUES (?, ?, ?, ?)",
		user.Name, user.Email, user.ContactNumber, user.Address)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create user")
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve user ID")
	}

	return int(userID), nil
}

func CreateProduct(product models.Product) (int, error) {
	result, err := db.Exec("INSERT INTO Product (name, description, category, price, availability, ingredients, nutritional_information, image_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		product.Name, product.Description, product.Category, product.Price, product.Availability, product.Ingredients, product.NutritionalInformation, product.ImageURL)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create product")
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve product ID")
	}

	return int(productID), nil
}

// GetUserByID retrieves a user from the User table based on the user ID
func GetUserByID(userID int) (models.User, error) {
	var user models.User

	err := db.QueryRow("SELECT * FROM User WHERE user_id = ?", userID).Scan(&user.UserID, &user.Name, &user.Email, &user.ContactNumber, &user.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, errors.Wrap(err, "failed to get user")
	}

	return user, nil
}

// ListProducts retrieves all products from the Product table
func ListProducts() ([]models.Product, error) {
	rows, err := db.Query("SELECT * FROM Product")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get products")
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ProductID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Availability, &product.Ingredients, &product.NutritionalInformation, &product.ImageURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan product")
		}

		productWeights, err := ListProductWeights(product.ProductID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get product weights")
		}
		product.ProductWeights = productWeights

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error while iterating over products")
	}

	return products, nil
}

// GetProductByID retrieves a product from the Product table based on the product ID
func GetProductByID(productID int) (models.Product, error) {
	var product models.Product

	err := db.QueryRow("SELECT * FROM Product WHERE product_id = ?", productID).Scan(&product.ProductID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Availability, &product.Ingredients, &product.NutritionalInformation, &product.ImageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Product{}, errors.New("product not found")
		}
		return models.Product{}, errors.Wrap(err, "failed to get product")
	}

	return product, nil
}

// ListProductWeights retrieves all weight options for a product from the Product_Weight table
func ListProductWeights(productID int) ([]models.ProductWeight, error) {
	rows, err := db.Query("SELECT * FROM Product_Weight WHERE product_id = ?", productID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get product weights")
	}
	defer rows.Close()

	var productWeights []models.ProductWeight

	for rows.Next() {
		var productWeight models.ProductWeight
		err := rows.Scan(&productWeight.ProductWeightID, &productWeight.ProductID, &productWeight.Weight)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan product weight")
		}

		productWeights = append(productWeights, productWeight)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error while iterating over product weights")
	}

	return productWeights, nil
}

// UpdateProductAvailability updates the availability of a product in the Product table
func UpdateProductAvailability(productID int, availability int) error {
	_, err := db.Exec("UPDATE Product SET availability = ? WHERE product_id = ?", availability, productID)
	if err != nil {
		return errors.Wrap(err, "failed to update product availability")
	}

	return nil
}

// CreateOrder creates a new order in the Order table
func CreateOrder(order models.Order) (int, error) {
	result, err := db.Exec("INSERT INTO \"Order\" (user_id, order_date, total_amount, status, payment_method, delivery_address) VALUES (?, ?, ?, ?, ?, ?)",
		order.UserID, order.OrderDate, order.TotalAmount, order.Status, order.PaymentMethod, order.DeliveryAddress)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create order")
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve order ID")
	}

	return int(orderID), nil
}

// AddOrderItem adds a new item to an order in the Order_Items table
func AddOrderItem(orderItem models.OrderItem) (int, error) {
	result, err := db.Exec("INSERT INTO Order_Items (order_id, product_weight_id, quantity, price, customization_options) VALUES (?, ?, ?, ?, ?)",
		orderItem.OrderID, orderItem.ProductWeightID, orderItem.Quantity, orderItem.Price, orderItem.CustomizationOptions)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add order item")
	}

	orderItemID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve order item ID")
	}

	return int(orderItemID), nil
}

// ListOrderItems retrieves all items for an order from the Order_Items table
func ListOrderItems(orderID int) ([]models.OrderItem, error) {
	rows, err := db.Query("SELECT * FROM Order_Items WHERE order_id = ?", orderID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get order items")
	}
	defer rows.Close()

	var orderItems []models.OrderItem

	for rows.Next() {
		var orderItem models.OrderItem
		err := rows.Scan(&orderItem.OrderItemID, &orderItem.OrderID, &orderItem.ProductWeightID, &orderItem.Quantity, &orderItem.Price, &orderItem.CustomizationOptions)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan order item")
		}

		orderItems = append(orderItems, orderItem)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error while iterating over order items")
	}

	return orderItems, nil
}

// CreateInventory creates a new inventory record in the Inventory table
func CreateInventory(inventory models.Inventory) error {
	_, err := db.Exec("INSERT INTO Inventory (product_weight_id, available_quantity, reorder_threshold, supplier_id) VALUES (?, ?, ?, ?)",
		inventory.ProductWeightID, inventory.AvailableQuantity, inventory.ReorderThreshold, inventory.SupplierID)
	if err != nil {
		return errors.Wrap(err, "failed to create inventory")
	}

	return nil
}

// UpdateInventory updates the inventory record in the Inventory table
func UpdateInventory(inventory models.Inventory) error {
	_, err := db.Exec("UPDATE Inventory SET available_quantity = ?, reorder_threshold = ?, supplier_id = ? WHERE product_weight_id = ?",
		inventory.AvailableQuantity, inventory.ReorderThreshold, inventory.SupplierID, inventory.ProductWeightID)
	if err != nil {
		return errors.Wrap(err, "failed to update inventory")
	}

	return nil
}

// GetInventoryByProductWeightID retrieves the inventory record from the Inventory table based on the product weight ID
func GetInventoryByProductWeightID(productWeightID int) (models.Inventory, error) {
	var inventory models.Inventory

	err := db.QueryRow("SELECT * FROM Inventory WHERE product_weight_id = ?", productWeightID).Scan(&inventory.ProductWeightID, &inventory.AvailableQuantity, &inventory.ReorderThreshold, &inventory.SupplierID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Inventory{}, errors.New("inventory not found")
		}
		return models.Inventory{}, errors.Wrap(err, "failed to get inventory")
	}

	return inventory, nil
}

// CreateSupplier creates a new supplier in the Supplier table
func CreateSupplier(supplier models.Supplier) (int, error) {
	result, err := db.Exec("INSERT INTO Supplier (name, contact_number, email) VALUES (?, ?, ?)",
		supplier.Name, supplier.ContactNumber, supplier.Email)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create supplier")
	}

	supplierID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve supplier ID")
	}

	return int(supplierID), nil
}

// GetSupplierByID retrieves a supplier from the Supplier table based on the supplier ID
func GetSupplierByID(supplierID int) (models.Supplier, error) {
	var supplier models.Supplier

	err := db.QueryRow("SELECT * FROM Supplier WHERE supplier_id = ?", supplierID).Scan(&supplier.SupplierID, &supplier.Name, &supplier.ContactNumber, &supplier.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Supplier{}, errors.New("supplier not found")
		}
		return models.Supplier{}, errors.Wrap(err, "failed to get supplier")
	}

	return supplier, nil
}

// CreateCart creates a new cart in the Cart table
func CreateCart(cart models.Cart) (int, error) {
	result, err := db.Exec("INSERT INTO Cart (user_id, product_weight_id, quantity, applied_discounts) VALUES (?, ?, ?, ?)",
		cart.UserID, cart.ProductWeightID, cart.Quantity, cart.AppliedDiscounts)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create cart")
	}

	cartID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve cart ID")
	}

	return int(cartID), nil
}

// GetCartByID retrieves a cart from the Cart table based on the cart ID
func GetCartByID(cartID int) (models.Cart, error) {
	var cart models.Cart

	err := db.QueryRow("SELECT * FROM Cart WHERE cart_id = ?", cartID).Scan(&cart.CartID, &cart.UserID, &cart.ProductWeightID, &cart.Quantity, &cart.AppliedDiscounts)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Cart{}, errors.New("cart not found")
		}
		return models.Cart{}, errors.Wrap(err, "failed to get cart")
	}

	return cart, nil
}

// CreateReview creates a new review in the Reviews table
func CreateReview(review models.Review) (int, error) {
	result, err := db.Exec("INSERT INTO Reviews (product_id, user_id, rating, review_text, review_date) VALUES (?, ?, ?, ?, ?)",
		review.ProductID, review.UserID, review.Rating, review.ReviewText, review.ReviewDate)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create review")
	}

	reviewID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve review ID")
	}

	return int(reviewID), nil
}

// GetReviewByID retrieves a review from the Reviews table based on the review ID
func GetReviewByID(reviewID int) (models.Review, error) {
	var review models.Review

	err := db.QueryRow("SELECT * FROM Reviews WHERE review_id = ?", reviewID).Scan(&review.ReviewID, &review.ProductID, &review.UserID, &review.Rating, &review.ReviewText, &review.ReviewDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Review{}, errors.New("review not found")
		}
		return models.Review{}, errors.Wrap(err, "failed to get review")
	}

	return review, nil
}

// CreatePromotion creates a new promotion in the Promotions table
func CreatePromotion(promotion models.Promotion) (int, error) {
	result, err := db.Exec("INSERT INTO Promotions (product_id, discount_code, start_date, end_date, discount_percentage) VALUES (?, ?, ?, ?, ?)",
		promotion.ProductID, promotion.DiscountCode, promotion.StartDate, promotion.EndDate, promotion.DiscountPercentage)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create promotion")
	}

	promotionID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve promotion ID")
	}

	return int(promotionID), nil
}

// GetPromotionByID retrieves a promotion from the Promotions table based on the promotion ID
func GetPromotionByID(promotionID int) (models.Promotion, error) {
	var promotion models.Promotion

	err := db.QueryRow("SELECT * FROM Promotions WHERE promotion_id = ?", promotionID).Scan(&promotion.PromotionID, &promotion.ProductID, &promotion.DiscountCode, &promotion.StartDate, &promotion.EndDate, &promotion.DiscountPercentage)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Promotion{}, errors.New("promotion not found")
		}
		return models.Promotion{}, errors.Wrap(err, "failed to get promotion")
	}

	return promotion, nil
}

// CreateAnalytics creates a new analytics record in the Analytics table
func CreateAnalytics(analytics models.Analytics) (int, error) {
	result, err := db.Exec("INSERT INTO Analytics (date, total_sales, customer_count, revenue) VALUES (?, ?, ?, ?)",
		analytics.Date, analytics.TotalSales, analytics.CustomerCount, analytics.Revenue)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create analytics")
	}

	analyticsID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve analytics ID")
	}

	return int(analyticsID), nil
}

// GetAnalyticsByID retrieves an analytics record from the Analytics table based on the analytics ID
func GetAnalyticsByID(analyticsID int) (models.Analytics, error) {
	var analytics models.Analytics

	err := db.QueryRow("SELECT * FROM Analytics WHERE analytics_id = ?", analyticsID).Scan(&analytics.AnalyticsID, &analytics.Date, &analytics.TotalSales, &analytics.CustomerCount, &analytics.Revenue)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Analytics{}, errors.New("analytics record not found")
		}
		return models.Analytics{}, errors.Wrap(err, "failed to get analytics record")
	}

	return analytics, nil
}

func ListUsers() ([]models.User, error) {
	rows, err := db.Query("SELECT * FROM User")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.ContactNumber, &user.Address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan user")
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error while iterating over users")
	}

	return users, nil
}

// UpdateUserPhone updates the phone number for a user
func UpdateUserPhone(userID int, phoneNumber string) error {
	_, err := db.Exec("UPDATE User SET contact_number = ? WHERE user_id = ?", phoneNumber, userID)
	if err != nil {
		return errors.Wrap(err, "failed to update user phone")
	}

	return nil
}

// UpdateUserAddress updates an existing address for a user
func UpdateUserAddress(userID int, address string) error {
	_, err := db.Exec("UPDATE User SET address = ? WHERE user_id = ?", address, userID)
	if err != nil {
		return errors.Wrap(err, "failed to update user address")
	}

	return nil
}

// AddUserAddress adds a new address for a user
func AddUserAddress(userID int, address string) error {
	_, err := db.Exec("UPDATE User SET address = ? WHERE user_id = ?", address, userID)
	if err != nil {
		return errors.Wrap(err, "failed to add user address")
	}

	return nil
}

// GetUserAddresses retrieves all addresses for a user
func GetUserAddresses(userID int) ([]models.Address, error) {
	rows, err := db.Query("SELECT address FROM User WHERE user_id = ?", userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user addresses")
	}
	defer rows.Close()

	var addresses []models.Address

	for rows.Next() {
		var address models.Address
		err := rows.Scan(&address.Address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan user address")
		}

		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error while iterating over user addresses")
	}

	return addresses, nil
}

func ListOrders() ([]models.Order, error) {
	rows, err := db.Query("SELECT * FROM \"Order\"")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get orders")
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.OrderID, &order.UserID, &order.OrderDate, &order.TotalAmount, &order.Status, &order.PaymentMethod, &order.DeliveryAddress)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan order")
		}

		order.Items, err = ListOrderItems(order.OrderID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get order items")
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error while iterating over orders")
	}

	return orders, nil
}