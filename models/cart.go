package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`

	CartItems []CartItem `json:"cartItems,omitempty" bson:"cart_items,omitempty"`
	CartID    string     `json:"cartID,omitempty" bson:"cart_id,omitempty"`
	UserID    string     `json:"userID,omitempty" bson:"user_id,omitempty"`

	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updtaed_at"`
}

type CartItem struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`

	CartItemID string `json:"cartItemID,omitempty" bson:"cart_item_id,omitempty"`
	ProductID  string `json:"productID,omitempty" bson:"product_id,omitempty"`
	Quantity   int    `json:"quantity,omitempty" bson:"quantity,omitempty"`
}
