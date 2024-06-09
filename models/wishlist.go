package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WishItem struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	WishItemId string             `json:"wishItemId,omitempty" bson:"wish_list_id,omitempty"`
	UserId     string             `json:"userId,omitempty" bson:"user_id,omitempty"`
	ProductId  string             `json:"productId,omitempty" bson:"product_id,omitempty"`
	Product    Product            `json:"product,omitempty" bson:"product,omitempty"`

	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updtaed_at"`
}
