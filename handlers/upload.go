package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Pasta onde as imagens serão salvas
const pastaUpload = "uploads"

// Tipos de imagem permitidos
var tiposPermitidos = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// SalvarImagem salva um arquivo de imagem em disco e retorna a URL pública
func SalvarImagem(c *gin.Context, arquivo *multipart.FileHeader, subpasta string) (string, error) {
	// Verifica extensão
	ext := strings.ToLower(filepath.Ext(arquivo.Filename))
	if !tiposPermitidos[ext] {
		return "", fmt.Errorf("tipo de arquivo não permitido: %s (use jpg, jpeg, png ou webp)", ext)
	}

	// Verifica tamanho (máx 5MB)
	if arquivo.Size > 5*1024*1024 {
		return "", fmt.Errorf("imagem muito grande (máximo 5MB)")
	}

	// Cria a pasta se não existir
	caminhoPasta := filepath.Join(pastaUpload, subpasta)
	if err := os.MkdirAll(caminhoPasta, os.ModePerm); err != nil {
		return "", fmt.Errorf("erro ao criar pasta de upload: %v", err)
	}

	// Gera nome único para o arquivo
	nomeArquivo := fmt.Sprintf("%d_%s%s",
		time.Now().UnixNano(),
		strings.ReplaceAll(arquivo.Filename[:len(arquivo.Filename)-len(ext)], " ", "_"),
		ext,
	)

	caminhoCompleto := filepath.Join(caminhoPasta, nomeArquivo)

	// Salva o arquivo
	if err := c.SaveUploadedFile(arquivo, caminhoCompleto); err != nil {
		return "", fmt.Errorf("erro ao salvar imagem: %v", err)
	}

	// Retorna URL pública (relativa ao servidor)
	urlPublica := fmt.Sprintf("/uploads/%s/%s", subpasta, nomeArquivo)
	return urlPublica, nil
}

// DeletarImagem remove um arquivo de imagem do disco
func DeletarImagem(urlImagem string) error {
	// Remove a barra inicial "/uploads/..." → "uploads/..."
	caminho := strings.TrimPrefix(urlImagem, "/")
	if err := os.Remove(caminho); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("erro ao deletar imagem: %v", err)
	}
	return nil
}
