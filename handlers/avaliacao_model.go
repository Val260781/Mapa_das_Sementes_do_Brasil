package models

import (
	"time"

	"gorm.io/gorm"
)

type StatusAvaliacao string
type TipoAvaliacao string

const (
	AvaliacaoAprovada  StatusAvaliacao = "aprovada"
	AvaliacaoReprovada StatusAvaliacao = "reprovada"
	AvaliacaoPendente  StatusAvaliacao = "pendente"
)

const (
	TipoEspecie StatusAvaliacao = "especie"
	TipoSemente StatusAvaliacao = "semente"
)

type Avaliacao struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`

	// --- Quem avaliou ---
	EspecialistaID uint    `gorm:"not null" json:"especialista_id"`
	Especialista   Usuario `gorm:"foreignKey:EspecialistaID" json:"especialista,omitempty"`

	// --- O que foi avaliado ---
	// Apenas um dos dois será preenchido por avaliação
	EspecieID *uint    `json:"especie_id,omitempty"`
	Especie   *Especie `gorm:"foreignKey:EspecieID" json:"especie,omitempty"`

	SementeID *uint    `json:"semente_id,omitempty"`
	Semente   *Semente `gorm:"foreignKey:SementeID" json:"semente,omitempty"`

	// --- Resultado ---
	Status  StatusAvaliacao `gorm:"size:20;not null" json:"status"` // aprovada | reprovada
	Parecer string          `gorm:"type:text;not null" json:"parecer"` // obrigatório: motivo técnico

	// --- Observações extras (opcionais) ---
	RecomendacaoCorrecao string `gorm:"type:text" json:"recomendacao_correcao"` // o que deve ser corrigido
	FontesBibliograficas string `gorm:"type:text" json:"fontes_bibliograficas"`  // referências usadas
	DataAvaliacao        time.Time `json:"data_avaliacao"`
}
