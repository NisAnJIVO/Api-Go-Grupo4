package handlers

import (
	"context"
	"net/http"
	"time"
	"trabajador-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Coleccion *mongo.Collection

// ObtenerTodos godoc
// @Summary Listar todos los trabajadores
// @Tags trabajadores
// @Security BearerAuth
// @Success 200 {array} models.Trabajador
// @Router /trabajadores [get]
func ObtenerTodos(c *gin.Context) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var trabajadores []models.Trabajador
	cursor, _ := Coleccion.Find(ctx, bson.M{})
	cursor.All(ctx, &trabajadores)
	c.JSON(http.StatusOK, trabajadores)
}

// Crear godoc
// @Summary Crear trabajador
// @Tags trabajadores
// @Security BearerAuth
// @Param trabajador body models.Trabajador true "Datos"
// @Success 201 {object} models.Trabajador
// @Router /trabajadores [post]
func Crear(c *gin.Context) {
	var t models.Trabajador
	c.ShouldBindJSON(&t)
	t.ID = primitive.NewObjectID()
	Coleccion.InsertOne(context.Background(), t)
	c.JSON(http.StatusCreated, t)
}

// Eliminar godoc
// @Summary Eliminar trabajador por ID
// @Tags trabajadores
// @Security BearerAuth
// @Param id path string true "ID del trabajador"
// @Success 200 {string} string "Eliminado"
// @Router /trabajadores/{id} [delete]
func Eliminar(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	Coleccion.DeleteOne(context.Background(), bson.M{"_id": id})
	c.JSON(http.StatusOK, gin.H{"mensaje": "Trabajador eliminado"})
}

// Actualizar godoc
// @Summary Actualizar trabajador
// @Tags trabajadores
// @Security BearerAuth
// @Param id path string true "ID"
// @Param trabajador body models.Trabajador true "Datos"
// @Router /trabajadores/{id} [put]
func Actualizar(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var t models.Trabajador
	c.ShouldBindJSON(&t)
	Coleccion.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": t})
	c.JSON(http.StatusOK, t)
}
