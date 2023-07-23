package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/services"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Get the current directory path
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get the current directory:", err)
	}

	// Set the path of the database file
	dbPath := filepath.Join(dir, "database.db")

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Create a new database file
		log.Println("Database file does not exist. Creating a new database.")

		// Create the database file
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatal("Failed to create the database file:", err)
		}
		file.Close()

		// Run the database initialization code here
		// ...

	} else {
		log.Println("Database file found. Connecting to the existing database.")
	}

	// Open a connection to the SQLite database
	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer dbConn.Close()

	// Create a new instance of the repository
	repo := db.NewRepository(dbConn)

	// Create necessary tables in the database
	repo.CreateTables()

	// Create instances of the services
	productService := services.NewProductService(repo)
	productWeightService := services.NewProductWeightService(repo)
	userService := services.NewUserService(repo)
	cartService := services.NewCartService(repo)
	purchaseService := services.NewPurchaseService(repo, productService)
	paymentService := services.NewPaymentService(repo)
	// Create more instances of services as needed

	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Register the routes for each service
	productService.RegisterRoutes(router)
	userService.RegisterRoutes(router)
	productWeightService.RegisterRoutes(router)
	cartService.RegisterRoutes(router)
	purchaseService.RegisterRoutes(router)
	paymentService.RegisterRoutes(router)
	// Register more services' routes as needed

	// Add CORS support using the cors package
	corsHandler := cors.Default()
	router.Use(corsHandler.Handler) // Add the corsHandler to the router's middleware

	// Set up Swagger
	swaggerURL := "/docs/swagger.json"
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerURL), // The url pointing to API definition
	))

	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	// Log URLs
	log.Println("Swagger UI (API Documentation): http://localhost:8080/swagger/")
	log.Println("Swagger JSON Specification: http://localhost:8080/docs" + swaggerURL)
	log.Println("Database Path:", dbPath)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", router))

}
