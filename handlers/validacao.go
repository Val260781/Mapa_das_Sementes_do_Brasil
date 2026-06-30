package handlers

import (
	"strconv"

	"mapa-sementes-brasil/database"
	"mapa-sementes-brasil/models"
	"mapa-sementes-brasil/utils"

	"github.com/gin-gonic/gin"
)

// ─────────────────────────────────────────
// VALIDAR ESPÉCIE
// POST /api/especies/:id/validar
// Somente especialista ou admin
// ─────────────────────────────────────────
func ValidarEspecie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var especie models.Especie
	if err := database.DB.First(&especie, id).Error; err != nil {
		utils.Error(c, 404, "Espécie não encontrada")
		return
	}

	if err := database.DB.Model(&especie).Update("validada", true).Error; err != nil {
		utils.Error(c, 500, "Erro ao validar espécie")
		return
	}

	utils.Success(c, 200, "Espécie validada com sucesso!", gin.H{"id": especie.ID, "validada": true})
}

// ─────────────────────────────────────────
// VALIDAR SEMENTE
// POST /api/sementes/:id/validar
// Somente especialista ou admin
// ─────────────────────────────────────────
func ValidarSemente(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var semente models.Semente
	if err := database.DB.First(&semente, id).Error; err != nil {
		utils.Error(c, 404, "Semente não encontrada")
		return
	}

	if err := database.DB.Model(&semente).Update("validada", true).Error; err != nil {
		utils.Error(c, 500, "Erro ao validar semente")
		return
	}

	utils.Success(c, 200, "Semente validada com sucesso!", gin.H{"id": semente.ID, "validada": true})
}

// ─────────────────────────────────────────
// VALIDAR REGISTRO (aprova para aparecer no mapa)
// POST /api/registros/:id/validar
// Somente especialista ou admin
// ─────────────────────────────────────────
func ValidarRegistro(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var registro models.RegistroSemente
	if err := database.DB.First(&registro, id).Error; err != nil {
		utils.Error(c, 404, "Registro não encontrado")
		return
	}

	if err := database.DB.Model(&registro).Update("status", models.StatusAprovado).Error; err != nil {
		utils.Error(c, 500, "Erro ao validar registro")
		return
	}

	utils.Success(c, 200, "Registro aprovado e disponível no mapa!", gin.H{"id": registro.ID, "status": models.StatusAprovado})
}

// ─────────────────────────────────────────
// REJEITAR REGISTRO
// POST /api/registros/:id/rejeitar
// Somente especialista ou admin
// ─────────────────────────────────────────
func RejeitarRegistro(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var registro models.RegistroSemente
	if err := database.DB.First(&registro, id).Error; err != nil {
		utils.Error(c, 404, "Registro não encontrado")
		return
	}

	if err := database.DB.Model(&registro).Update("status", models.StatusRejeitado).Error; err != nil {
		utils.Error(c, 500, "Erro ao rejeitar registro")
		return
	}

	utils.Success(c, 200, "Registro rejeitado.", gin.H{"id": registro.ID, "status": models.StatusRejeitado})
}

// ─────────────────────────────────────────
// ADMIN: EXCLUIR ESPÉCIE (definitivo)
// DELETE /api/admin/especies/:id
// Somente admin — Cuidado: falha se houver sementes vinculadas (FK).
// ─────────────────────────────────────────
func AdminExcluirEspecie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var especie models.Especie
	if err := database.DB.First(&especie, id).Error; err != nil {
		utils.Error(c, 404, "Espécie não encontrada")
		return
	}

	var totalSementes int64
	database.DB.Model(&models.Semente{}).Where("especie_id = ?", id).Count(&totalSementes)
	if totalSementes > 0 {
		utils.Error(c, 409, "Não é possível excluir: existem sementes vinculadas a esta espécie")
		return
	}

	if err := database.DB.Unscoped().Delete(&especie).Error; err != nil {
		utils.Error(c, 500, "Erro ao excluir espécie")
		return
	}

	utils.Success(c, 200, "Espécie excluída com sucesso.", nil)
}

// ─────────────────────────────────────────
// ADMIN: EXCLUIR SEMENTE (definitivo)
// DELETE /api/admin/sementes/:id
// Somente admin — Cuidado: falha se houver registros vinculados (FK).
// ─────────────────────────────────────────
func AdminExcluirSemente(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var semente models.Semente
	if err := database.DB.First(&semente, id).Error; err != nil {
		utils.Error(c, 404, "Semente não encontrada")
		return
	}

	var totalRegistros int64
	database.DB.Model(&models.RegistroSemente{}).Where("semente_id = ?", id).Count(&totalRegistros)
	if totalRegistros > 0 {
		utils.Error(c, 409, "Não é possível excluir: existem registros vinculados a esta semente")
		return
	}

	if err := database.DB.Unscoped().Delete(&semente).Error; err != nil {
		utils.Error(c, 500, "Erro ao excluir semente")
		return
	}

	utils.Success(c, 200, "Semente excluída com sucesso.", nil)
}

// ─────────────────────────────────────────
// ADMIN: EXCLUIR REGISTRO (definitivo, com fotos)
// DELETE /api/admin/registros/:id
// Somente admin — diferente da rota de usuário comum, não exige ser o dono.
// ─────────────────────────────────────────
func AdminExcluirRegistro(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var registro models.RegistroSemente
	if err := database.DB.Preload("Fotos").First(&registro, id).Error; err != nil {
		utils.Error(c, 404, "Registro não encontrado")
		return
	}

	for _, foto := range registro.Fotos {
		utils.DeletarImagem(foto.URL)
	}

	if err := database.DB.Unscoped().Delete(&registro).Error; err != nil {
		utils.Error(c, 500, "Erro ao excluir registro")
		return
	}

	utils.Success(c, 200, "Registro excluído com sucesso.", nil)
}
