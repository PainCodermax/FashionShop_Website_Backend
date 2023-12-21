package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	UserCode      string  `json:"userCode,omitempty" bson:"user_code,omitempty"`
	First_Name    *string `json:"first_name,omitempty" validate:"required,min=2,max=30"`
	Last_Name     *string `json:"last_name,omitempty"  validate:"required,min=2,max=30"`
	Password      string  `json:"password,omitempty"`
	Email         *string `json:"email,omitempty"`
	Phone         *string `json:"phone,omitempty"      validate:"required"`
	Token         *string `json:"token,omitempty"`
	Refresh_Token *string `json:"refresh_token,omitempty"`

	//address
	Province string `json:"province,omitempty"`
	District string `json:"district,omitempty"`
	Ward     string `json:"ward,omitempty"`

	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updtaed_at"`

	User_ID         string    `json:"user_id,omitempty"`
	Cart            Cart      `json:"cart,omitempty" bson:"cart,omitempty"`
	Address_Details []Address `json:"address,omitempty" bson:"address"`
	Order_Status    []Order   `json:"orders,omitempty" bson:"orders"`
	IsAdmin         bool      `json:"isAdmin,omitempty" bson:"is_admin,omitempty"`
	VerifyCode      string    `json:"verifyCode,omitempty" bson:"verify_code,omitempty"`
	IsVerified      bool      `json:"isVerified,omitempty" bson:"is_verified,omitempty"`
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

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_list"  bson:"order_list"`
	Orderered_At   time.Time          `json:"ordered_on"  bson:"ordered_on"`
	Price          int                `json:"total_price" bson:"total_price"`
	Discount       *int               `json:"discount"    bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod"     bson:"cod"`
}
