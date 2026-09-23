package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Usuario struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `bson:"email"         json:"email"    example:"ana@empresa.com"`
	Password string             `bson:"password"      json:"password" example:"1234"`
}
