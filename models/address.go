package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Province struct {
	ProvinceID   int    `json:"ProvinceID,omitempty"`
	ProvinceName string `json:"ProvinceName,omitempty"`
}

type District struct {
	DistrictID   int    `json:"DistrictID,omitempty"`
	DistrictName string `json:"DistrictName,omitempty"`
}

type Ward struct {
	DistrictID   int    `json:"WardCode,omitempty"`
	DistrictName string `json:"DistrictName,omitempty"`
}

type UserAddress struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	AddressID  string `json:"addressID,omitempty" bson:"address_id,omitempty"`
	UserID     string `json:"userId,omitempty" bson:"user_id,omitempty"`
	Street     string `json:"street" bson:"street,omitempty"`
	ProvinceID int    `json:"ProvinceID,omitempty" bson:"province_id,omitempty"`
	DistrictID int    `json:"DistrictID,omitempty" bson:"district_id,omitempty"`
	WardID     int    `json:"WardCode,omitempty" bson:"ward_id,omitempty"`
}
