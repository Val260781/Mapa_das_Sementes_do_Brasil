package models

import (
	"time"

	"gorm.io/gorm"
)

type RoleUsuario string

const (
	RoleArtesao      RoleUsuario = "artesao"
	RolePesquisador  RoleUsuario = "pesquisador"
	RoleEspecialista RoleUsuario = "especialista"
	RoleAdmin        RoleUsuario = "admin"
)

type Usuario struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Nome          string         `gorm:"size:150;not null" json:"nome"`
	Email         string         `gorm:"size:200;unique;not null" json:"email"`
	Senha         string         `gorm:"size:255;not null" json:"-"`
	Role          RoleUsuario    `gorm:"size:20;default:artesao" json:"role"`
	Especialidade string         `gorm:"size:150" json:"especialidade"`        // Ex: Engenheiro Agrônomo
	RegistroProf  string         `gorm:"size:50" json:"registro_profissional"` // Ex: CREA, CRBio
	Estado        string         `gorm:"size:2" json:"estado"`
	Municipio     string         `gorm:"size:150" json:"municipio"`
	Bio           string         `gorm:"type:text" json:"bio"`
	AvatarURL     string         `gorm:"size:500" json:"avatar_url"`
	Ativo         bool           `gorm:"default:true" json:"ativo"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relacionamentos
	Registros         []RegistroSemente         `gorm:"foreignKey:UsuarioID" json:"registros,omitempty"`
	ConhecimentosTrad []ConhecimentoTradicional `gorm:"foreignKey:UsuarioID" json:"conhecimentos,omitempty"`
}
