package db

import (
	"database/sql"
	"log"
)

type Repository struct {
	*sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// createProductTable creates the product table in the database if it doesn't exist or modifies the table structure
func (r *Repository) createProductTable() error {
	query := `CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		category TEXT NOT NULL,
		ingredients TEXT NOT NULL,
		nutritional_info TEXT NOT NULL,
		image_urls TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	_, err := r.Exec(query)
	if err != nil {
		return err
	}

	// Create the product_weights table if it doesn't exist
	query = `CREATE TABLE IF NOT EXISTS product_weights (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL,
		weight FLOAT NOT NULL,
		price FLOAT NOT NULL,
		stock_availability INTEGER NOT NULL, -- Modified field name
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (product_id) REFERENCES products (id)
	);`

	_, err = r.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// createUserTable creates the user table in the database if it doesn't exist or modifies the table structure
func (r *Repository) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL,
		contact_number TEXT NOT NULL UNIQUE,
		verified_account BOOLEAN NOT NULL DEFAULT 0
	);`

	_, err := r.Exec(query)
	return err
}

// createAddressTable creates the address table in the database if it doesn't exist or modifies the table structure
func (r *Repository) createAddressTable() error {
	query := `CREATE TABLE IF NOT EXISTS addresses (
		address_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		address_line1 TEXT NOT NULL,
		address_line2 TEXT,
		city TEXT NOT NULL,
		state TEXT NOT NULL,
		zip_code TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (user_id)
	);`

	_, err := r.Exec(query)
	return err
}

// createCartTable creates the cart table in the database if it doesn't exist or modifies the table structure
func (r *Repository) createCartTable() error {
	query := `CREATE TABLE IF NOT EXISTS carts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		items TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (user_id)
	);`

	_, err := r.Exec(query)
	return err
}

// createCartItemTable creates the cart_item table in the database if it doesn't exist or modifies the table structure
func (r *Repository) createCartItemTable() error {
	query := `CREATE TABLE IF NOT EXISTS cart_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cart_id INTEGER NOT NULL,
		product_weight_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		user_id INTEGER NOT NULL, -- Add the user_id column
		FOREIGN KEY (cart_id) REFERENCES carts (id),
		FOREIGN KEY (product_weight_id) REFERENCES product_weights (id),
		FOREIGN KEY (user_id) REFERENCES users (user_id) -- Add the foreign key constraint
	);`

	_, err := r.Exec(query)
	return err
}

// CreateTables creates or updates all necessary tables in the database
func (r *Repository) CreateTables() {
	if err := r.createProductTable(); err != nil {
		log.Fatal(err)
	}
	if err := r.createUserTable(); err != nil {
		log.Fatal(err)
	}
	if err := r.createAddressTable(); err != nil {
		log.Fatal(err)
	}
	if err := r.createCartTable(); err != nil {
		log.Fatal(err)
	}
	if err := r.createCartItemTable(); err != nil {
		log.Fatal(err)
	}
}
