package models

type Amount struct {
	Month       string  `bson:"_id"`
	TotalAmount int64 `bson:"totalAmount"`
}

type Report struct {
	TotalUser         int64    `json:"totalUser,omitempty"`
	TotalOrder        int64    `json:"totalOrder,omitempty"`
	TotalOrderSuccess int64    `json:"totalOrderSuccess,omitempty"`
	TotalProduct      int64    `json:"totalProduct,omitempty"`
	Amounts           []Amount `json:"amounts,omitempty"`
}
