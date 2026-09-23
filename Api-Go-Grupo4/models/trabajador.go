package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Trabajador struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" example:"65f1c..."`
	Nombre       string             `json:"nombre" bson:"nombre" example:"Karen"`
	Apellido     string             `json:"apellido" bson:"apellido" example:"Amurrio"`
	CI           string             `json:"ci" bson:"ci" example:"1234567"`
	Cargo        string             `json:"cargo" bson:"cargo" example:"Developer"`
	Departamento string             `json:"departamento" bson:"departamento" example:"Sistemas"`
}
