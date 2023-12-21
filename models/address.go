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