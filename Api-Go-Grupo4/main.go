package main

import (
	"trabajador-api/db"
	"trabajador-api/handlers"
	"trabajador-api/middleware"

	_ "trabajador-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Trabajador API
// @version         1.0
// @host            localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cliente := db.ConectarDB()

	// Inyectamos las colecciones en los handlers
	handlers.Coleccion = cliente.Database("examen").Collection("trabajadores")
	handlers.ColeccionUsuarios = cliente.Database("examen").Collection("usuarios")

	r := gin.Default()

	// ── Ruta pública — no requiere token ──────────────────────────────────
	r.POST("/login", handlers.Login)

	// ── Rutas protegidas — pasan primero por el middleware JWT ────────────
	protegidas := r.Group("/")
	protegidas.Use(middleware.VerificarJWT)
	{
		// REST
		protegidas.GET("/trabajadores", handlers.ObtenerTodos)
		protegidas.POST("/trabajadores", handlers.Crear)
		protegidas.PUT("/trabajadores/:id", handlers.Actualizar)
		protegidas.DELETE("/trabajadores/:id", handlers.Eliminar)

		// GraphQL
		protegidas.POST("/graphql", handlers.HandleGraphQL)
	}

	// Swagger — público para facilitar las pruebas
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
