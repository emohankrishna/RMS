package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	First_name   string             `json:"first_name" validate:"required,min=2,max=100"`
	Last_name    string             `json:"last_name" validate:"required,min=2,max=100"`
	Password     string             `json:"password" validate:"required,min=6"`
	Email        string             `json:"email" validate:"email,required"`
	Avatar       string             `json:"avatar"`
	Phone        string             `json:"phone" validate:"required"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refresh_token"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
	User_id      string             `json:"user_id"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=6"`
}
