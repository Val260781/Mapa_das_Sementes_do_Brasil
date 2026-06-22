package handlers

import (
	"math"
	"net/http"
	"strconv"

	"mapa-sementes-brasil/database"
	"mapa-sementes-brasil/models"
	"mapa-sementes-brasil/utils"

	"github.com/gin-gonic/gin"
)

// ─────────────────────────────────────────
// BUSCA GERAL
// GET /api/busca?q=olho-de-cabra
// Busca em espécies, sementes e conhecimentos ao mesmo tempo
// Público
// ─────────────────────────────────────────
func BuscaGeral(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		utils.Error(c, 400, "Informe o parâmetro 'q' para buscar. Ex: /api/busca?q=olho-de-cabra")
		return
	}

	termo := "%" + q + "%"

	// Busca em espécies
	var especies []models.Especie
	database.DB.
		Where("nome_cientifico ILIKE ? OR nome_popular ILIKE ? OR descricao ILIKE ? OR bioma ILIKE ?",
			termo, termo, termo, termo).
		Where("validada = true").
		Limit(10).
		Find(&especies)

	// Busca em sementes
	var sementes []models.Semente
	database.DB.
		Preload("Especie").
		Where("nome ILIKE ? OR descricao ILIKE ? OR usos_artesanais ILIKE ? OR cor ILIKE ?",
			termo, termo, termo, termo).
		Where("validada = true").
		Limit(10).
		Find(&sementes)

	// Busca em conhecimentos tradicionais
	var conhecimentos []models.ConhecimentoTradicional
	database.DB.
		Preload("Semente").
		Preload("Usuario").
		Where("titulo ILIKE ? OR conteudo ILIKE ? OR origem ILIKE ? OR tecnica ILIKE ?",
			termo, termo, termo, termo).
		Where("validado = true").
		Limit(10).
		Find(&conhecimentos)

	// Busca em registros (município e estado)
	var registros []models.RegistroSemente
	database.DB.
		Preload("Semente").
		Preload("Semente.Especie").
		Where("municipio ILIKE ? OR estado ILIKE ? OR descricao ILIKE ?",
			termo, termo, termo).
		Where("status = ?", models.StatusAprovado).
		Limit(10).
		Find(&registros)

	total := len(especies) + len(sementes) + len(conhecimentos) + len(registros)

	utils.Success(c, http.StatusOK, "Busca concluída", gin.H{
		"termo":         q,
		"total":         total,
		"especies":      especies,
		"sementes":      sementes,
		"conhecimentos": conhecimentos,
		"registros":     registros,
	})
}

// ─────────────────────────────────────────
// BUSCA AVANÇADA DE ESPÉCIES
// GET /api/busca/especies
// Filtros: q, bioma, status_conservacao, validado, familia
// Público
// ─────────────────────────────────────────
func BuscaEspecies(c *gin.Context) {
	query := database.DB.Model(&models.Especie{}).Preload("Sementes")

	// Busca por texto
	if q := c.Query("q"); q != "" {
		termo := "%" + q + "%"
		query = query.Where(
			"nome_cientifico ILIKE ? OR nome_popular ILIKE ? OR descricao ILIKE ?",
			termo, termo, termo,
		)
	}

	// Filtro por bioma
	if bioma := c.Query("bioma"); bioma != "" {
		query = query.Where("bioma ILIKE ?", "%"+bioma+"%")
	}

	// Filtro por família botânica
	if familia := c.Query("familia"); familia != "" {
		query = query.Where("familia ILIKE ?", "%"+familia+"%")
	}

	// Filtro por status de conservação (LC, NT, VU, EN, CR)
	if status := c.Query("status_conservacao"); status != "" {
		query = query.Where("status_conservacao = ?", status)
	}

	// Filtro por validação
	if validado := c.Query("validado"); validado != "" {
		query = query.Where("validada = ?", validado == "true")
	} else {
		// Por padrão retorna só validadas
		query = query.Where("validada = true")
	}

	// Ordenação
	ordem := c.Query("ordem")
	switch ordem {
	case "nome":
		query = query.Order("nome_popular ASC")
	case "recente":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("nome_popular ASC")
	}

	var especies []models.Especie
	if err := query.Find(&especies).Error; err != nil {
		utils.Error(c, 500, "Erro ao buscar espécies")
		return
	}

	utils.Success(c, http.StatusOK, "Busca de espécies concluída", gin.H{
		"total":    len(especies),
		"especies": especies,
	})
}

