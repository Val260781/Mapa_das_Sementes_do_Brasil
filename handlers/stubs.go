package handlers

import (
	"time"

	"mapa-sementes-brasil/config"
	"mapa-sementes-brasil/database"
	"mapa-sementes-brasil/models"
	"mapa-sementes-brasil/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CadastroInput struct {
	// --- Obrigatórios ---
	NomeCompleto string `json:"nome_completo" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Senha        string `json:"senha" binding:"required"`
	Telefone     string `json:"telefone" binding:"required"`

	// --- Dados pessoais ---
	DataNascimento string `json:"data_nascimento"`
	CPF            string `json:"cpf"`

	// --- Localização ---
	Estado    string `json:"estado"`
	Municipio string `json:"municipio"`
	Bairro    string `json:"bairro"`

	// --- Atuação ---
	Role            string `json:"role"`
	Profissao       string `json:"profissao"`
	Bioma           string `json:"bioma"`
	LocalEncontro   string `json:"local_encontro"`
	Experiencia     string `json:"experiencia"`
	TempoAtuacao    string `json:"tempo_atuacao"`
	TiposArtesanato string `json:"tipos_artesanato"`

	// --- Especialista ---
	Especialidade      string `json:"especialidade"`
	RegistroProf       string `json:"registro_profissional"`
	InstituicaoVinculo string `json:"instituicao_vinculo"`

	// --- Contato e perfil ---
	RedesSociais string `json:"redes_sociais"`
	SiteURL      string `json:"site_url"`
	Bio          string `json:"bio"`
}

type LoginInput struct {
	Email string `json:"email" binding:"required"`
	Senha string `json:"senha" binding:"required"`
}

func Cadastro(c *gin.Context) {
	var input CadastroInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, "Dados inválidos: "+err.Error())
		return
	}

	if !utils.IsEmailValido(input.Email) {
		utils.Error(c, 400, "E-mail inválido")
		return
	}

	if !utils.IsSenhaValida(input.Senha) {
		utils.Error(c, 400, "A senha deve ter no mínimo 6 caracteres")
		return
	}

	var existe models.Usuario
	if result := database.DB.Where("email = ?", input.Email).First(&existe); result.Error == nil {
		utils.Error(c, 409, "E-mail já cadastrado")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Senha), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, 500, "Erro ao processar senha")
		return
	}

	role := models.RoleArtesao
	if input.Role == "pesquisador" {
		role = models.RolePesquisador
	} else if input.Role == "especialista" {
		role = models.RoleEspecialista
	}

	var dataNasc time.Time
	if input.DataNascimento != "" {
		dataNasc, err = time.Parse("2006-01-02", input.DataNascimento)
		if err != nil {
			utils.Error(c, 400, "Formato de data inválido. Use: AAAA-MM-DD (ex: 1990-05-20)")
			return
		}
	}

	usuario := models.Usuario{
		Email:              input.Email,
		Senha:              string(hash),
		Role:               role,
		Ativo:              true,
		NomeCompleto:       input.NomeCompleto,
		DataNascimento:     dataNasc,
		Telefone:           input.Telefone,
		CPF:                input.CPF,
		Estado:             input.Estado,
		Municipio:          input.Municipio,
		Bairro:             input.Bairro,
		Profissao:          input.Profissao,
		Bioma:              input.Bioma,
		LocalEncontro:      input.LocalEncontro,
		Experiencia:        input.Experiencia,
		TempoAtuacao:       input.TempoAtuacao,
		TiposArtesanato:    input.TiposArtesanato,
		Especialidade:      input.Especialidade,
		RegistroProf:       input.RegistroProf,
		InstituicaoVinculo: input.InstituicaoVinculo,
		RedesSociais:       input.RedesSociais,
		SiteURL:            input.SiteURL,
		Bio:                input.Bio,
	}

	if err := database.DB.Create(&usuario).Error; err != nil {
		utils.Error(c, 500, "Erro ao criar usuário")
		return
	}

	utils.Success(c, 201, "Usuário cadastrado com sucesso!", gin.H{
		"id":            usuario.ID,
		"nome_completo": usuario.NomeCompleto,
		"email":         usuario.Email,
		"role":          usuario.Role,
		"municipio":     usuario.Municipio,
		"estado":        usuario.Estado,
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, "Dados inválidos")
		return
	}

	var usuario models.Usuario
	if err := database.DB.Where("email = ? AND ativo = true", input.Email).First(&usuario).Error; err != nil {
		utils.Error(c, 401, "E-mail ou senha incorretos")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(input.Senha)); err != nil {
		utils.Error(c, 401, "E-mail ou senha incorretos")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usuario_id": usuario.ID,
		"role":       string(usuario.Role),
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		utils.Error(c, 500, "Erro ao gerar token")
		return
	}

	utils.Success(c, 200, "Login realizado com sucesso!", gin.H{
		"token": tokenStr,
		"usuario": gin.H{
			"id":            usuario.ID,
			"nome_completo": usuario.NomeCompleto,
			"email":         usuario.Email,
			"role":          usuario.Role,
			"especialidade": usuario.Especialidade,
			"municipio":     usuario.Municipio,
			"estado":        usuario.Estado,
		},
	})
}
