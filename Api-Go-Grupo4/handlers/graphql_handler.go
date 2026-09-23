package handlers

import (
	"context"
	"trabajador-api/models"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Definimos el objeto Trabajador para GraphQL
var trabajadorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Trabajador",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.String},
		"nombre":       &graphql.Field{Type: graphql.String},
		"apellido":     &graphql.Field{Type: graphql.String},
		"cargo":        &graphql.Field{Type: graphql.String},
		"departamento": &graphql.Field{Type: graphql.String},
	},
})

// Definimos la Query principal
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"obtenerTrabajadores": &graphql.Field{
			Type: graphql.NewList(trabajadorType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var lista []models.Trabajador
				cursor, _ := Coleccion.Find(context.Background(), bson.M{})
				cursor.All(context.Background(), &lista)
				return lista, nil
			},
		},
	},
})

// Mutaciones: crear, actualizar y eliminar
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{

		// ── CREAR ──────────────────────────────────────────────────────────
		"crearTrabajador": &graphql.Field{
			Type:        trabajadorType,
			Description: "Crea un nuevo trabajador",
			Args: graphql.FieldConfigArgument{
				"nombre":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"apellido":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"cargo":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"departamento": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				t := models.Trabajador{
					ID:           primitive.NewObjectID(),
					Nombre:       p.Args["nombre"].(string),
					Apellido:     p.Args["apellido"].(string),
					Cargo:        p.Args["cargo"].(string),
					Departamento: p.Args["departamento"].(string),
				}
				Coleccion.InsertOne(context.Background(), t)
				return t, nil
			},
		},

		// ── ACTUALIZAR ─────────────────────────────────────────────────────
		"actualizarTrabajador": &graphql.Field{
			Type:        trabajadorType,
			Description: "Actualiza los campos de un trabajador existente por ID",
			Args: graphql.FieldConfigArgument{
				"id":           &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"nombre":       &graphql.ArgumentConfig{Type: graphql.String},
				"apellido":     &graphql.ArgumentConfig{Type: graphql.String},
				"cargo":        &graphql.ArgumentConfig{Type: graphql.String},
				"departamento": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				oid, err := primitive.ObjectIDFromHex(p.Args["id"].(string))
				if err != nil {
					return nil, err
				}

				// Construimos el $set solo con los campos que llegaron
				campos := bson.M{}
				if v, ok := p.Args["nombre"].(string); ok {
					campos["nombre"] = v
				}
				if v, ok := p.Args["apellido"].(string); ok {
					campos["apellido"] = v
				}
				if v, ok := p.Args["cargo"].(string); ok {
					campos["cargo"] = v
				}
				if v, ok := p.Args["departamento"].(string); ok {
					campos["departamento"] = v
				}

				filtro := bson.M{"_id": oid}
				Coleccion.UpdateOne(context.Background(), filtro, bson.M{"$set": campos})

				// Devolvemos el documento actualizado
				var actualizado models.Trabajador
				Coleccion.FindOne(context.Background(), filtro).Decode(&actualizado)
				return actualizado, nil
			},
		},

		// ── ELIMINAR ───────────────────────────────────────────────────────
		"eliminarTrabajador": &graphql.Field{
			Type:        graphql.String,
			Description: "Elimina un trabajador por ID y devuelve un mensaje de confirmación",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				oid, err := primitive.ObjectIDFromHex(p.Args["id"].(string))
				if err != nil {
					return nil, err
				}
				Coleccion.DeleteOne(context.Background(), bson.M{"_id": oid})
				return "Trabajador eliminado correctamente", nil
			},
		},
	},
})

// Esquema Global — ahora incluye Query y Mutation
var TrabajadorSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

// HandleGraphQL es el controlador que Gin usará
func HandleGraphQL(c *gin.Context) {
	var requestBody struct {
		Query string `json:"query"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Cuerpo de petición inválido"})
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:        TrabajadorSchema,
		RequestString: requestBody.Query,
	})

	c.JSON(200, result)
}
