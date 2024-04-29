package models

type RequestOrder struct {
	CartID     string     `json:"cartID,omitempty"`
	CartItems  []CartItem `json:"items,omitempty"`
	Address    string     `json:"address,omitempty"`
	TotalPrice int        `json:"totalPice,omitempty"`
	Quantity   int        `json:"quantity,omitempty"`
	ShipFee    int        `json:"shipFee,omitempty" bson:"ship_fee,omitemty"`
}

type RequestPayment struct {
	OrderId string `json:"orderId,omitempty"`
	Amount  string `json:"amount,omitempty"`
}
