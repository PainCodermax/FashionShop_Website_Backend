package models

import (
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Delivery struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	DeliveryID     string             `json:"deliveryID,omitempty" bson:"delivery_id,omitempty"`
	OrderID        *string            `json:"orderID,omitempty" bson:"order_id,omitempty"`
	DeliveryDate   time.Time          `json:"deliveryDate,omitempty" bson:"delivery_date,omitempty"`
	DeliveryStatus enum.OrderStatus   `json:"deliveryStatus,omitempty" bson:"delivery_status,omitempty"`
	Address        string             `json:"address,omitempty" bson:"address,omitempty"`

	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updtaed_at"`
}
