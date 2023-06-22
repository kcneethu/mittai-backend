package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gklps/mittai-backend/models"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	dbPath := filepath.Join(filepath.Dir(exePath), "database.db")

	// Check if the database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Create the database file
		file, err := os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("failed to create database file: %w", err)
		}
		file.Close()
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	db = conn

	// Create tables if they don't exist
	err = createTables()
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	err := db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
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
		return fmt.Errorf("failed to create User table: %w", err)
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
		return fmt.Errorf("failed to create Product table: %w", err)
	}

	// Create Product_Weight table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Product_Weight (
		product_weight_id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER,
		weight TEXT,
		FOREIGN KEY (product_id) REFERENCES Product (product_id)
	)`)
	if err != nil {
		return fmt.Errorf("failed to create Product_Weight table: %w", err)
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
		return fmt.Errorf("failed to create Order table: %w", err)
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
		return fmt.Errorf("failed to create Order_Items table: %w", err)
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
		return fmt.Errorf("failed to create Inventory table: %w", err)
	}

	// Create Supplier table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Supplier (
		supplier_id INTEGER PRIMARY KEY,
		name TEXT,
		contact_number TEXT,
		email TEXT
	)`)
	if err != nil {
		return fmt.Errorf("failed to create Supplier table: %w", err)
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
		return fmt.Errorf("failed to create Cart table: %w", err)
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
		return fmt.Errorf("failed to create Reviews table: %w", err)
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
		return fmt.Errorf("failed to create Promotions table: %w", err)
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
		return fmt.Errorf("failed to create Analytics table: %w", err)
	}

	return nil
}

// FlushAllTables deletes all records from all tables
func FlushAllTables() error {
	_, err := db.Exec("DELETE FROM User")
	if err != nil {
		return fmt.Errorf("failed to flush User table: %w", err)
	}

	_, err = db.Exec("DELETE FROM Product")
	if err != nil {
		return fmt.Errorf("failed to flush Product table: %w", err)
	}

	// Flush records from other tables (Product_Weight, Order, Order_Items, Inventory, Supplier, Cart, Reviews, Promotions, Analytics) similarly

	return nil
}

func CreateUser(user models.User) (int, error) {
	result, err := db.Exec("INSERT INTO User (name, email, contact_number, address) VALUES (?, ?, ?, ?)",
		user.Name, user.Email, user.ContactNumber, user.Address)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve user ID: %w", err)
	}

	return int(userID), nil
}

func CreateProduct(product models.Product) (int, error) {
	result, err := db.Exec("INSERT INTO Product (name, description, category, price, availability, ingredients, nutritional_information, image_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		product.Name, product.Description, product.Category, product.Price, product.Availability, product.Ingredients, product.NutritionalInformation, product.ImageURL)
	if err != nil {
		return 0, fmt.Errorf("failed to create product: %w", err)
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve product ID: %w", err)
	}

	return int(productID), nil
}

// GetUserByID retrieves a user from the User table based on the user ID
func GetUserByID(userID int) (models.User, error) {
	var user models.User

	err := db.QueryRow("SELECT * FROM User WHERE user_id = ?", userID).Scan(&user.UserID, &user.Name, &user.Email, &user.ContactNumber, &user.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// ListProducts retrieves all products from the Product table

func ListProducts() ([]models.Product, error) {
	rows, err := db.Query("SELECT * FROM Product")
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ProductID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Availability, &product.Ingredients, &product.NutritionalInformation, &product.ImageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		productWeights, err := ListProductWeights(product.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product weights: %w", err)
		}
		product.ProductWeights = productWeights

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over products: %w", err)
	}

	return products, nil
}

// GetProductByID retrieves a product from the Product table based on the product ID
func GetProductByID(productID int) (models.Product, error) {
	var product models.Product

	err := db.QueryRow("SELECT * FROM Product WHERE product_id = ?", productID).Scan(&product.ProductID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Availability, &product.Ingredients, &product.NutritionalInformation, &product.ImageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Product{}, fmt.Errorf("product not found")
		}
		return models.Product{}, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// ListProductWeights retrieves all weight options for a product from the Product_Weight table
func ListProductWeights(productID int) ([]models.ProductWeight, error) {
	rows, err := db.Query("SELECT * FROM Product_Weight WHERE product_id = ?", productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product weights: %w", err)
	}
	defer rows.Close()

	var productWeights []models.ProductWeight

	for rows.Next() {
		var productWeight models.ProductWeight
		err := rows.Scan(&productWeight.ProductWeightID, &productWeight.ProductID, &productWeight.Weight)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product weight: %w", err)
		}
		productWeights = append(productWeights, productWeight)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over product weights: %w", err)
	}

	return productWeights, nil
}

// GetProductWeightByID retrieves a product weight option from the Product_Weight table based on the weight ID
func GetProductWeightByID(productWeightID int) (models.ProductWeight, error) {
	var productWeight models.ProductWeight

	err := db.QueryRow("SELECT * FROM Product_Weight WHERE product_weight_id = ?", productWeightID).Scan(&productWeight.ProductWeightID, &productWeight.ProductID, &productWeight.Weight)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ProductWeight{}, fmt.Errorf("product weight not found")
		}
		return models.ProductWeight{}, fmt.Errorf("failed to get product weight: %w", err)
	}

	return productWeight, nil
}

