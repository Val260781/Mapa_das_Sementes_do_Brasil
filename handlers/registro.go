package handlers

import (
	"net/http"
	"strconv"
	"time"

	"mapa-sementes-brasil/database"
	"mapa-sementes-brasil/models"
	"mapa-sementes-brasil/utils"

	"github.com/gin-gonic/gin"
)

// ─────────────────────────────────────────
// CRIAR REGISTRO DE SEMENTE
// POST /api/registros
// Qualquer usuário logado
// ─────────────────────────────────────────
func CriarRegistro(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

	var input struct {
		SementeID  uint    `json:"semente_id" binding:"required"`
		Latitude   float64 `json:"latitude" binding:"required"`
		Longitude  float64 `json:"longitude" binding:"required"`
		Estado     string  `json:"estado"`
		Municipio  string  `json:"municipio"`
		Descricao  string  `json:"descricao"`
		Quantidade int     `json:"quantidade"`
		DataColeta string  `json:"data_coleta"` // AAAA-MM-DD
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, "Dados inválidos: "+err.Error())
		return
	}

	// Valida coordenadas
	if input.Latitude < -90 || input.Latitude > 90 {
		utils.Error(c, 400, "Latitude inválida (deve estar entre -90 e 90)")
		return
	}
	if input.Longitude < -180 || input.Longitude > 180 {
		utils.Error(c, 400, "Longitude inválida (deve estar entre -180 e 180)")
		return
	}

	// Verifica se a semente existe
	var semente models.Semente
	if err := database.DB.First(&semente, input.SementeID).Error; err != nil {
		utils.Error(c, 404, "Semente não encontrada com o ID informado")
		return
	}

	// Converte data de coleta
	dataColeta := time.Now()
	if input.DataColeta != "" {
		var err error
		dataColeta, err = time.Parse("2006-01-02", input.DataColeta)
		if err != nil {
			utils.Error(c, 400, "Formato de data inválido. Use: AAAA-MM-DD (ex: 2024-03-15)")
			return
		}
	}

	registro := models.RegistroSemente{
		UsuarioID:  usuarioID,
		SementeID:  input.SementeID,
		Latitude:   input.Latitude,
		Longitude:  input.Longitude,
		Estado:     input.Estado,
		Municipio:  input.Municipio,
		Descricao:  input.Descricao,
		Quantidade: input.Quantidade,
		DataColeta: dataColeta,
		Status:     models.StatusPendente,
	}

	if err := database.DB.Create(&registro).Error; err != nil {
		utils.Error(c, 500, "Erro ao salvar registro: "+err.Error())
		return
	}

	// Recarrega com relacionamentos
	database.DB.Preload("Semente").Preload("Semente.Especie").Preload("Usuario").First(&registro, registro.ID)

	utils.Success(c, http.StatusCreated, "Registro criado com sucesso!", registro)
}

// ─────────────────────────────────────────
// LISTAR REGISTROS
// GET /api/registros
// Público — com filtros opcionais
// ─────────────────────────────────────────
func ListarRegistros(c *gin.Context) {
	var registros []models.RegistroSemente

	query := database.DB.
		Preload("Semente").
		Preload("Semente.Especie").
		Preload("Usuario").
		Preload("Fotos")

	// Filtro por semente
	if sementeID := c.Query("semente_id"); sementeID != "" {
		query = query.Where("semente_id = ?", sementeID)
	}

	// Filtro por estado
	if estado := c.Query("estado"); estado != "" {
		query = query.Where("estado ILIKE ?", estado)
	}

	// Filtro por município
	if municipio := c.Query("municipio"); municipio != "" {
		query = query.Where("municipio ILIKE ?", "%"+municipio+"%")
	}

	// Filtro por status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filtro por usuário
	if usuarioID := c.Query("usuario_id"); usuarioID != "" {
		query = query.Where("usuario_id = ?", usuarioID)
	}

	// Filtro por área geográfica (bounding box)
	// Ex: ?lat_min=-16&lat_max=-15&lng_min=-48&lng_max=-47
	if latMin := c.Query("lat_min"); latMin != "" {
		query = query.Where("latitude >= ?", latMin)
	}
	if latMax := c.Query("lat_max"); latMax != "" {
		query = query.Where("latitude <= ?", latMax)
	}
	if lngMin := c.Query("lng_min"); lngMin != "" {
		query = query.Where("longitude >= ?", lngMin)
	}
	if lngMax := c.Query("lng_max"); lngMax != "" {
		query = query.Where("longitude <= ?", lngMax)
	}

	if err := query.Order("created_at DESC").Find(&registros).Error; err != nil {
		utils.Error(c, 500, "Erro ao buscar registros")
		return
	}

	utils.Success(c, http.StatusOK, "Registros encontrados", gin.H{
		"total":     len(registros),
		"registros": registros,
	})
}

