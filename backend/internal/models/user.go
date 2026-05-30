package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username"      json:"username"`
	Password  string             `bson:"password"      json:"-"`
	Role      string             `bson:"role"          json:"role"`
	CreatedAt time.Time          `bson:"createdAt"     json:"createdAt"`
	LastLogin time.Time          `bson:"lastLogin"     json:"lastLogin"`
}