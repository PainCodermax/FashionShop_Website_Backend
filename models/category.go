package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       *string            `json:"name"`
	Detail     *string            `json:"detail"`
	CategoryId string             `json:"categoryId,omitempty" bson:"category_id,omitempty"`
	Gender     string             `json:"gender,omitempty" bson:"gender,omitempty"`
	Created_At time.Time          `json:"created_at"`
	Updated_At time.Time          `json:"updtaed_at"`
}
