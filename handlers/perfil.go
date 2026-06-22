package handlers

import (
	"net/http"
	"time"

	"mapa-sementes-brasil/database"
	"mapa-sementes-brasil/models"
	"mapa-sementes-brasil/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ─────────────────────────────────────────
// VER MEU PERFIL
// GET /api/perfil
// ─────────────────────────────────────────
func MeuPerfil(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

	var usuario models.Usuario
	if err := database.DB.First(&usuario, usuarioID).Error; err != nil {
		utils.Error(c, 404, "Usuário não encontrado")
		return
	}

	utils.Success(c, 200, "Perfil encontrado", usuario)
}

// ─────────────────────────────────────────
// EDITAR MEU PERFIL
// PUT /api/perfil
// ─────────────────────────────────────────
func EditarPerfil(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

	var usuario models.Usuario
	if err := database.DB.First(&usuario, usuarioID).Error; err != nil {
		utils.Error(c, 404, "Usuário não encontrado")
		return
	}

	var input struct {
		// Dados pessoais
		NomeCompleto   string `json:"nome_completo"`
		DataNascimento string `json:"data_nascimento"` // AAAA-MM-DD
		Telefone       string `json:"telefone"`

		// Localização
		Estado    string `json:"estado"`
		Municipio string `json:"municipio"`
		Bairro    string `json:"bairro"`

		// Atuação
		Profissao       string `json:"profissao"`
		Bioma           string `json:"bioma"`
		LocalEncontro   string `json:"local_encontro"`
		Experiencia     string `json:"experiencia"`
		TempoAtuacao    string `json:"tempo_atuacao"`
		TiposArtesanato string `json:"tipos_artesanato"`

		// Especialista
		Especialidade      string `json:"especialidade"`
		RegistroProf       string `json:"registro_profissional"`
		InstituicaoVinculo string `json:"instituicao_vinculo"`

		// Contato
		RedesSociais string `json:"redes_sociais"`
		SiteURL      string `json:"site_url"`
		Bio          string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, "Dados inválidos: "+err.Error())
		return
	}

	updates := map[string]interface{}{}

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
	if input.Bairro != "" {
		updates["bairro"] = input.Bairro
	}
	if input.Profissao != "" {
		updates["profissao"] = input.Profissao
	}
	if input.Bioma != "" {
		updates["bioma"] = input.Bioma
	}
	if input.LocalEncontro != "" {
		updates["local_encontro"] = input.LocalEncontro
	}
	if input.Experiencia != "" {
		updates["experiencia"] = input.Experiencia
	}
	if input.TempoAtuacao != "" {
		updates["tempo_atuacao"] = input.TempoAtuacao
	}
	if input.TiposArtesanato != "" {
		updates["tipos_artesanato"] = input.TiposArtesanato
	}
	if input.Especialidade != "" {
		updates["especialidade"] = input.Especialidade
	}
	if input.RegistroProf != "" {
		updates["registro_prof"] = input.RegistroProf
	}
	if input.InstituicaoVinculo != "" {
		updates["instituicao_vinculo"] = input.InstituicaoVinculo
	}
	if input.RedesSociais != "" {
		updates["redes_sociais"] = input.RedesSociais
	}
	if input.SiteURL != "" {
		updates["site_url"] = input.SiteURL
	}
	if input.Bio != "" {
		updates["bio"] = input.Bio
	}

	// Converte data de nascimento se informada
	if input.DataNascimento != "" {
		dataNasc, err := time.Parse("2006-01-02", input.DataNascimento)
		if err != nil {
			utils.Error(c, 400, "Formato de data inválido. Use: AAAA-MM-DD (ex: 1990-05-20)")
			return
		}
		updates["data_nascimento"] = dataNasc
	}

	if len(updates) == 0 {
		utils.Error(c, 400, "Nenhum campo para atualizar foi informado")
		return
	}

	if err := database.DB.Model(&usuario).Updates(updates).Error; err != nil {
		utils.Error(c, 500, "Erro ao atualizar perfil")
		return
	}

	// Recarrega o usuário atualizado
	database.DB.First(&usuario, usuarioID)

	utils.Success(c, 200, "Perfil atualizado com sucesso!", usuario)
}

