package database

import (
	"log"

	"mapa-sementes-brasil/models"
)

func RunMigrations() {
	err := DB.AutoMigrate(
		&models.Usuario{},
		&models.Especie{},
		&models.Semente{},
		&models.RegistroSemente{},
		&models.ConhecimentoTradicional{},
	)

	if err != nil {
		log.Fatalf("Erro ao executar migrations: %v", err)
	}

	log.Println("✅ Migrations executadas com sucesso!")
}
