package models

import (
	"time"

	"gorm.io/gorm"
)

type Semente struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	EspecieID      uint           `gorm:"not null" json:"especie_id"`
	Nome           string         `gorm:"size:200;not null" json:"nome"`
	Descricao      string         `gorm:"type:text" json:"descricao"`
	Cor            string         `gorm:"size:100" json:"cor"`
	Tamanho        string         `gorm:"size:50" json:"tamanho"` // pequena, média, grande
	Textura        string         `gorm:"size:100" json:"textura"`
	UsosArtesanais string         `gorm:"type:text" json:"usos_artesanais"`
	Tecnicas       string         `gorm:"type:text" json:"tecnicas_artesanato"`
	Validada       bool           `gorm:"default:false" json:"validada"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relacionamentos
	Especie   Especie           `gorm:"foreignKey:EspecieID" json:"especie,omitempty"`
	Registros []RegistroSemente `gorm:"foreignKey:SementeID" json:"registros,omitempty"`
}
