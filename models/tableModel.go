package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Capacity     *int               `json:"capacity" validate:"required"`
	Table_Number *int               `json:"table_number" validate:"required"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
	Table_id     string             `json:"table_id"`
}
