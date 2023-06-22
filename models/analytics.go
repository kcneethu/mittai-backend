package models

// Analytics represents the analytics data for a specific date
type Analytics struct {
	AnalyticsID   int     `json:"analytics_id"`
	Date          string  `json:"date"`
	TotalSales    float64 `json:"total_sales"`
	CustomerCount int     `json:"customer_count"`
	Revenue       float64 `json:"revenue"`
}
