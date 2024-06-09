package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Province struct {
	ProvinceID    int      `json:"ProvinceID,omitempty"`
	ProvinceName  string   `json:"ProvinceName,omitempty"`
	NameExtension []string `json:"NameExtension,omitempty"`
}

type District struct {
	DistrictID   int    `json:"DistrictID,omitempty"`
	DistrictName string `json:"DistrictName,omitempty"`
}

type Ward struct {
	WardID   string `json:"WardCode,omitempty"`
	WardName string `json:"WardName,omitempty"`
}

type UserAddress struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Phone       string `json:"phone,omitempty" bson:"phone,omitempty"`
	AddressID   string `json:"addressID,omitempty" bson:"address_id,omitempty"`
	UserID      string `json:"userId,omitempty" bson:"user_id,omitempty"`
	Street      string `json:"street,omitempty" bson:"street,omitempty"`
	ProvinceID  string `json:"ProvinceID,omitempty" bson:"province_id,omitempty"`
	DistrictID  string `json:"DistrictID,omitempty" bson:"district_id,omitempty"`
	WardID      string `json:"WardCode,omitempty" bson:"ward_id,omitempty"`
	IsDefault   bool   `json:"isDefault,omitempty" bson:"is_default,omitempty"`
	FullAddress string `json:"fullAddress,omitempty" bson:"full_address,omitempty"`
}
