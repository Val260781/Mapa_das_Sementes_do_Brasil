package models

import (
	"time"

	"gorm.io/gorm"
)

type RoleUsuario string

const (
	RoleArtesao           RoleUsuario = "artesao"
	RolePesquisador       RoleUsuario = "pesquisador"
	RoleEspecialista      RoleUsuario = "especialista"
	RoleEstudante         RoleUsuario = "estudante"
	RoleAgenteTerritorial RoleUsuario = "agente_territorial"
	RoleAdmin             RoleUsuario = "admin"
)

type Usuario struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// --- Dados de acesso ---
	Email string      `gorm:"size:200;unique;not null" json:"email"`
	Senha string      `gorm:"size:255;not null" json:"-"`
	Role  RoleUsuario `gorm:"size:20;default:artesao" json:"role"`
	Ativo bool        `gorm:"default:true" json:"ativo"`

	// --- Dados pessoais ---
	NomeCompleto   string    `gorm:"size:200" json:"nome_completo"`
	DataNascimento time.Time `json:"data_nascimento"`
	Telefone       string    `gorm:"size:20" json:"telefone"` // Ex: (62) 99999-9999
	CPF            string    `gorm:"size:14" json:"cpf"`      // opcional, para validação futura

	// --- Localização ---
	Estado    string `gorm:"size:2" json:"estado"`
	Municipio string `gorm:"size:150" json:"municipio"`
	Bairro    string `gorm:"size:150" json:"bairro"`

	// --- Atuação no projeto ---
	// Profissao também é usado para o curso, no caso de role=estudante
	// (ex: "Engenharia Florestal - UFG"), evitando precisar de coluna nova.
	Profissao       string `gorm:"size:150" json:"profissao"`         // Ex: Artesã, Biólogo, Agricultor, ou Curso (se estudante)
	Bioma           string `gorm:"size:100" json:"bioma"`             // Ex: Cerrado, Amazônia, Caatinga
	LocalEncontro   string `gorm:"type:text" json:"local_encontro"`   // Descrição livre de onde encontra sementes
	Experiencia     string `gorm:"type:text" json:"experiencia"`      // Experiência com sementes/artesanato
	TempoAtuacao    string `gorm:"size:50" json:"tempo_atuacao"`      // Ex: "2 anos", "Desde 1995"
	TiposArtesanato string `gorm:"type:text" json:"tipos_artesanato"` // Tipos de artesanato que produz

	// --- Especialista (preenchido só para role=especialista) ---
	Especialidade      string `gorm:"size:150" json:"especialidade"`        // Ex: Engenheiro Florestal
	RegistroProf       string `gorm:"size:50" json:"registro_profissional"` // Ex: CREA-GO 123456
	InstituicaoVinculo string `gorm:"size:200" json:"instituicao_vinculo"`  // Ex: UFG, Embrapa

	// --- Contato e redes ---
	RedesSociais string `gorm:"size:300" json:"redes_sociais"` // Ex: @usuario_instagram
	SiteURL      string `gorm:"size:300" json:"site_url"`

	// --- Mídia ---
	AvatarURL string `gorm:"size:500" json:"avatar_url"`
	Bio       string `gorm:"type:text" json:"bio"` // Apresentação curta

	// --- Relacionamentos ---
	Registros         []RegistroSemente         `gorm:"foreignKey:UsuarioID" json:"registros,omitempty"`
	ConhecimentosTrad []ConhecimentoTradicional `gorm:"foreignKey:UsuarioID" json:"conhecimentos,omitempty"`
	Avaliacoes        []Avaliacao               `gorm:"foreignKey:EspecialistaID" json:"avaliacoes,omitempty"`
}