// ─────────────────────────────────────────
// BUSCA AVANÇADA DE SEMENTES
// GET /api/busca/sementes
// Filtros: q, cor, tamanho, bioma, especie_id, validado
// Público
// ─────────────────────────────────────────
func BuscaSementes(c *gin.Context) {
	query := database.DB.Model(&models.Semente{}).Preload("Especie")

	// Busca por texto
	if q := c.Query("q"); q != "" {
		termo := "%" + q + "%"
		query = query.Where(
			"nome ILIKE ? OR descricao ILIKE ? OR usos_artesanais ILIKE ? OR tecnicas ILIKE ?",
			termo, termo, termo, termo,
		)
	}

	// Filtro por cor
	if cor := c.Query("cor"); cor != "" {
		query = query.Where("cor ILIKE ?", "%"+cor+"%")
	}

	// Filtro por tamanho (pequena, média, grande)
	if tamanho := c.Query("tamanho"); tamanho != "" {
		query = query.Where("tamanho ILIKE ?", "%"+tamanho+"%")
	}

	// Filtro por textura
	if textura := c.Query("textura"); textura != "" {
		query = query.Where("textura ILIKE ?", "%"+textura+"%")
	}

	// Filtro por espécie
	if especieID := c.Query("especie_id"); especieID != "" {
		query = query.Where("especie_id = ?", especieID)
	}

	// Filtro por bioma (via espécie)
	if bioma := c.Query("bioma"); bioma != "" {
		query = query.Joins("JOIN especies ON especies.id = sementes.especie_id").
			Where("especies.bioma ILIKE ?", "%"+bioma+"%")
	}

	// Filtro por validação
	if validado := c.Query("validado"); validado != "" {
		query = query.Where("validada = ?", validado == "true")
	} else {
		query = query.Where("validada = true")
	}

	// Ordenação
	ordem := c.Query("ordem")
	switch ordem {
	case "nome":
		query = query.Order("sementes.nome ASC")
	case "recente":
		query = query.Order("sementes.created_at DESC")
	default:
		query = query.Order("sementes.nome ASC")
	}

	var sementes []models.Semente
	if err := query.Find(&sementes).Error; err != nil {
		utils.Error(c, 500, "Erro ao buscar sementes")
		return
	}

	utils.Success(c, http.StatusOK, "Busca de sementes concluída", gin.H{
		"total":    len(sementes),
		"sementes": sementes,
	})
}

