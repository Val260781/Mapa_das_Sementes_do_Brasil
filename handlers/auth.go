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
	Nome          string `json:"nome" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Senha         string `json:"senha" binding:"required"`
	Role          string `json:"role"`
	Especialidade string `json:"especialidade"`
	RegistroProf  string `json:"registro_profissional"`
	Estado        string `json:"estado"`
	Municipio     string `json:"municipio"`
	Bio           string `json:"bio"`
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

	// Verificar se e-mail já existe
	var existe models.Usuario
	if result := database.DB.Where("email = ?", input.Email).First(&existe); result.Error == nil {
		utils.Error(c, 409, "E-mail já cadastrado")
		return
	}

	// Hash da senha
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Senha), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, 500, "Erro ao processar senha")
		return
	}

	// Definir role
	role := models.RoleArtesao
	if input.Role == "pesquisador" {
		role = models.RolePesquisador
	} else if input.Role == "especialista" {
		role = models.RoleEspecialista
	}

	usuario := models.Usuario{
		Nome:          input.Nome,
		Email:         input.Email,
		Senha:         string(hash),
		Role:          role,
		Especialidade: input.Especialidade,
		RegistroProf:  input.RegistroProf,
		Estado:        input.Estado,
		Municipio:     input.Municipio,
		Bio:           input.Bio,
	}

	if err := database.DB.Create(&usuario).Error; err != nil {
		utils.Error(c, 500, "Erro ao criar usuário")
		return
	}

	utils.Success(c, 201, "Usuário cadastrado com sucesso!", gin.H{
		"id":    usuario.ID,
		"nome":  usuario.Nome,
		"email": usuario.Email,
		"role":  usuario.Role,
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

	// Gerar JWT
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
			"nome":          usuario.Nome,
			"email":         usuario.Email,
			"role":          usuario.Role,
			"especialidade": usuario.Especialidade,
		},
	})
}

