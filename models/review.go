package models

// Review represents a product review in the system
//sort:
//category
//popularity - db manupulation
//recently added

type Review struct {
	ReviewID   int    `json:"review_id"`
	ProductID  int    `json:"product_id"`
	UserID     int    `json:"user_id"`
	Rating     int    `json:"rating"`
	ReviewText string `json:"review_text"`
	ReviewDate string `json:"review_date"`
}
