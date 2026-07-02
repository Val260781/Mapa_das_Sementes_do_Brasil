package handlers

import (
	"net/http"
	"strconv"
	"time"

	"mapa-sementes-brasil/database"
	"mapa-sementes-brasil/models"

	"github.com/gin-gonic/gin"
)

// ─────────────────────────────────────────
// AVALIAR ESPÉCIE
// POST /api/avaliacoes/especie/:id
// Somente especialista ou admin
// ─────────────────────────────────────────
func AvaliarEspecie(c *gin.Context) {
	// Pega o usuário logado do contexto (inserido pelo middleware JWT)
	usuarioID := c.GetUint("usuario_id")
	role := c.GetString("role")

	if role != string(models.RoleEspecialista) && role != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"erro": "Apenas especialistas podem avaliar espécies"})
		return
	}

	especieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	// Verifica se espécie existe
	var especie models.Especie
	if err := database.DB.First(&especie, especieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Espécie não encontrada"})
		return
	}

	var input struct {
		Status               string `json:"status" binding:"required"`  // "aprovada" ou "reprovada"
		Parecer              string `json:"parecer" binding:"required"` // obrigatório sempre
		RecomendacaoCorrecao string `json:"recomendacao_correcao"`
		FontesBibliograficas string `json:"fontes_bibliograficas"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos: " + err.Error()})
		return
	}

	// Valida status
	if input.Status != string(models.AvaliacaoAprovada) && input.Status != string(models.AvaliacaoReprovada) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Status deve ser 'aprovada' ou 'reprovada'"})
		return
	}

	// Se reprovada, recomendação de correção é obrigatória
	if input.Status == string(models.AvaliacaoReprovada) && input.RecomendacaoCorrecao == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Para reprovar, informe o campo 'recomendacao_correcao' com o que deve ser corrigido",
		})
		return
	}

	especieIDUint := uint(especieID)

	avaliacao := models.Avaliacao{
		EspecialistaID:       usuarioID,
		EspecieID:            &especieIDUint,
		Status:               models.StatusAvaliacao(input.Status),
		Parecer:              input.Parecer,
		RecomendacaoCorrecao: input.RecomendacaoCorrecao,
		FontesBibliograficas: input.FontesBibliograficas,
		DataAvaliacao:        time.Now(),
	}

	if err := database.DB.Create(&avaliacao).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar avaliação"})
		return
	}

	// Atualiza o status de validação da espécie
	validada := input.Status == string(models.AvaliacaoAprovada)
	database.DB.Model(&especie).Update("validada", validada)

	c.JSON(http.StatusCreated, gin.H{
		"mensagem":         "Avaliação registrada com sucesso!",
		"avaliacao":        avaliacao,
		"especie_validada": validada,
	})
}

// ─────────────────────────────────────────
// AVALIAR SEMENTE
// POST /api/avaliacoes/semente/:id
// Somente especialista ou admin
// ─────────────────────────────────────────
func AvaliarSemente(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")
	role := c.GetString("role")

	if role != string(models.RoleEspecialista) && role != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"erro": "Apenas especialistas podem avaliar sementes"})
		return
	}

	sementeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	var semente models.Semente
	if err := database.DB.First(&semente, sementeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Semente não encontrada"})
		return
	}

	var input struct {
		Status               string `json:"status" binding:"required"`
		Parecer              string `json:"parecer" binding:"required"`
		RecomendacaoCorrecao string `json:"recomendacao_correcao"`
		FontesBibliograficas string `json:"fontes_bibliograficas"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos: " + err.Error()})
		return
	}

	if input.Status != string(models.AvaliacaoAprovada) && input.Status != string(models.AvaliacaoReprovada) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Status deve ser 'aprovada' ou 'reprovada'"})
		return
	}

	if input.Status == string(models.AvaliacaoReprovada) && input.RecomendacaoCorrecao == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Para reprovar, informe o campo 'recomendacao_correcao' com o que deve ser corrigido",
		})
		return
	}

	sementeIDUint := uint(sementeID)

	avaliacao := models.Avaliacao{
		EspecialistaID:       usuarioID,
		SementeID:            &sementeIDUint,
		Status:               models.StatusAvaliacao(input.Status),
		Parecer:              input.Parecer,
		RecomendacaoCorrecao: input.RecomendacaoCorrecao,
		FontesBibliograficas: input.FontesBibliograficas,
		DataAvaliacao:        time.Now(),
	}

	if err := database.DB.Create(&avaliacao).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar avaliação"})
		return
	}

	validada := input.Status == string(models.AvaliacaoAprovada)
	database.DB.Model(&semente).Update("validada", validada)

	c.JSON(http.StatusCreated, gin.H{
		"mensagem":         "Avaliação registrada com sucesso!",
		"avaliacao":        avaliacao,
		"semente_validada": validada,
	})
}

