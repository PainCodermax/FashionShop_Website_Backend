package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rating struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Content    string             `json:"content,omitempty" bson:"content,omitempty"`
	Score      int                `json:"score,omitempty" bson:"score,omitempty"`
	OrderID    *string            `json:"orderID,omitempty" bson:"order_id,omitempty"`
	Product_ID string             `json:"productId,omitempty" bson:"product_id,omitempty"`
	User_ID    string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Order      *Order             `json:"order,omitempty" bson:"order,omitempty"`
	User       *User              `json:"user,omitempty" bson:"user,omitempty"`
	Product    *Product           `json:"product,omitempty" bson:"product,omitempty"`

	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updtaed_at"`
}
