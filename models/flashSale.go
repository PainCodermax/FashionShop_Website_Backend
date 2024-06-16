package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashSale struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`

	FlashSaleId   string    `json:"flashSaleId,omitempty" bson:"flash_sale_id,omitempty"`
	ProductIdList []string  `json:"productIdList,omitempty" bson:"product_id_list,omitempty"`
	Discount      int       `json:"discount,omitempty" bson:"discount,omitempty"`
	TimeStarted   *time.Time `json:"timeStarted" bson:"time_started,omitempty"`
	TimeExpired   *time.Time `json:"timeExpired" bson:"time_expired,omitempty"`
	ProductList   []Product `json:"productList,omitempty" bson:"-"`

	Created_At time.Time `json:"created_at" bson:"created_at"`
	Updated_At time.Time `json:"updated_at" bson:"updated_at"`
}