// ─────────────────────────────────────────
// AVALIAR REGISTRO
// POST /api/avaliacoes/registro/:id
// Somente especialista ou admin
// ─────────────────────────────────────────
func AvaliarRegistro(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")
	role := c.GetString("role")

	if role != string(models.RoleEspecialista) && role != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"erro": "Apenas especialistas podem avaliar registros"})
		return
	}

	registroID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	var registro models.RegistroSemente
	if err := database.DB.First(&registro, registroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Registro não encontrado"})
		return
	}

	var input struct {
		Status               string `json:"status" binding:"required"`  // "aprovada" ou "reprovada"
		Parecer              string `json:"parecer" binding:"required"` // obrigatório sempre
		RecomendacaoCorrecao string `json:"recomendacao_correcao"`
		FontesBibliograficas string `json:"fontes_bibliograficas"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos: " + err.Error()})
		return
	}

	if input.Status != string(models.AvaliacaoAprovada) && input.Status != string(models.AvaliacaoReprovada) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Status deve ser 'aprovada' ou 'reprovada'"})
		return
	}

	if input.Status == string(models.AvaliacaoReprovada) && input.RecomendacaoCorrecao == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Para reprovar, informe o campo 'recomendacao_correcao' com o que deve ser corrigido",
		})
		return
	}

	registroIDUint := uint(registroID)

	avaliacao := models.Avaliacao{
		EspecialistaID:       usuarioID,
		RegistroID:           &registroIDUint,
		Status:               models.StatusAvaliacao(input.Status),
		Parecer:              input.Parecer,
		RecomendacaoCorrecao: input.RecomendacaoCorrecao,
		FontesBibliograficas: input.FontesBibliograficas,
		DataAvaliacao:        time.Now(),
	}

	if err := database.DB.Create(&avaliacao).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar avaliação"})
		return
	}

	// Registro tem 3 estados (pendente/aprovado/rejeitado), não um bool
	novoStatus := models.StatusRejeitado
	if input.Status == string(models.AvaliacaoAprovada) {
		novoStatus = models.StatusAprovado
	}
	database.DB.Model(&registro).Update("status", novoStatus)

	c.JSON(http.StatusCreated, gin.H{
		"mensagem":        "Avaliação registrada com sucesso!",
		"avaliacao":       avaliacao,
		"registro_status": novoStatus,
	})
}

// ─────────────────────────────────────────
// AVALIAR CONHECIMENTO TRADICIONAL
// POST /api/avaliacoes/conhecimento/:id
// Somente especialista ou admin
// ─────────────────────────────────────────
func AvaliarConhecimento(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")
	role := c.GetString("role")

	if role != string(models.RoleEspecialista) && role != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"erro": "Apenas especialistas podem avaliar conhecimentos tradicionais"})
		return
	}

	conhecimentoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	var conhecimento models.ConhecimentoTradicional
	if err := database.DB.First(&conhecimento, conhecimentoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Conhecimento não encontrado"})
		return
	}

	var input struct {
		Status               string `json:"status" binding:"required"`
		Parecer              string `json:"parecer" binding:"required"`
		RecomendacaoCorrecao string `json:"recomendacao_correcao"`
		FontesBibliograficas string `json:"fontes_bibliograficas"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos: " + err.Error()})
		return
	}

	if input.Status != string(models.AvaliacaoAprovada) && input.Status != string(models.AvaliacaoReprovada) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Status deve ser 'aprovada' ou 'reprovada'"})
		return
	}

	if input.Status == string(models.AvaliacaoReprovada) && input.RecomendacaoCorrecao == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Para reprovar, informe o campo 'recomendacao_correcao' com o que deve ser corrigido",
		})
		return
	}

	conhecimentoIDUint := uint(conhecimentoID)

	avaliacao := models.Avaliacao{
		EspecialistaID:       usuarioID,
		ConhecimentoID:       &conhecimentoIDUint,
		Status:               models.StatusAvaliacao(input.Status),
		Parecer:              input.Parecer,
		RecomendacaoCorrecao: input.RecomendacaoCorrecao,
		FontesBibliograficas: input.FontesBibliograficas,
		DataAvaliacao:        time.Now(),
	}

	if err := database.DB.Create(&avaliacao).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar avaliação"})
		return
	}

	validado := input.Status == string(models.AvaliacaoAprovada)
	database.DB.Model(&conhecimento).Update("validado", validado)

	c.JSON(http.StatusCreated, gin.H{
		"mensagem":              "Avaliação registrada com sucesso!",
		"avaliacao":             avaliacao,
		"conhecimento_validado": validado,
	})
}

