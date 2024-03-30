package models

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
	AddressID  string `json:"addressID,omitempty"`
	UserID     string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Street     string `json:"street"`
	ProvinceID int    `json:"ProvinceID,omitempty"`
	DistrictID int    `json:"DistrictID,omitempty"`
	WardID     int    `json:"WardCode,omitempty"`
}