package main

import (
	"log"

	"mapa-sementes-brasil/config"
	"mapa-sementes-brasil/database"
	"mapa-sementes-brasil/handlers"
	"mapa-sementes-brasil/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	database.Connect()
	database.RunMigrations()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "🌱 Mapa das Sementes do Brasil — API rodando!",
		})
	})

	// Rotas públicas
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/cadastro", handlers.Cadastro)
			auth.POST("/login", handlers.Login)
		}
	}

	// Rotas protegidas
	protegido := api.Group("/")
	protegido.Use(middleware.AuthRequired())
	{
		protegido.GET("/perfil", handlers.MeuPerfil)
	}

	port := config.App.Port
	log.Printf("🚀 Servidor iniciando na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
