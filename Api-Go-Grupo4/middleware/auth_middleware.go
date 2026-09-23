package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var claveSecreta = []byte("mi_clave_super_secreta_2024")

// VerificarJWT intercepta cada request y valida el token antes de continuar
func VerificarJWT(c *gin.Context) {
	// El token llega en el header: Authorization: Bearer <token>
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no proporcionado"})
		c.Abort() // detiene la cadena de handlers, no continúa
		return
	}

	// Separamos "Bearer" del token real
	partes := strings.SplitN(authHeader, " ", 2)
	if len(partes) != 2 || partes[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
		c.Abort()
		return
	}

	tokenString := partes[1]

	// Parseamos y validamos el token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Verificamos que el algoritmo sea el que esperamos
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return claveSecreta, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
		c.Abort()
		return
	}

	// Extraemos los claims y los guardamos en el contexto de Gin
	// para que los handlers puedan usarlos si los necesitan
	claims, _ := token.Claims.(jwt.MapClaims)
	c.Set("email", claims["email"])
	c.Set("id", claims["id"])

	c.Next() // todo bien, continúa al handler real
}
