package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(dbPath string) error {
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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Order (
		order_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		order_date TEXT,
		total_amount REAL,
		status TEXT,
		payment_method TEXT,
		delivery_address TEXT,
		FOREIGN KEY (user_id) REFERENCES User (user_id)
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
		FOREIGN KEY (order_id) REFERENCES Order (order_id),
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

// Implement other database operations (CRUD) for each table using SQL queries
// Example functions: CreateUser, GetUserByID, ListProducts, GetProductByID, etc.
// Each function will interact with the database and return the desired data
