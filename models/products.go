package models

import (
	"time"
)

// Product represents a product in the system
type Product struct {
	ID              int              `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Category        string           `json:"category"`
	Ingredients     string           `json:"ingredients"`
	NutritionalInfo string           `json:"nutritional_info"`
	ImageURLs       []string         `json:"image_urls"`
	Weights         []*ProductWeight `json:"weights"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// ProductWeight represents a specific weight variant of a product
type ProductWeight struct {
	ID                   int
	ProductID            int
	Weight               int
	Price                float64
	StockAvailability    int
	Measurement          string
	CreatedAt, UpdatedAt time.Time
}
