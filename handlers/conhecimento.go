package handlers

import (
	"net/http"
	"strconv"

	"mapa-sementes-brasil/database"
	"mapa-sementes-brasil/models"
	"mapa-sementes-brasil/utils"

	"github.com/gin-gonic/gin"
)

// ─────────────────────────────────────────
// CADASTRAR CONHECIMENTO TRADICIONAL
// POST /api/conhecimentos
// Qualquer usuário logado
// ─────────────────────────────────────────
func CriarConhecimento(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

	var input struct {
		SementeID uint   `json:"semente_id" binding:"required"`
		Titulo    string `json:"titulo" binding:"required"`
		Conteudo  string `json:"conteudo" binding:"required"`
		Origem    string `json:"origem"`  // comunidade, povo indígena, região
		Tecnica   string `json:"tecnica"` // técnica artesanal descrita
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, "Dados inválidos: "+err.Error())
		return
	}

	// Verifica se a semente existe
	var semente models.Semente
	if err := database.DB.First(&semente, input.SementeID).Error; err != nil {
		utils.Error(c, 404, "Semente não encontrada com o ID informado")
		return
	}

	conhecimento := models.ConhecimentoTradicional{
		UsuarioID: usuarioID,
		SementeID: input.SementeID,
		Titulo:    input.Titulo,
		Conteudo:  input.Conteudo,
		Origem:    input.Origem,
		Tecnica:   input.Tecnica,
		Validado:  false,
		Curtidas:  0,
	}

	if err := database.DB.Create(&conhecimento).Error; err != nil {
		utils.Error(c, 500, "Erro ao salvar conhecimento: "+err.Error())
		return
	}

	// Recarrega com relacionamentos
	database.DB.
		Preload("Usuario").
		Preload("Semente").
		Preload("Semente.Especie").
		First(&conhecimento, conhecimento.ID)

	utils.Success(c, http.StatusCreated, "Conhecimento tradicional cadastrado com sucesso! Aguardando validação.", conhecimento)
}