// ─────────────────────────────────────────
// DETALHE DE UM REGISTRO
// GET /api/registros/:id
// Público
// ─────────────────────────────────────────
func DetalheRegistro(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID inválido")
		return
	}

	var registro models.RegistroSemente
	if err := database.DB.
		Preload("Semente").
		Preload("Semente.Especie").
		Preload("Usuario").
		Preload("Fotos").
		First(&registro, id).Error; err != nil {
		utils.Error(c, 404, "Registro não encontrado")
		return
	}

	utils.Success(c, http.StatusOK, "Registro encontrado", registro)
}

// ─────────────────────────────────────────
// EDITAR REGISTRO
// PUT /api/registros/:id
// Somente o dono do registro
// ─────────────────────────────────────────
func EditarRegistro(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

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

	// Verifica se é o dono
	if registro.UsuarioID != usuarioID {
		utils.Error(c, 403, "Você não tem permissão para editar este registro")
		return
	}

	var input struct {
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
		Estado     string  `json:"estado"`
		Municipio  string  `json:"municipio"`
		Descricao  string  `json:"descricao"`
		Quantidade int     `json:"quantidade"`
		DataColeta string  `json:"data_coleta"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if input.Latitude != 0 {
		updates["latitude"] = input.Latitude
	}
	if input.Longitude != 0 {
		updates["longitude"] = input.Longitude
	}
	if input.Estado != "" {
		updates["estado"] = input.Estado
	}
	if input.Municipio != "" {
		updates["municipio"] = input.Municipio
	}
	if input.Descricao != "" {
		updates["descricao"] = input.Descricao
	}
	if input.Quantidade > 0 {
		updates["quantidade"] = input.Quantidade
	}

	if input.DataColeta != "" {
		dataColeta, err := time.Parse("2006-01-02", input.DataColeta)
		if err != nil {
			utils.Error(c, 400, "Formato de data inválido. Use: AAAA-MM-DD")
			return
		}
		updates["data_coleta"] = dataColeta
	}

	if len(updates) == 0 {
		utils.Error(c, 400, "Nenhum campo para atualizar foi informado")
		return
	}

	if err := database.DB.Model(&registro).Updates(updates).Error; err != nil {
		utils.Error(c, 500, "Erro ao atualizar registro")
		return
	}

	database.DB.Preload("Semente").Preload("Fotos").First(&registro, id)
	utils.Success(c, http.StatusOK, "Registro atualizado com sucesso!", registro)
}

// ─────────────────────────────────────────
// DELETAR REGISTRO
// DELETE /api/registros/:id
// Dono ou admin
// ─────────────────────────────────────────
func DeletarRegistro(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")
	role := c.GetString("role")

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

	// Verifica permissão
	if registro.UsuarioID != usuarioID && role != string(models.RoleAdmin) {
		utils.Error(c, 403, "Você não tem permissão para deletar este registro")
		return
	}

	// Remove fotos do disco
	for _, foto := range registro.Fotos {
		utils.DeletarImagem(foto.URL)
	}

	if err := database.DB.Delete(&registro).Error; err != nil {
		utils.Error(c, 500, "Erro ao deletar registro")
		return
	}

	utils.Success(c, http.StatusOK, "Registro removido com sucesso.", nil)
}

// ─────────────────────────────────────────
// UPLOAD DE FOTOS DO REGISTRO
// POST /api/registros/:id/fotos
// Dono do registro
// ─────────────────────────────────────────
func UploadFotosRegistro(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

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

	if registro.UsuarioID != usuarioID {
		utils.Error(c, 403, "Você não tem permissão para adicionar fotos neste registro")
		return
	}

	// Aceita múltiplos arquivos no campo "fotos"
	form, err := c.MultipartForm()
	if err != nil {
		utils.Error(c, 400, "Erro ao processar formulário de upload")
		return
	}

	arquivos := form.File["fotos"]
	if len(arquivos) == 0 {
		utils.Error(c, 400, "Nenhuma foto enviada. Use o campo 'fotos'")
		return
	}

	if len(arquivos) > 5 {
		utils.Error(c, 400, "Máximo de 5 fotos por registro")
		return
	}

	var fotasSalvas []models.FotoRegistro

	for _, arquivo := range arquivos {
		urlFoto, err := utils.SalvarImagem(c, arquivo, "registros")
		if err != nil {
			utils.Error(c, 400, "Erro ao salvar foto '"+arquivo.Filename+"': "+err.Error())
			return
		}

		// Legenda opcional via query param
		foto := models.FotoRegistro{
			RegistroID: uint(id),
			URL:        urlFoto,
			Legenda:    c.PostForm("legenda"),
		}

		if err := database.DB.Create(&foto).Error; err != nil {
			utils.Error(c, 500, "Erro ao salvar foto no banco")
			return
		}

		fotasSalvas = append(fotasSalvas, foto)
	}

	utils.Success(c, http.StatusCreated, "Fotos enviadas com sucesso!", gin.H{
		"fotos_salvas": fotasSalvas,
		"total":        len(fotasSalvas),
	})
}

// ─────────────────────────────────────────
// DELETAR FOTO DO REGISTRO
// DELETE /api/registros/:id/fotos/:foto_id
// Dono do registro
// ─────────────────────────────────────────
func DeletarFotoRegistro(c *gin.Context) {
	usuarioID := c.GetUint("usuario_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, 400, "ID de registro inválido")
		return
	}

	fotoID, err := strconv.Atoi(c.Param("foto_id"))
	if err != nil {
		utils.Error(c, 400, "ID de foto inválido")
		return
	}

	// Verifica se o registro pertence ao usuário
	var registro models.RegistroSemente
	if err := database.DB.First(&registro, id).Error; err != nil {
		utils.Error(c, 404, "Registro não encontrado")
		return
	}

	if registro.UsuarioID != usuarioID {
		utils.Error(c, 403, "Você não tem permissão para remover esta foto")
		return
	}

	var foto models.FotoRegistro
	if err := database.DB.Where("id = ? AND registro_id = ?", fotoID, id).First(&foto).Error; err != nil {
		utils.Error(c, 404, "Foto não encontrada")
		return
	}

	utils.DeletarImagem(foto.URL)

	if err := database.DB.Delete(&foto).Error; err != nil {
		utils.Error(c, 500, "Erro ao deletar foto")
		return
	}

	utils.Success(c, http.StatusOK, "Foto removida com sucesso.", nil)
}

// ─────────────────────────────────────────
// REGISTROS NO MAPA (GeoJSON)
// GET /api/registros/mapa
// Público — retorna formato otimizado para mapa
// ─────────────────────────────────────────
func RegistrosParaMapa(c *gin.Context) {
	var registros []models.RegistroSemente

	query := database.DB.
		Preload("Semente").
		Preload("Semente.Especie").
		Preload("Fotos").
		Where("status = ?", models.StatusAprovado)

	if estado := c.Query("estado"); estado != "" {
		query = query.Where("estado ILIKE ?", estado)
	}

	if sementeID := c.Query("semente_id"); sementeID != "" {
		query = query.Where("semente_id = ?", sementeID)
	}

	if err := query.Find(&registros).Error; err != nil {
		utils.Error(c, 500, "Erro ao buscar registros para o mapa")
		return
	}

	// Formata como GeoJSON para uso em mapas (Leaflet, Google Maps, etc.)
	type Ponto struct {
		ID        uint    `json:"id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Semente   string  `json:"semente"`
		Especie   string  `json:"especie"`
		Municipio string  `json:"municipio"`
		Estado    string  `json:"estado"`
		FotoURL   string  `json:"foto_url,omitempty"`
	}

	var pontos []Ponto
	for _, r := range registros {
		ponto := Ponto{
			ID:        r.ID,
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
			Municipio: r.Municipio,
			Estado:    r.Estado,
			Semente:   r.Semente.Nome,
			Especie:   r.Semente.Especie.NomePopular,
		}
		if len(r.Fotos) > 0 {
			ponto.FotoURL = r.Fotos[0].URL
		}
		pontos = append(pontos, ponto)
	}

	utils.Success(c, http.StatusOK, "Pontos do mapa carregados", gin.H{
		"total":  len(pontos),
		"pontos": pontos,
	})
}
