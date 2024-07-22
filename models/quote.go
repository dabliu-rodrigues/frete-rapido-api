package models

type CreateQuoteRequest struct {
	Recipient struct {
		Address struct {
			Zipcode string `json:"zipcode" validate:"required"`
		} `json:"address"`
	} `json:"recipient"`
	Volumes []struct {
		Category   int     `json:"category" validate:"required,gte=0"`
		Amount     int     `json:"amount" validate:"required,gte=0"`
		UnitWeight float64 `json:"unitary_weight" validate:"required,gte=0"`
		Price      float64 `json:"price" validate:"required,gte=0"`
		SKU        string  `json:"sku" validate:"required"`
		Height     float64 `json:"height" validate:"required,gte=0"`
		Width      float64 `json:"width" validate:"required,gte=0"`
		Length     float64 `json:"length" validate:"required,gte=0"`
	} `json:"volumes" validate:"required,dive,required"`
}

type CreateQuoteResponse struct {
	Carrier []Quote
}
type Quote struct {
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline int     `json:"deadline"`
	Price    float64 `json:"price"`
}

type MetricsResponse struct {
	CheapestQuote      float64 `json:"cheapest_quote"`
	MostExpensiveQuote float64 `json:"most_expensive_quote"`
	Services           []struct {
		AveragePrice float64 `json:"average_price"`
		Carrier      string  `json:"carrier"`
		Count        int     `json:"count"`
		TotalPrice   float64 `json:"total_price"`
	} `json:"services"`
}
