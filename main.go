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
	orderStatusService := services.NewOrderStatusService(repo)
	userService := services.NewUserService(repo)
	cartService := services.NewCartService(repo)
	purchaseService := services.NewPurchaseService(repo, productService, cartService, orderStatusService, userService)
	paymentService := services.NewPaymentService(repo)
	addressService := services.NewAddressService(repo)
	wishlistService := services.NewWishlistService(repo)
	// Create more instances of services as needed

	router := mux.NewRouter()
	// corsHandler := handlers.CORS(
	// 	handlers.AllowedOrigins([]string{"http://localhost:3000"}),
	// 	handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	// 	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	// )

	// Register the routes for each service
	productService.RegisterRoutes(router)
	userService.RegisterRoutes(router)
	productWeightService.RegisterRoutes(router)
	cartService.RegisterRoutes(router)
	purchaseService.RegisterRoutes(router)
	paymentService.RegisterRoutes(router)
	addressService.RegisterRoutes(router)
	wishlistService.RegisterRoutes(router)
	orderStatusService.RegisterRoutes(router)
	http.Handle("/", corsMiddleware(router))
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://0.0.0.0:8080/docs/swagger.json"), // Set the URL to your swagger.json
	))
	// http.Handle("/", corsHandler(router))
	// handler := corsMiddleware(router)
	// Register more services' routes as needed

	//router.Use(LoggingMiddleware)

	// Directly serve swagger.json at a specific route
	router.HandleFunc("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	// Continue to serve other static files in /docs
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	// Log URLs
	log.Println("Swagger UI: http://localhost:8080/swagger/")
	log.Println("Swagger JSON Specification: http://localhost:8080/docs/swagger.json")
	log.Println("Database Path:", dbPath)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))

}

// corsMiddleware is a middleware function to set the CORS headers in the response.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Access-Control-Allow-Origin header to allow requests from http://localhost:3000
		//	for key, values := range r.Header {
		//		log.Printf("%s: %v\n", key, values)
		//	}
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Optionally, you can set other CORS headers, such as Access-Control-Allow-Methods, etc.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Allow preflight requests (OPTIONS method) by setting appropriate headers for preflight responses
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
