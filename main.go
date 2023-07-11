package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/handlers"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Mittai Backend API
// @description API documentation for Mittai Backend
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	// Initialize the database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create a new router
	r := chi.NewRouter()

	// Define routes
	r.Get("/users", handlers.UsersHandler)
	r.Post("/users", handlers.UsersHandler)

	r.Get("/products", handlers.ProductsHandler)
	r.Post("/products", handlers.ProductsHandler)

	r.Get("/orders", handlers.OrdersHandler)
	r.Post("/orders", handlers.OrdersHandler)

	// Serve swagger.json file
	r.Handle("/docs/*", http.StripPrefix("/docs", http.FileServer(http.Dir("./docs"))))

	// Swagger API documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"), // The URL pointing to the API definition
	))

	// Start the server
	port := 8080
	fmt.Printf("Server listening on port %d...\n", port)

	swaggerURL := fmt.Sprintf("http://localhost:%d/swagger/index.html", port)
	fmt.Printf("Swagger API documentation available at: %s\n", swaggerURL)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
