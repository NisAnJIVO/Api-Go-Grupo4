package handlers

import (
	"context"
	"net/http"
	"time"
	"trabajador-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Clave secreta para firmar los tokens — en producción va en variable de entorno
var claveSecreta = []byte("mi_clave_super_secreta_2024")

// ColeccionUsuarios apunta a la colección de usuarios en MongoDB
var ColeccionUsuarios *mongo.Collection

// Login godoc
// @Summary Iniciar sesión y obtener token JWT
// @Tags auth
// @Param credenciales body models.LoginRequest true "Email y contraseña"
// @Success 200 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var credenciales models.LoginRequest
	if err := c.ShouldBindJSON(&credenciales); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Buscar usuario en MongoDB
	var usuarioEncontrado models.Usuario
	err := ColeccionUsuarios.FindOne(
		context.Background(),
		bson.M{
			"email":    credenciales.Email,
			"password": credenciales.Password, // en producción usar bcrypt
		},
	).Decode(&usuarioEncontrado)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Generar el token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    usuarioEncontrado.ID.Hex(),
		"email": usuarioEncontrado.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // expira en 24 horas
	})

	tokenString, err := token.SignedString(claveSecreta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