func UpdateUserPhone(userID int, phoneNumber string) error {
	_, err := db.Exec("UPDATE User SET contact_number = ? WHERE user_id = ?", phoneNumber, userID)
	if err != nil {
		return fmt.Errorf("failed to update user phone: %w", err)
	}

	return nil
}

// AddUserAddress adds a new address for a user
func AddUserAddress(userID int, address models.Address) error {
	_, err := db.Exec("INSERT INTO User_Address (user_id, address) VALUES (?, ?)", userID, address)
	if err != nil {
		return fmt.Errorf("failed to add user address: %w", err)
	}

	return nil
}

// UpdateUserAddress updates an existing address for a user
func UpdateUserAddress(userID int, address models.Address) error {
	_, err := db.Exec("UPDATE User_Address SET address = ? WHERE user_id = ? AND address_id = ?", address, userID, address.AddressID)
	if err != nil {
		return fmt.Errorf("failed to update user address: %w", err)
	}

	return nil
}

// GetUserAddresses retrieves all addresses for a user
func GetUserAddresses(userID int) ([]models.Address, error) {
	rows, err := db.Query("SELECT address_id, address FROM User_Address WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user addresses: %w", err)
	}
	defer rows.Close()

	var addresses []models.Address

	for rows.Next() {
		var address models.Address
		err := rows.Scan(&address.AddressID, &address.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user address: %w", err)
		}
		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over user addresses: %w", err)
	}

	return addresses, nil
}

// ListUsers retrieves all users from the database
func ListUsers() ([]models.User, error) {
	rows, err := db.Query("SELECT * FROM User")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.ContactNumber, &user.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over users: %w", err)
	}

	return users, nil
}

// CreateOrder creates a new order in the database and returns the order ID
func CreateOrder(order models.Order) (int, error) {
	result, err := db.Exec("INSERT INTO `Order` (user_id, order_date, total_amount, status, payment_method, delivery_address) VALUES (?, ?, ?, ?, ?, ?)",
		order.UserID, order.OrderDate, order.TotalAmount, order.Status, order.PaymentMethod, order.DeliveryAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve order ID: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start database transaction: %w", err)
	}

	// Insert order items
	for _, item := range order.Items {
		_, err = tx.Exec("INSERT INTO Order_Items (order_id, product_weight_id, quantity, price, customization_options) VALUES (?, ?, ?, ?, ?)",
			orderID, item.ProductWeightID, item.Quantity, item.Price, item.CustomizationOptions)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to commit database transaction: %w", err)
	}

	return int(orderID), nil
}

// ListOrders retrieves all orders from the Order table
func ListOrders() ([]models.Order, error) {
	rows, err := db.Query("SELECT * FROM `Order`")
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.OrderID, &order.UserID, &order.OrderDate, &order.TotalAmount, &order.Status, &order.PaymentMethod, &order.DeliveryAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Retrieve order items for the order
		orderItems, err := GetOrderItemsByOrderID(order.OrderID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items: %w", err)
		}
		order.Items = orderItems

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over orders: %w", err)
	}

	return orders, nil
}

// GetOrderItemsByOrderID retrieves the order items for a given order ID
func GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error) {
	rows, err := db.Query("SELECT * FROM Order_Items WHERE order_id = ?", orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var orderItems []models.OrderItem

	for rows.Next() {
		var orderItem models.OrderItem
		err := rows.Scan(&orderItem.OrderItemID, &orderItem.OrderID, &orderItem.ProductWeightID, &orderItem.Quantity, &orderItem.Price, &orderItem.CustomizationOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		orderItems = append(orderItems, orderItem)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over order items: %w", err)
	}

	return orderItems, nil
}
