package models

import (
	"time"

	"gorm.io/gorm"
)

type StatusRegistro string

const (
	StatusPendente  StatusRegistro = "pendente"
	StatusAprovado  StatusRegistro = "aprovado"
	StatusRejeitado StatusRegistro = "rejeitado"
)

type RegistroSemente struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UsuarioID  uint           `gorm:"not null" json:"usuario_id"`
	SementeID  uint           `gorm:"not null" json:"semente_id"`
	Latitude   float64        `gorm:"not null" json:"latitude"`
	Longitude  float64        `gorm:"not null" json:"longitude"`
	Estado     string         `gorm:"size:2" json:"estado"`
	Municipio  string         `gorm:"size:150" json:"municipio"`
	Descricao  string         `gorm:"type:text" json:"descricao"`
	Quantidade int            `json:"quantidade"` // estimativa de sementes coletadas
	DataColeta time.Time      `json:"data_coleta"`
	Status     StatusRegistro `gorm:"size:20;default:pendente" json:"status"`
	Fotos      []FotoRegistro `gorm:"foreignKey:RegistroID" json:"fotos,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relacionamentos
	Usuario Usuario `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
	Semente Semente `gorm:"foreignKey:SementeID" json:"semente,omitempty"`
}

type FotoRegistro struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RegistroID uint      `gorm:"not null" json:"registro_id"`
	URL        string    `gorm:"size:500;not null" json:"url"`
	Legenda    string    `gorm:"size:300" json:"legenda"`
	CreatedAt  time.Time `json:"created_at"`
}
