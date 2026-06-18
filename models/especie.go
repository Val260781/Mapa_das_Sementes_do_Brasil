package models

import (
	"time"

	"gorm.io/gorm"
)

type Especie struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	NomeCientifico    string         `gorm:"size:200;unique;not null" json:"nome_cientifico"`
	NomePopular       string         `gorm:"size:200" json:"nome_popular"`
	Familia           string         `gorm:"size:150" json:"familia"`
	Bioma             string         `gorm:"size:100" json:"bioma"` // Cerrado, Mata Atlântica, Amazônia...
	Descricao         string         `gorm:"type:text" json:"descricao"`
	UsosArtesanais    string         `gorm:"type:text" json:"usos_artesanais"`
	Epoca             string         `gorm:"size:100" json:"epoca_colheita"`
	StatusConservacao string         `gorm:"size:50" json:"status_conservacao"` // LC, NT, VU, EN, CR
	ImagemURL         string         `gorm:"size:500" json:"imagem_url"`
	Validada          bool           `gorm:"default:false" json:"validada"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relacionamentos
	Sementes []Semente `gorm:"foreignKey:EspecieID" json:"sementes,omitempty"`
}
