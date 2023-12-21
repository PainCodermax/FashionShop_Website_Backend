package models

type ProductResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Product `json:"data"`
	Total   int       `json:"total,omitempty"`
}

type CategoryResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Category `json:"data"`
}

type CartItemResponse struct {
	Status     int        `json:"status"`
	Message    string     `json:"message"`
	Data       []CartItem `json:"data"`
	Total      int        `json:"total,omitempty"`
	TotalPrice int        `json:"total_price,omitempty"`
}

type AddressResponse struct {
	code string
}