// ─────────────────────────────────────────
// BUSCA GEOGRÁFICA POR PROXIMIDADE
// GET /api/busca/mapa?lat=-15.7&lng=-47.9&raio_km=100
// Retorna registros num raio de X km ao redor do ponto
// Público
// ─────────────────────────────────────────
func BuscaPorProximidade(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	raioStr := c.Query("raio_km")

	if latStr == "" || lngStr == "" {
		utils.Error(c, 400, "Informe 'lat' e 'lng' para busca por proximidade")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil || lat < -90 || lat > 90 {
		utils.Error(c, 400, "Latitude inválida")
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil || lng < -180 || lng > 180 {
		utils.Error(c, 400, "Longitude inválida")
		return
	}

	raioKm := 50.0 // padrão 50km
	if raioStr != "" {
		raioKm, err = strconv.ParseFloat(raioStr, 64)
		if err != nil || raioKm <= 0 {
			utils.Error(c, 400, "Raio inválido. Informe em km. Ex: raio_km=100")
			return
		}
	}

	// Converte raio em graus (aproximação: 1 grau ≈ 111 km)
	raioGraus := raioKm / 111.0

	var registros []models.RegistroSemente
	database.DB.
		Preload("Semente").
		Preload("Semente.Especie").
		Preload("Fotos").
		Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
			lat-raioGraus, lat+raioGraus,
			lng-raioGraus, lng+raioGraus,
		).
		Where("status = ?", models.StatusAprovado).
		Find(&registros)

	// Calcula distância real e filtra precisamente
	type RegistroComDistancia struct {
		models.RegistroSemente
		DistanciaKm float64 `json:"distancia_km"`
	}

	var resultado []RegistroComDistancia
	for _, r := range registros {
		dist := calcularDistanciaKm(lat, lng, r.Latitude, r.Longitude)
		if dist <= raioKm {
			resultado = append(resultado, RegistroComDistancia{
				RegistroSemente: r,
				DistanciaKm:     math.Round(dist*10) / 10, // 1 casa decimal
			})
		}
	}

	utils.Success(c, http.StatusOK, "Registros próximos encontrados", gin.H{
		"centro": gin.H{
			"latitude":  lat,
			"longitude": lng,
		},
		"raio_km":   raioKm,
		"total":     len(resultado),
		"registros": resultado,
	})
}

// ─────────────────────────────────────────
// BUSCA POR ESTADO
// GET /api/busca/estado/:uf
// Retorna espécies, sementes e registros de um estado
// Público
// ─────────────────────────────────────────
func BuscaPorEstado(c *gin.Context) {
	uf := c.Param("uf")
	if len(uf) != 2 {
		utils.Error(c, 400, "Informe a UF com 2 letras. Ex: /api/busca/estado/GO")
		return
	}

	// Registros do estado
	var registros []models.RegistroSemente
	database.DB.
		Preload("Semente").
		Preload("Semente.Especie").
		Preload("Fotos").
		Where("estado ILIKE ? AND status = ?", uf, models.StatusAprovado).
		Order("created_at DESC").
		Find(&registros)

	// Municípios únicos com registros
	type MunicipioCount struct {
		Municipio string `json:"municipio"`
		Total     int    `json:"total"`
	}
	var municipios []MunicipioCount
	database.DB.Model(&models.RegistroSemente{}).
		Select("municipio, COUNT(*) as total").
		Where("estado ILIKE ? AND status = ?", uf, models.StatusAprovado).
		Group("municipio").
		Order("total DESC").
		Scan(&municipios)

	// Sementes encontradas no estado (via registros)
	sementeIDs := []uint{}
	for _, r := range registros {
		sementeIDs = append(sementeIDs, r.SementeID)
	}

	var sementes []models.Semente
	if len(sementeIDs) > 0 {
		database.DB.
			Preload("Especie").
			Where("id IN ?", sementeIDs).
			Find(&sementes)
	}

	utils.Success(c, http.StatusOK, "Dados do estado "+uf, gin.H{
		"estado":               uf,
		"total_registros":      len(registros),
		"total_municipios":     len(municipios),
		"total_sementes":       len(sementes),
		"municipios":           municipios,
		"sementes_encontradas": sementes,
		"registros":            registros,
	})
}

// ─────────────────────────────────────────
// ESTATÍSTICAS GERAIS
// GET /api/busca/estatisticas
// Resumo geral do projeto
// Público
// ─────────────────────────────────────────
func Estatisticas(c *gin.Context) {
	var totalEspecies, totalEspeciesValidadas int64
	var totalSementes, totalSementesValidadas int64
	var totalRegistros, totalRegistrosAprovados int64
	var totalConhecimentos, totalConhecimentosValidados int64
	var totalUsuarios int64

	database.DB.Model(&models.Especie{}).Count(&totalEspecies)
	database.DB.Model(&models.Especie{}).Where("validada = true").Count(&totalEspeciesValidadas)
	database.DB.Model(&models.Semente{}).Count(&totalSementes)
	database.DB.Model(&models.Semente{}).Where("validada = true").Count(&totalSementesValidadas)
	database.DB.Model(&models.RegistroSemente{}).Count(&totalRegistros)
	database.DB.Model(&models.RegistroSemente{}).Where("status = ?", models.StatusAprovado).Count(&totalRegistrosAprovados)
	database.DB.Model(&models.ConhecimentoTradicional{}).Count(&totalConhecimentos)
	database.DB.Model(&models.ConhecimentoTradicional{}).Where("validado = true").Count(&totalConhecimentosValidados)
	database.DB.Model(&models.Usuario{}).Where("ativo = true").Count(&totalUsuarios)

	// Estados com mais registros
	type EstadoCount struct {
		Estado string `json:"estado"`
		Total  int    `json:"total"`
	}
	var estadosAtivos []EstadoCount
	database.DB.Model(&models.RegistroSemente{}).
		Select("estado, COUNT(*) as total").
		Where("status = ? AND estado != ''", models.StatusAprovado).
		Group("estado").
		Order("total DESC").
		Limit(10).
		Scan(&estadosAtivos)

	// Biomas com mais espécies
	type BiomaCount struct {
		Bioma string `json:"bioma"`
		Total int    `json:"total"`
	}
	var biomasAtivos []BiomaCount
	database.DB.Model(&models.Especie{}).
		Select("bioma, COUNT(*) as total").
		Where("validada = true AND bioma != ''").
		Group("bioma").
		Order("total DESC").
		Scan(&biomasAtivos)

	utils.Success(c, http.StatusOK, "Estatísticas do Mapa das Sementes do Brasil", gin.H{
		"especies": gin.H{
			"total":     totalEspecies,
			"validadas": totalEspeciesValidadas,
		},
		"sementes": gin.H{
			"total":     totalSementes,
			"validadas": totalSementesValidadas,
		},
		"registros": gin.H{
			"total":     totalRegistros,
			"aprovados": totalRegistrosAprovados,
		},
		"conhecimentos": gin.H{
			"total":     totalConhecimentos,
			"validados": totalConhecimentosValidados,
		},
		"usuarios_ativos": totalUsuarios,
		"estados_ativos":  estadosAtivos,
		"biomas_ativos":   biomasAtivos,
	})
}

// ─────────────────────────────────────────
// FUNÇÃO AUXILIAR: Fórmula de Haversine
// Calcula distância em km entre dois pontos geográficos
// ─────────────────────────────────────────
func calcularDistanciaKm(lat1, lng1, lat2, lng2 float64) float64 {
	const raioTerraKm = 6371.0

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return raioTerraKm * c
}
