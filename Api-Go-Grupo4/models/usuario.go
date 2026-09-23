package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Usuario struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string             `bson:"email"         json:"email"    example:"ana@empresa.com"`
	Password string             `bson:"password"      json:"password" example:"1234"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"ana@empresa.com"`
	Password string `json:"password" binding:"required" example:"1234"`
}
