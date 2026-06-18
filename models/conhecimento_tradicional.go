package models

import (
	"time"

	"gorm.io/gorm"
)

type ConhecimentoTradicional struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UsuarioID uint           `gorm:"not null" json:"usuario_id"`
	SementeID uint           `gorm:"not null" json:"semente_id"`
	Titulo    string         `gorm:"size:200;not null" json:"titulo"`
	Conteudo  string         `gorm:"type:text;not null" json:"conteudo"`
	Origem    string         `gorm:"size:200" json:"origem"` // comunidade, povo indígena, etc.
	Tecnica   string         `gorm:"size:200" json:"tecnica"`
	Validado  bool           `gorm:"default:false" json:"validado"`
	Curtidas  int            `gorm:"default:0" json:"curtidas"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relacionamentos
	Usuario Usuario `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
	Semente Semente `gorm:"foreignKey:SementeID" json:"semente,omitempty"`
}
