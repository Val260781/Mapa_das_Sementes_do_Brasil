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
// CRIAR ESPÉCIE
// POST /api/especies
// Qualquer usuário logado pode cadastrar
// ─────────────────────────────────────────
func CriarEspecie(c *gin.Context) {
	var input struct {
		NomeCientifico    string `json:"nome_cientifico" binding:"required"`
		NomePopular       string `json:"nome_popular"`
		Familia           string `json:"familia"`
		Bioma             string `json:"bioma"`
		Descricao         string `json:"descricao"`
		UsosArtesanais    string `json:"usos_artesanais"`
		EpocaColheita     string `json:"epoca_colheita"`
		StatusConservacao string `json:"status_conservacao"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos: " + err.Error()})
		return
	}

	especie := models.Especie{
		NomeCientifico:    input.NomeCientifico,
		NomePopular:       input.NomePopular,
		Familia:           input.Familia,
		Bioma:             input.Bioma,
		Descricao:         input.Descricao,
		UsosArtesanais:    input.UsosArtesanais,
		Epoca:             input.EpocaColheita,
		StatusConservacao: input.StatusConservacao,
		Validada:          false, // sempre começa pendente
	}

	if err := database.DB.Create(&especie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao cadastrar espécie: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensagem": "Espécie cadastrada com sucesso! Aguardando validação por especialista.",
		"especie":  especie,
	})
}

// ─────────────────────────────────────────
// LISTAR ESPÉCIES
// GET /api/especies
// Público — com filtros opcionais por bioma e status
// ─────────────────────────────────────────
func ListarEspecies(c *gin.Context) {
	var especies []models.Especie

	query := database.DB.Preload("Sementes")

	// Filtro por bioma (ex: ?bioma=Cerrado)
	if bioma := c.Query("bioma"); bioma != "" {
		query = query.Where("bioma ILIKE ?", "%"+bioma+"%")
	}

	// Filtro por validação (ex: ?validada=true)
	if validada := c.Query("validada"); validada != "" {
		query = query.Where("validada = ?", validada == "true")
	}

	// Busca por nome (ex: ?busca=ipê)
	if busca := c.Query("busca"); busca != "" {
		query = query.Where("nome_cientifico ILIKE ? OR nome_popular ILIKE ?",
			"%"+busca+"%", "%"+busca+"%")
	}

	if err := query.Find(&especies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar espécies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":   len(especies),
		"especies": especies,
	})
}

// ─────────────────────────────────────────
// DETALHE DE UMA ESPÉCIE
// GET /api/especies/:id
// Público
// ─────────────────────────────────────────
func DetalheEspecie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	var especie models.Especie
	if err := database.DB.Preload("Sementes").First(&especie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Espécie não encontrada"})
		return
	}

	c.JSON(http.StatusOK, especie)
}

// ─────────────────────────────────────────
// EDITAR ESPÉCIE
// PUT /api/especies/:id
// Somente admin ou especialista
// ─────────────────────────────────────────
func EditarEspecie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	var especie models.Especie
	if err := database.DB.First(&especie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Espécie não encontrada"})
		return
	}

	var input struct {
		NomeCientifico    string `json:"nome_cientifico"`
		NomePopular       string `json:"nome_popular"`
		Familia           string `json:"familia"`
		Bioma             string `json:"bioma"`
		Descricao         string `json:"descricao"`
		UsosArtesanais    string `json:"usos_artesanais"`
		EpocaColheita     string `json:"epoca_colheita"`
		StatusConservacao string `json:"status_conservacao"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	// Atualiza apenas campos preenchidos
	updates := map[string]interface{}{}
	if input.NomeCientifico != "" { updates["nome_cientifico"] = input.NomeCientifico }
	if input.NomePopular != ""    { updates["nome_popular"] = input.NomePopular }
	if input.Familia != ""        { updates["familia"] = input.Familia }
	if input.Bioma != ""          { updates["bioma"] = input.Bioma }
	if input.Descricao != ""      { updates["descricao"] = input.Descricao }
	if input.UsosArtesanais != "" { updates["usos_artesanais"] = input.UsosArtesanais }
	if input.EpocaColheita != ""  { updates["epoca"] = input.EpocaColheita }
	if input.StatusConservacao != "" { updates["status_conservacao"] = input.StatusConservacao }

	if err := database.DB.Model(&especie).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao atualizar espécie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Espécie atualizada com sucesso!",
		"especie":  especie,
	})
}

// ─────────────────────────────────────────
// DELETAR ESPÉCIE
// DELETE /api/especies/:id
// Somente admin
// ─────────────────────────────────────────
func DeletarEspecie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	if err := database.DB.Delete(&models.Especie{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao deletar espécie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Espécie removida com sucesso."})
}

// ─────────────────────────────────────────
// UPLOAD DE FOTO DA ESPÉCIE
// POST /api/especies/:id/foto
// Qualquer usuário logado
// ─────────────────────────────────────────
func UploadFotoEspecie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	var especie models.Especie
	if err := database.DB.First(&especie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Espécie não encontrada"})
		return
	}

	arquivo, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Nenhuma foto enviada. Use o campo 'foto'"})
		return
	}

	urlFoto, err := utils.SalvarImagem(c, arquivo, "especies")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	// Remove foto anterior do disco se existir
	if especie.ImagemURL != "" {
		utils.DeletarImagem(especie.ImagemURL)
	}

	// Atualiza URL no banco
	if err := database.DB.Model(&especie).Update("imagem_url", urlFoto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar URL da foto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem":   "Foto enviada com sucesso!",
		"imagem_url": urlFoto,
	})
}
