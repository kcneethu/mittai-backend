package services

import (
	"fmt"

	"github.com/gklps/mittai-backend/db"
	"github.com/gklps/mittai-backend/models"
)

// ListProductWeights retrieves all weight options for a product
func ListProductWeights(productID int) ([]models.ProductWeight, error) {
	productWeights, err := db.ListProductWeights(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product weights: %w", err)
	}

	return productWeights, nil
}

// GetProductWeightByID retrieves a product weight option by its ID
func GetProductWeightByID(productWeightID int) (models.ProductWeight, error) {
	productWeight, err := db.GetProductWeightByID(productWeightID)
	if err != nil {
		return models.ProductWeight{}, fmt.Errorf("failed to get product weight: %w", err)
	}

	return productWeight, nil
}

func ListProducts() ([]models.Product, error) {
	products, err := db.ListProducts()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	return products, nil
}

// GetProductByID retrieves a product by its ID
func GetProductByID(productID int) (models.Product, error) {
	product, err := db.GetProductByID(productID)
	if err != nil {
		return models.Product{}, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}