// ─────────────────────────────────────────
// TROCAR SENHA
// PUT /api/perfil/senha
// ─────────────────────────────────────────
func TrocarSenha(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

	var usuario models.Usuario
	if err := database.DB.First(&usuario, usuarioID).Error; err != nil {
		utils.Error(c, 404, "Usuário não encontrado")
		return
	}

	var input struct {
		SenhaAtual string `json:"senha_atual" binding:"required"`
		NovaSenha  string `json:"nova_senha" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, "Informe senha_atual e nova_senha")
		return
	}

	// Verifica senha atual
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(input.SenhaAtual)); err != nil {
		utils.Error(c, 401, "Senha atual incorreta")
		return
	}

	if !utils.IsSenhaValida(input.NovaSenha) {
		utils.Error(c, 400, "A nova senha deve ter no mínimo 6 caracteres")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.NovaSenha), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, 500, "Erro ao processar nova senha")
		return
	}

	if err := database.DB.Model(&usuario).Update("senha", string(hash)).Error; err != nil {
		utils.Error(c, 500, "Erro ao atualizar senha")
		return
	}

	utils.Success(c, 200, "Senha alterada com sucesso!", nil)
}

// ─────────────────────────────────────────
// UPLOAD DE AVATAR
// POST /api/perfil/avatar
// ─────────────────────────────────────────
func UploadAvatar(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

	var usuario models.Usuario
	if err := database.DB.First(&usuario, usuarioID).Error; err != nil {
		utils.Error(c, 404, "Usuário não encontrado")
		return
	}

	arquivo, err := c.FormFile("avatar")
	if err != nil {
		utils.Error(c, 400, "Nenhuma imagem enviada. Use o campo 'avatar'")
		return
	}

	urlAvatar, err := utils.SalvarImagem(c, arquivo, "avatares")
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// Remove avatar anterior se existir
	if usuario.AvatarURL != "" {
		utils.DeletarImagem(usuario.AvatarURL)
	}

	if err := database.DB.Model(&usuario).Update("avatar_url", urlAvatar).Error; err != nil {
		utils.Error(c, 500, "Erro ao salvar avatar")
		return
	}

	utils.Success(c, 200, "Avatar atualizado com sucesso!", gin.H{
		"avatar_url": urlAvatar,
	})
}

// ─────────────────────────────────────────
// MINHAS CONTRIBUIÇÕES
// GET /api/perfil/contribuicoes
// ─────────────────────────────────────────
func MinhasContribuicoes(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

	// Espécies cadastradas pelo usuário
	// (usamos CreatedBy se existir, senão buscamos por registros vinculados)
	var registros []models.RegistroSemente
	database.DB.
		Where("usuario_id = ?", usuarioID).
		Preload("Semente").
		Preload("Semente.Especie").
		Order("created_at DESC").
		Find(&registros)

	// Avaliações feitas (se for especialista)
	var avaliacoes []models.Avaliacao
	database.DB.
		Where("especialista_id = ?", usuarioID).
		Preload("Especie").
		Preload("Semente").
		Order("created_at DESC").
		Find(&avaliacoes)

	// Conhecimentos tradicionais cadastrados
	var conhecimentos []models.ConhecimentoTradicional
	database.DB.
		Where("usuario_id = ?", usuarioID).
		Preload("Semente").
		Order("created_at DESC").
		Find(&conhecimentos)

	utils.Success(c, 200, "Contribuições encontradas", gin.H{
		"registros_semente":          registros,
		"total_registros":            len(registros),
		"avaliacoes":                 avaliacoes,
		"total_avaliacoes":           len(avaliacoes),
		"conhecimentos_tradicionais": conhecimentos,
		"total_conhecimentos":        len(conhecimentos),
	})
}

// ─────────────────────────────────────────
// VER PERFIL PÚBLICO DE OUTRO USUÁRIO
// GET /api/usuarios/:id
// ─────────────────────────────────────────
func PerfilPublico(c *gin.Context) {
	id := c.Param("id")

	var usuario models.Usuario
	if err := database.DB.First(&usuario, id).Error; err != nil {
		utils.Error(c, 404, "Usuário não encontrado")
		return
	}

	// Retorna apenas dados públicos
	utils.Success(c, 200, "Perfil encontrado", gin.H{
		"id":            usuario.ID,
		"nome_completo": usuario.NomeCompleto,
		"role":          usuario.Role,
		"profissao":     usuario.Profissao,
		"bioma":         usuario.Bioma,
		"estado":        usuario.Estado,
		"municipio":     usuario.Municipio,
		"especialidade": usuario.Especialidade,
		"instituicao":   usuario.InstituicaoVinculo,
		"registro_prof": usuario.RegistroProf,
		"bio":           usuario.Bio,
		"redes_sociais": usuario.RedesSociais,
		"site_url":      usuario.SiteURL,
		"avatar_url":    usuario.AvatarURL,
		"membro_desde":  usuario.CreatedAt,
	})
}

// ─────────────────────────────────────────
// DESATIVAR MINHA CONTA
// DELETE /api/perfil
// ─────────────────────────────────────────
func DesativarConta(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

	var input struct {
		Senha string `json:"senha" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, "Informe sua senha para confirmar")
		return
	}

	var usuario models.Usuario
	if err := database.DB.First(&usuario, usuarioID).Error; err != nil {
		utils.Error(c, 404, "Usuário não encontrado")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(input.Senha)); err != nil {
		utils.Error(c, 401, "Senha incorreta")
		return
	}

	if err := database.DB.Model(&usuario).Update("ativo", false).Error; err != nil {
		utils.Error(c, 500, "Erro ao desativar conta")
		return
	}

	utils.Success(c, 200, "Conta desativada com sucesso.", nil)
}

// ─────────────────────────────────────────
// LISTAR USUÁRIOS (admin)
// GET /api/admin/usuarios
// ─────────────────────────────────────────
func ListarUsuarios(c *gin.Context) {
	var usuarios []models.Usuario

	query := database.DB.Model(&models.Usuario{})

	// Filtro por role
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}

	// Filtro por estado
	if estado := c.Query("estado"); estado != "" {
		query = query.Where("estado = ?", estado)
	}

	// Filtro por ativo/inativo
	if ativo := c.Query("ativo"); ativo != "" {
		query = query.Where("ativo = ?", ativo == "true")
	}

	// Busca por nome
	if busca := c.Query("busca"); busca != "" {
		query = query.Where("nome_completo ILIKE ?", "%"+busca+"%")
	}

	if err := query.Order("created_at DESC").Find(&usuarios).Error; err != nil {
		utils.Error(c, 500, "Erro ao buscar usuários")
		return
	}

	utils.Success(c, http.StatusOK, "Usuários encontrados", gin.H{
		"total":    len(usuarios),
		"usuarios": usuarios,
	})
}
