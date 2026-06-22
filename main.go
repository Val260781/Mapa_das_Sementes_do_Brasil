package main

import (
	"log"
	"net/http"

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
	r.Static("/uploads", "./uploads")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "🌱 Mapa das Sementes do Brasil — API rodando!",
		})
	})

	api := r.Group("/api")

	// ─── Auth ───────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/cadastro", handlers.Cadastro)
		auth.POST("/login", handlers.Login)
	}

	// ─── Públicas ───────────────────────────────────
	api.GET("/especies", handlers.ListarEspecies)
	api.GET("/especies/:id", handlers.DetalheEspecie)
	api.GET("/sementes", handlers.ListarSementes)
	api.GET("/sementes/:id", handlers.DetalheSemente)
	api.GET("/sementes/:id/conhecimentos", handlers.ConhecimentosPorSemente)
	api.GET("/usuarios/:id", handlers.PerfilPublico)
	api.GET("/registros", handlers.ListarRegistros)
	api.GET("/registros/mapa", handlers.RegistrosParaMapa)
	api.GET("/registros/:id", handlers.DetalheRegistro)
	api.GET("/conhecimentos", handlers.ListarConhecimentos)
	api.GET("/conhecimentos/:id", handlers.DetalheConhecimento)

	// ─── Busca avançada (pública) ───────────────────
	busca := api.Group("/busca")
	{
		busca.GET("", handlers.BuscaGeral)
		busca.GET("/especies", handlers.BuscaEspecies)
		busca.GET("/sementes", handlers.BuscaSementes)
		busca.GET("/mapa", handlers.BuscaPorProximidade)
		busca.GET("/estado/:uf", handlers.BuscaPorEstado)
		busca.GET("/estatisticas", handlers.Estatisticas)
	}

	// ─── Protegidas ─────────────────────────────────
	protegido := api.Group("/")
	protegido.Use(middleware.AuthRequired())
	{
		// Perfil
		protegido.GET("/perfil", handlers.MeuPerfil)
		protegido.PUT("/perfil", handlers.EditarPerfil)
		protegido.PUT("/perfil/senha", handlers.TrocarSenha)
		protegido.POST("/perfil/avatar", handlers.UploadAvatar)
		protegido.GET("/perfil/contribuicoes", handlers.MinhasContribuicoes)
		protegido.DELETE("/perfil", handlers.DesativarConta)

		// Espécies
		protegido.POST("/especies", handlers.CriarEspecie)
		protegido.POST("/especies/:id/foto", handlers.UploadFotoEspecie)
		protegido.PUT("/especies/:id", handlers.EditarEspecie)
		protegido.DELETE("/especies/:id", handlers.DeletarEspecie)

		// Sementes
		protegido.POST("/sementes", handlers.CriarSemente)
		protegido.PUT("/sementes/:id", handlers.EditarSemente)
		protegido.DELETE("/sementes/:id", handlers.DeletarSemente)

		// Registros
		protegido.POST("/registros", handlers.CriarRegistro)
		protegido.PUT("/registros/:id", handlers.EditarRegistro)
		protegido.DELETE("/registros/:id", handlers.DeletarRegistro)
		protegido.POST("/registros/:id/fotos", handlers.UploadFotosRegistro)
		protegido.DELETE("/registros/:id/fotos/:foto_id", handlers.DeletarFotoRegistro)

		// Conhecimento Tradicional
		protegido.POST("/conhecimentos", handlers.CriarConhecimento)
		protegido.PUT("/conhecimentos/:id", handlers.EditarConhecimento)
		protegido.DELETE("/conhecimentos/:id", handlers.DeletarConhecimento)
		protegido.POST("/conhecimentos/:id/curtir", handlers.CurtirConhecimento)
		protegido.POST("/conhecimentos/:id/validar", handlers.ValidarConhecimento)

		// Avaliações
		protegido.POST("/avaliacoes/especie/:id", handlers.AvaliarEspecie)
		protegido.POST("/avaliacoes/semente/:id", handlers.AvaliarSemente)
		protegido.GET("/avaliacoes/pendentes", handlers.ListarPendentes)
		protegido.GET("/avaliacoes/historico", handlers.HistoricoAvaliacoes)

		// Admin
		admin := protegido.Group("/admin")
		admin.Use(middleware.RoleRequired("admin"))
		{
			admin.GET("/usuarios", handlers.ListarUsuarios)
		}
	}

	port := config.App.Port
	log.Printf("🚀 Servidor iniciando na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

