package models

import (
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	UserID  string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	OrderID string `json:"orderID,omitempty" bson:"order_id,omitempty"`
	CartID  string `json:"cartID,omitempty" bson:"cart_id,omitempty"`

	Quantity int    `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Price    int    `json:"price,omitempty" bson:"price,omitempty"`
	Status   enum.OrderStatus `json:"status,omitempty" bson:"status,omitempty"`

	Items         []CartItem         `json:"Items,omitempty" bson:"items,omitempty"`
	Address       string             `json:"address,omitempty" bson:"address,omitempty"`
	ShipFee       int                `json:"shipFee,omitempty" bson:"ship_fee,omitemty"`
	DileveryDate  time.Time          `json:"deliveryDate,omitempty" bson:"delivery_date,omitempty"`
	PaymentMethod enum.PaymentMethod `json:"paymentMethod,omitempty" bson:"payment_method,omitempty"`
	IsPaid        bool               `json:"isPaid,omitempty" bson:"is_paid,omitempty"`

	Created_At time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Updated_At time.Time `json:"updtaed_at,omitempty" bson:"updtaed_at,omitempty"`
}
