package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Admin struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	UserName *string            `json:"username" validate:"required,min=2,max=100"`
	Password *string            `json:"Password" validate:"required,min=6"`
}