// ─────────────────────────────────────────
// LISTAR CONHECIMENTOS
// GET /api/conhecimentos
// Público — com filtros opcionais
// ─────────────────────────────────────────
func ListarConhecimentos(c *gin.Context) {
	var conhecimentos []models.ConhecimentoTradicional

	query := database.DB.
		Preload("Usuario").
		Preload("Semente").
		Preload("Semente.Especie")

	// Filtro por semente
	if sementeID := c.Query("semente_id"); sementeID != "" {
		query = query.Where("semente_id = ?", sementeID)
	}

	// Filtro por validação
	if validado := c.Query("validado"); validado != "" {
		query = query.Where("validado = ?", validado == "true")
	}

	// Filtro por origem (comunidade, povo, região)
	if origem := c.Query("origem"); origem != "" {
		query = query.Where("origem ILIKE ?", "%"+origem+"%")
	}

	// Filtro por técnica
	if tecnica := c.Query("tecnica"); tecnica != "" {
		query = query.Where("tecnica ILIKE ?", "%"+tecnica+"%")
	}

	// Busca por título ou conteúdo
	if busca := c.Query("busca"); busca != "" {
		query = query.Where("titulo ILIKE ? OR conteudo ILIKE ?",
			"%"+busca+"%", "%"+busca+"%")
	}

	// Filtro por usuário
	if usuarioID := c.Query("usuario_id"); usuarioID != "" {
		query = query.Where("usuario_id = ?", usuarioID)
	}

	// Ordenação por curtidas ou data
	ordem := c.Query("ordem")
	if ordem == "curtidas" {
		query = query.Order("curtidas DESC")
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Find(&conhecimentos).Error; err != nil {
		utils.Error(c, 500, "Erro ao buscar conhecimentos")
		return
	}

	utils.Success(c, http.StatusOK, "Conhecimentos encontrados", gin.H{
		"total":         len(conhecimentos),
		"conhecimentos": conhecimentos,
	})
}

// ─────────────────────────────────────────
// DETALHE DE UM CONHECIMENTO
// GET /api/conhecimentos/:id
// Público
// ─────────────────────────────────────────
func DetalheConhecimento(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var conhecimento models.ConhecimentoTradicional
	if err := database.DB.
		Preload("Usuario").
		Preload("Semente").
		Preload("Semente.Especie").
		First(&conhecimento, id).Error; err != nil {
		utils.Error(c, 404, "Conhecimento não encontrado")
		return
	}

	utils.Success(c, http.StatusOK, "Conhecimento encontrado", conhecimento)
}

// ─────────────────────────────────────────
// EDITAR CONHECIMENTO
// PUT /api/conhecimentos/:id
// Somente o dono
// ─────────────────────────────────────────
func EditarConhecimento(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var conhecimento models.ConhecimentoTradicional
	if err := database.DB.First(&conhecimento, id).Error; err != nil {
		utils.Error(c, 404, "Conhecimento não encontrado")
		return
	}

	// Verifica se é o dono
	if conhecimento.UsuarioID != usuarioID {
		utils.Error(c, 403, "Você não tem permissão para editar este conhecimento")
		return
	}

	var input struct {
		Titulo   string `json:"titulo"`
		Conteudo string `json:"conteudo"`
		Origem   string `json:"origem"`
		Tecnica  string `json:"tecnica"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if input.Titulo != "" {
		updates["titulo"] = input.Titulo
	}
	if input.Conteudo != "" {
		updates["conteudo"] = input.Conteudo
	}
	if input.Origem != "" {
		updates["origem"] = input.Origem
	}
	if input.Tecnica != "" {
		updates["tecnica"] = input.Tecnica
	}

	// Edição reseta a validação — precisa ser revalidado pelo especialista
	if len(updates) > 0 {
		updates["validado"] = false
	}

	if len(updates) == 0 {
		utils.Error(c, 400, "Nenhum campo para atualizar foi informado")
		return
	}

	if err := database.DB.Model(&conhecimento).Updates(updates).Error; err != nil {
		utils.Error(c, 500, "Erro ao atualizar conhecimento")
		return
	}

	database.DB.Preload("Usuario").Preload("Semente").First(&conhecimento, id)
	utils.Success(c, http.StatusOK, "Conhecimento atualizado! Aguarda nova validação.", conhecimento)
}

// ─────────────────────────────────────────
// DELETAR CONHECIMENTO
// DELETE /api/conhecimentos/:id
// Dono ou admin
// ─────────────────────────────────────────
func DeletarConhecimento(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")
	role := c.GetString("role")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var conhecimento models.ConhecimentoTradicional
	if err := database.DB.First(&conhecimento, id).Error; err != nil {
		utils.Error(c, 404, "Conhecimento não encontrado")
		return
	}

	if conhecimento.UsuarioID != usuarioID && role != string(models.RoleAdmin) {
		utils.Error(c, 403, "Você não tem permissão para deletar este conhecimento")
		return
	}

	if err := database.DB.Delete(&conhecimento).Error; err != nil {
		utils.Error(c, 500, "Erro ao deletar conhecimento")
		return
	}

	utils.Success(c, http.StatusOK, "Conhecimento removido com sucesso.", nil)
}

// ─────────────────────────────────────────
// CURTIR CONHECIMENTO
// POST /api/conhecimentos/:id/curtir
// Qualquer usuário logado
// ─────────────────────────────────────────
func CurtirConhecimento(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var conhecimento models.ConhecimentoTradicional
	if err := database.DB.First(&conhecimento, id).Error; err != nil {
		utils.Error(c, 404, "Conhecimento não encontrado")
		return
	}

	// Incrementa curtidas
	if err := database.DB.Model(&conhecimento).
		Update("curtidas", conhecimento.Curtidas+1).Error; err != nil {
		utils.Error(c, 500, "Erro ao registrar curtida")
		return
	}

	utils.Success(c, http.StatusOK, "Curtida registrada!", gin.H{
		"curtidas": conhecimento.Curtidas + 1,
	})
}

// ─────────────────────────────────────────
// VALIDAR CONHECIMENTO (especialista)
// POST /api/conhecimentos/:id/validar
// Somente especialista ou admin
// ─────────────────────────────────────────
func ValidarConhecimento(c *gin.Context) {
	role := c.GetString("role")

	if role != string(models.RoleEspecialista) && role != string(models.RoleAdmin) {
		utils.Error(c, 403, "Apenas especialistas podem validar conhecimentos tradicionais")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var conhecimento models.ConhecimentoTradicional
	if err := database.DB.First(&conhecimento, id).Error; err != nil {
		utils.Error(c, 404, "Conhecimento não encontrado")
		return
	}

	var input struct {
		Validado bool   `json:"validado"`                   // true = valida, false = reprova
		Parecer  string `json:"parecer" binding:"required"` // obrigatório sempre
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, "Informe o campo 'parecer' com a justificativa")
		return
	}

	if err := database.DB.Model(&conhecimento).
		Update("validado", input.Validado).Error; err != nil {
		utils.Error(c, 500, "Erro ao validar conhecimento")
		return
	}

	status := "reprovado"
	if input.Validado {
		status = "validado"
	}

	utils.Success(c, http.StatusOK, "Conhecimento "+status+" com sucesso!", gin.H{
		"id":       conhecimento.ID,
		"validado": input.Validado,
		"parecer":  input.Parecer,
	})
}

// ─────────────────────────────────────────
// CONHECIMENTOS POR SEMENTE
// GET /api/sementes/:id/conhecimentos
// Público
// ─────────────────────────────────────────
func ConhecimentosPorSemente(c *gin.Context) {
	sementeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	// Verifica se a semente existe
	var semente models.Semente
	if err := database.DB.First(&semente, sementeID).Error; err != nil {
		utils.Error(c, 404, "Semente não encontrada")
		return
	}

	var conhecimentos []models.ConhecimentoTradicional
	if err := database.DB.
		Where("semente_id = ? AND validado = true", sementeID).
		Preload("Usuario").
		Order("curtidas DESC").
		Find(&conhecimentos).Error; err != nil {
		utils.Error(c, 500, "Erro ao buscar conhecimentos")
		return
	}

	utils.Success(c, http.StatusOK, "Conhecimentos da semente encontrados", gin.H{
		"semente":       semente.Nome,
		"total":         len(conhecimentos),
		"conhecimentos": conhecimentos,
	})
}
