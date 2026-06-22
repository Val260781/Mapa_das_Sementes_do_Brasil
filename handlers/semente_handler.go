package handlers

import (
	"net/http"
	"strconv"

	"mapa-sementes-brasil/database"
	"mapa-sementes-brasil/models"

	"github.com/gin-gonic/gin"
)

// ─────────────────────────────────────────
// CRIAR SEMENTE
// POST /api/sementes
// Qualquer usuário logado
// ─────────────────────────────────────────
func CriarSemente(c *gin.Context) {
	var input struct {
		EspecieID      uint   `json:"especie_id" binding:"required"`
		Nome           string `json:"nome" binding:"required"`
		Descricao      string `json:"descricao"`
		Cor            string `json:"cor"`
		Tamanho        string `json:"tamanho"`
		Textura        string `json:"textura"`
		UsosArtesanais string `json:"usos_artesanais"`
		Tecnicas       string `json:"tecnicas_artesanato"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos: " + err.Error()})
		return
	}

	// Verifica se a espécie existe
	var especie models.Especie
	if err := database.DB.First(&especie, input.EspecieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Espécie não encontrada com o ID informado"})
		return
	}

	semente := models.Semente{
		EspecieID:      input.EspecieID,
		Nome:           input.Nome,
		Descricao:      input.Descricao,
		Cor:            input.Cor,
		Tamanho:        input.Tamanho,
		Textura:        input.Textura,
		UsosArtesanais: input.UsosArtesanais,
		Tecnicas:       input.Tecnicas,
		Validada:       false, // sempre começa pendente
	}

	if err := database.DB.Create(&semente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao cadastrar semente: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensagem": "Semente cadastrada com sucesso! Aguardando validação por especialista.",
		"semente":  semente,
	})
}

// ─────────────────────────────────────────
// LISTAR SEMENTES
// GET /api/sementes
// Público — com filtros opcionais
// ─────────────────────────────────────────
func ListarSementes(c *gin.Context) {
	var sementes []models.Semente

	query := database.DB.Preload("Especie")

	// Filtro por espécie (ex: ?especie_id=1)
	if especieID := c.Query("especie_id"); especieID != "" {
		query = query.Where("especie_id = ?", especieID)
	}

	// Filtro por validação (ex: ?validada=true)
	if validada := c.Query("validada"); validada != "" {
		query = query.Where("validada = ?", validada == "true")
	}

	// Busca por nome (ex: ?busca=açaí)
	if busca := c.Query("busca"); busca != "" {
		query = query.Where("nome ILIKE ? OR descricao ILIKE ?",
			"%"+busca+"%", "%"+busca+"%")
	}

	// Filtro por cor (ex: ?cor=marrom)
	if cor := c.Query("cor"); cor != "" {
		query = query.Where("cor ILIKE ?", "%"+cor+"%")
	}

	if err := query.Find(&sementes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar sementes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    len(sementes),
		"sementes": sementes,
	})
}

// ─────────────────────────────────────────
// DETALHE DE UMA SEMENTE
// GET /api/sementes/:id
// Público
// ─────────────────────────────────────────
func DetalheSemente(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	var semente models.Semente
	if err := database.DB.
		Preload("Especie").
		Preload("Registros").
		First(&semente, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Semente não encontrada"})
		return
	}

	c.JSON(http.StatusOK, semente)
}

// ─────────────────────────────────────────
// EDITAR SEMENTE
// PUT /api/sementes/:id
// Admin ou especialista
// ─────────────────────────────────────────
func EditarSemente(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	var semente models.Semente
	if err := database.DB.First(&semente, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Semente não encontrada"})
		return
	}

	var input struct {
		Nome           string `json:"nome"`
		Descricao      string `json:"descricao"`
		Cor            string `json:"cor"`
		Tamanho        string `json:"tamanho"`
		Textura        string `json:"textura"`
		UsosArtesanais string `json:"usos_artesanais"`
		Tecnicas       string `json:"tecnicas_artesanato"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Nome != ""           { updates["nome"] = input.Nome }
	if input.Descricao != ""      { updates["descricao"] = input.Descricao }
	if input.Cor != ""            { updates["cor"] = input.Cor }
	if input.Tamanho != ""        { updates["tamanho"] = input.Tamanho }
	if input.Textura != ""        { updates["textura"] = input.Textura }
	if input.UsosArtesanais != "" { updates["usos_artesanais"] = input.UsosArtesanais }
	if input.Tecnicas != ""       { updates["tecnicas"] = input.Tecnicas }

	if err := database.DB.Model(&semente).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao atualizar semente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Semente atualizada com sucesso!",
		"semente":  semente,
	})
}

// ─────────────────────────────────────────
// DELETAR SEMENTE
// DELETE /api/sementes/:id
// Somente admin
// ─────────────────────────────────────────
func DeletarSemente(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	if err := database.DB.Delete(&models.Semente{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao deletar semente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Semente removida com sucesso."})
}
