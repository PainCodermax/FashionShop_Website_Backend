package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	UserCode      string  `json:"userCode,omitempty" bson:"user_code,omitempty"`
	First_Name    *string `json:"first_name,omitempty" bson:"first_name,omitempty" validate:"required,min=2,max=30"`
	Last_Name     *string `json:"last_name,omitempty"  bson:"last_name,omitempty" validate:"required,min=2,max=30"`
	Password      string  `json:"password,omitempty" bson:"password,omitempty"`
	Email         *string `json:"email,omitempty" bson:"email,omitempty"`
	Phone         *string `json:"phone,omitempty" bson:"phone,omitempty"     validate:"required"`
	Token         *string `json:"token,omitempty" bson:"token,omitempty"`
	Refresh_Token *string `json:"refresh_token,omitempty" bson:"refresh_token"`

	//address
	Province string `json:"province,omitempty" bson:"province,omitempty"`
	District string `json:"district,omitempty" bson:"district,omitempty"`
	Ward     string `json:"ward,omitempty" bson:"ward,omitempty"`

	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updtaed_at"`

	User_ID    string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Cart       Cart   `json:"cart,omitempty" bson:"cart,omitempty"`
	IsAdmin    bool   `json:"isAdmin,omitempty" bson:"is_admin,omitempty"`
	VerifyCode string `json:"verifyCode,omitempty" bson:"verify_code,omitempty"`
	IsVerified bool   `json:"isVerified,omitempty" bson:"is_verified,omitempty"`
	IsActive   bool   `json:"isActive,omitempty" bson:"is_active,omitempty"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        int                `json:"price"  bson:"price"`
	Rating       *uint              `json:"rating" bson:"rating"`
	Image        *string            `json:"image"  bson:"image"`
}

type Address struct {
	Address_id primitive.ObjectID `bson:"_id"`
	House      *string            `json:"house_name" bson:"house_name"`
	Street     *string            `json:"street_name" bson:"street_name"`
	City       *string            `json:"city_name" bson:"city_name"`
	Pincode    *string            `json:"pin_code" bson:"pin_code"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod"     bson:"cod"`
}
