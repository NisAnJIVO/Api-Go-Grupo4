package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Cliente *mongo.Client

func ConectarDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017" // fallback para desarrollo local
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error conectando a MongoDB:", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB no responde:", err)
	}

	log.Println("Conectado a MongoDB en:", uri)
	Cliente = client

	SembrarDatosIniciales(client)

	return client
}

// SembrarDatosIniciales crea un usuario y trabajadores de prueba si la base de datos está vacía
func SembrarDatosIniciales(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := client.Database("examen")
	usuariosCol := db.Collection("usuarios")
	trabajadoresCol := db.Collection("trabajadores")

	countUsuarios, _ := usuariosCol.CountDocuments(ctx, bson.M{})
	if countUsuarios == 0 {
		_, err := usuariosCol.InsertOne(ctx, bson.M{
			"email":    "ana@empresa.com",
			"password": "1234",
		})
		if err == nil {
			log.Println("Usuario inicial creado: email: ana@empresa.com | password: 1234")
		}
	}

	countTrabajadores, _ := trabajadoresCol.CountDocuments(ctx, bson.M{})
	if countTrabajadores == 0 {
		_, err := trabajadoresCol.InsertOne(ctx, bson.M{
			"nombre":       "Isaac",
			"apellido":     "Vargas",
			"ci":           "12642012",
			"cargo":        "Backend Developer Go",
			"departamento": "Tecnología",
		})
		if err == nil {
			log.Println("🌱 Trabajador inicial creado de ejemplo")
		}
	}
}
