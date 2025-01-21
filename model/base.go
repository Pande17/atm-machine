package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Base struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id, omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt, omitempty"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt, omitempty"`
	DeletedAt *time.Time         `json:"deletedAt" bson:"deletedAt, omitempty"`
}
