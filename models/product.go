package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Product_ID     string             `json:"productId,omitempty" bson:"product_id,omitempty"`
	ProductName    *string            `json:"productName,omitempty" bson:"product_name,omitempty"`
	Price          *int               `json:"price,omitempty" bson:"price,omitempty"`
	Detail         *string            `json:"detail,omitempty" bson:"detail,omitempty"`
	ListImage      []string           `json:"listImage,omitempty" bson:"list_image,omitempty"`
	Quantity       *int               `json:"quantity,omitempty" bson:"quantity,omitempty"`
	CategoryID     string             `json:"categoryID,omitempty" bson:"category_id,omitempty"`
	Gender         string             `json:"gender,omitempty" bson:"gender,omitempty"`
	CategoryMame   string             `json:"categoryName,omitempty" bson:"categorymame,omitempty"`
	FlashSalePrice int                `json:"flashSalePrice,omitempty" bson:"-"`
	SalePrice      *int               `json:"salePrice,omitempty" bson:"-"`

	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updtaed_at"`
}

type Recommendation struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Product_id_1 string             `json:"productId1,omitempty" bson:"product_id_1,omitempty"`
	Product_id_2 string             `json:"productId2,omitempty" bson:"product_id_2,omitempty"`
	Weight       float64            `json:"weight,omitempty" bson:"weight,omitempty"`
}
