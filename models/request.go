package models

import (
	"net/http"

	"github.com/PainCodermax/FashionShop_Website_Backend/enum"
)

type RequestOrder struct {
	CartID        string             `json:"cartID,omitempty"`
	CartItems     []CartItem         `json:"items,omitempty"`
	Address       string             `json:"address,omitempty"`
	TotalPrice    int                `json:"totalPice,omitempty"`
	Quantity      int                `json:"quantity,omitempty"`
	ShipFee       int                `json:"shipFee,omitempty" bson:"ship_fee,omitemty"`
	PaymentMethod enum.PaymentMethod `json:"paymentMethod,omitempty" bson:"payment_method,omitempty"`
}

type RequestPayment struct {
	OrderId string `json:"orderId,omitempty"`
	Amount  string `json:"amount,omitempty"`
}

type RequestResult struct {
	Response *http.Response
	Error    error
}
