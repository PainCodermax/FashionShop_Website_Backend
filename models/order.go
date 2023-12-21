package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Product_ID string `json:"orderID,omitempty" bson:"order_id,omitempty"`
	CartID     string `json:"cartID,omitempty" bson:"cart_id,omitempty"`

	Quantity int    `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Price    int    `json:"price,omitempty" bson:"price,omitempty"`
	Status   string `json:"status,omitempty" bson:"status,omitempty"`
	PaymentMethod string `json:"payment"`

	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updtaed_at"`
}
