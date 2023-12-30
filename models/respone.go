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
	Status   int        `json:"status"`
	Message  string     `json:"message"`
	Data     []CartItem `json:"data"`
	Total    int        `json:"total,omitempty"`
	Province string     `json:"province,omitempty" bson:"-"`
	District string     `json:"district,omitempty" bson:"-"`
	Ward     string     `json:"ward,omitempty" bson:"-"`
	ShipFee  int        `json:"shipFee,omitempty" bson:"-"`

	TotalPrice int `json:"total_price,omitempty"`
}

type ProvinceResponse struct {
	Code    int        `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Data    []Province `json:"data,omitempty"`
}

type DistrictResponse struct {
	Code    int        `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Data    []District `json:"data,omitempty"`
}

type WardResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    []Ward `json:"data,omitempty"`
}

type OrderResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    []Order `json:"data"`
	Total   int     `json:"total,omitempty"`
}

type ShipmentResponse struct {
	Code    int        `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Data    []ShipMent `json:"data,omitempty"`
}
