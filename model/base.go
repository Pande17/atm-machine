package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Base struct {
	ID        primitive.ObjectID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