// ─────────────────────────────────────────
// LISTAR PENDENTES
// GET /api/avaliacoes/pendentes
// Somente especialista ou admin
// ─────────────────────────────────────────
func ListarPendentes(c *gin.Context) {
	role := c.GetString("role")
	if role != string(models.RoleEspecialista) && role != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"erro": "Acesso restrito a especialistas"})
		return
	}

	var especiesPendentes []models.Especie
	var sementesPendentes []models.Semente
	var registrosPendentes []models.RegistroSemente
	var conhecimentosPendentes []models.ConhecimentoTradicional

	database.DB.Where("validada = false").Find(&especiesPendentes)
	database.DB.Where("validada = false").Preload("Especie").Find(&sementesPendentes)
	database.DB.Where("status = ?", models.StatusPendente).Preload("Semente").Preload("Usuario").Preload("Fotos").Find(&registrosPendentes)
	database.DB.Where("validado = false").Preload("Semente").Preload("Usuario").Find(&conhecimentosPendentes)

	c.JSON(http.StatusOK, gin.H{
		"especies_pendentes":      especiesPendentes,
		"sementes_pendentes":      sementesPendentes,
		"registros_pendentes":     registrosPendentes,
		"conhecimentos_pendentes": conhecimentosPendentes,
		"total_especies":          len(especiesPendentes),
		"total_sementes":          len(sementesPendentes),
		"total_registros":         len(registrosPendentes),
		"total_conhecimentos":     len(conhecimentosPendentes),
	})
}

// ─────────────────────────────────────────
// HISTÓRICO DE AVALIAÇÕES
// GET /api/avaliacoes/historico
// Somente especialista ou admin
// ─────────────────────────────────────────
func HistoricoAvaliacoes(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")
	role := c.GetString("role")

	var avaliacoes []models.Avaliacao
	query := database.DB.Preload("Especie").Preload("Semente").Preload("Registro").Preload("Conhecimento").Preload("Especialista")

	// Admin vê tudo, especialista vê só as suas
	if role == string(models.RoleEspecialista) {
		query = query.Where("especialista_id = ?", usuarioID)
	}

	if err := query.Order("created_at DESC").Find(&avaliacoes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar histórico"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":      len(avaliacoes),
		"avaliacoes": avaliacoes,
	})
}
