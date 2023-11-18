package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Product_ID   string             `json:"productId,omitempty" bson:"product_id,omitempty"`
	ProductName  *string            `json:"productName,omitempty" bson:"product_name,omitempty"`
	Price        *int               `json:"price,omitempty" bson:"price,omitempty"`
	Detail       *string            `json:"detail,omitempty" bson:"detail,omitempty"`
	ListImage    []string           `json:"listImage,omitempty" bson:"list_image,omitempty"`
	Quantity     *int               `json:"quantity,omitempty" bson:"quantity,omitempty"`
	CategoryID   string             `json:"categoryID,omitempty" bson:"category_id,omitempty"`
	Gender       string             `json:"gender,omitempty" bson:"gender,omitempty"`
	CategoryMame string             `json:"categoryName,omitempty"`
}
