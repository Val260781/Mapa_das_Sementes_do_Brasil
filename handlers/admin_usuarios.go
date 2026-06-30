package handlers

import (
	"strconv"

	"mapa-sementes-brasil/database"
	"mapa-sementes-brasil/models"
	"mapa-sementes-brasil/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ─────────────────────────────────────────
// EDITAR USUÁRIO (ADMIN)
// PUT /api/admin/usuarios/:id
// Permite corrigir e-mail, resetar senha, mudar role, etc.
// de QUALQUER usuário — somente admin.
// ─────────────────────────────────────────
func AdminEditarUsuario(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var usuario models.Usuario
	if err := database.DB.First(&usuario, id).Error; err != nil {
		utils.Error(c, 404, "Usuário não encontrado")
		return
	}

	var input struct {
		Email          string `json:"email"`
		NovaSenha      string `json:"nova_senha"` // se vazio, mantém a senha atual
		NomeCompleto   string `json:"nome_completo"`
		Telefone       string `json:"telefone"`
		Role           string `json:"role"`
		Estado         string `json:"estado"`
		Municipio      string `json:"municipio"`
		Ativo          *bool  `json:"ativo"` // ponteiro para diferenciar "não enviado" de "false"
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, "Dados inválidos: "+err.Error())
		return
	}

	updates := map[string]interface{}{}

	// E-mail: verifica se o novo e-mail já não pertence a outro usuário
	if input.Email != "" && input.Email != usuario.Email {
		if !utils.IsEmailValido(input.Email) {
			utils.Error(c, 400, "E-mail inválido")
			return
		}
		var existe models.Usuario
		if err := database.DB.Where("email = ? AND id <> ?", input.Email, id).First(&existe).Error; err == nil {
			utils.Error(c, 409, "Este e-mail já está em uso por outro usuário")
			return
		}
		updates["email"] = input.Email
	}

	// Reset de senha (opcional)
	if input.NovaSenha != "" {
		if !utils.IsSenhaValida(input.NovaSenha) {
			utils.Error(c, 400, "A nova senha deve ter no mínimo 6 caracteres")
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(input.NovaSenha), bcrypt.DefaultCost)
		if err != nil {
			utils.Error(c, 500, "Erro ao processar nova senha")
			return
		}
		updates["senha"] = string(hash)
	}

	if input.NomeCompleto != "" {
		updates["nome_completo"] = input.NomeCompleto
	}
	if input.Telefone != "" {
		updates["telefone"] = input.Telefone
	}
	if input.Estado != "" {
		updates["estado"] = input.Estado
	}
	if input.Municipio != "" {
		updates["municipio"] = input.Municipio
	}
	if input.Role != "" {
		updates["role"] = input.Role
	}
	if input.Ativo != nil {
		updates["ativo"] = *input.Ativo
	}

	if len(updates) == 0 {
		utils.Error(c, 400, "Nenhum campo para atualizar foi informado")
		return
	}

	if err := database.DB.Model(&usuario).Updates(updates).Error; err != nil {
		utils.Error(c, 500, "Erro ao atualizar usuário")
		return
	}

	database.DB.First(&usuario, id)
	utils.Success(c, 200, "Usuário atualizado com sucesso!", gin.H{
		"id":            usuario.ID,
		"nome_completo": usuario.NomeCompleto,
		"email":         usuario.Email,
		"role":          usuario.Role,
		"ativo":         usuario.Ativo,
		"municipio":     usuario.Municipio,
		"estado":        usuario.Estado,
	})
}

// ─────────────────────────────────────────
// EXCLUIR USUÁRIO (ADMIN)
// DELETE /api/admin/usuarios/:id
// Exclusão definitiva (hard delete) — libera o e-mail para novo cadastro.
// Somente admin.
// ─────────────────────────────────────────
func AdminExcluirUsuario(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	adminID := c.GetUint("usuario_id")
	if uint(id) == adminID {
		utils.Error(c, 400, "Você não pode excluir sua própria conta de administrador por aqui")
		return
	}

	var usuario models.Usuario
	if err := database.DB.First(&usuario, id).Error; err != nil {
		utils.Error(c, 404, "Usuário não encontrado")
		return
	}

	// Unscoped() faz exclusão definitiva (hard delete), ignorando o soft-delete
	// padrão do GORM. Isso é necessário para liberar o e-mail (índice unique)
	// para um novo cadastro.
	if err := database.DB.Unscoped().Delete(&usuario).Error; err != nil {
		utils.Error(c, 500, "Erro ao excluir usuário")
		return
	}

	utils.Success(c, 200, "Usuário excluído com sucesso. O e-mail está liberado para novo cadastro.", nil)
}
