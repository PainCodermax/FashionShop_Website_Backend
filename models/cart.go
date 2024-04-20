package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	CartID     string   `json:"cartID,omitempty" bson:"cart_id,omitempty"`
	UserID     string   `json:"userID,omitempty" bson:"user_id,omitempty"`
	Quantity   int      `json:"quantity,omitempty" bson:"quantity,omitempty"`
	TotalPrice int      `json:"totalPrice,omitempty" bson:"-"`
	ItemList   CartItem `json:"itemList,omitempty" bson:"-"`

	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updtaed_at"`
}

type CartItem struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	CartItemID   string `json:"cartItemID,omitempty" bson:"cart_item_id,omitempty"`
	CartID       string `json:"cartID,omitempty" bson:"cart_id,omitempty"`
	ProductID    string `json:"productID,omitempty" bson:"product_id,omitempty"`
	ItemQuantity int    `json:"itemQuantity,omitempty" bson:"item_quantity,omitempty"`

	ProductName *string  `json:"productName,omitempty" bson:"product_name,omitempty"`
	Price       int      `json:"price,omitempty" bson:"price,omitempty"`
	ListImage   []string `json:"listImage,omitempty" bson:"list_image,omitempty"`
	Quantity    int      `json:"quantity,omitempty" bson:"quantity,omitempty"`
	IsRate      bool     `json:"IsRate,omitempty" bson:"_"`
}
