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

// SalvarImagem salva um arquivo de imagem em disco e retorna o caminho relativo.
// pasta: subpasta dentro de ./uploads (ex: "avatars", "especies")
func SalvarImagem(c *gin.Context, arquivo *multipart.FileHeader, pasta string) (string, error) {
	// Garante que a pasta existe
	dir := filepath.Join("uploads", pasta)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("erro ao criar diretório: %w", err)
	}

	// Gera nome único para evitar colisões
	ext := strings.ToLower(filepath.Ext(arquivo.Filename))
	nomeArquivo := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	caminho := filepath.Join(dir, nomeArquivo)

	if err := c.SaveUploadedFile(arquivo, caminho); err != nil {
		return "", fmt.Errorf("erro ao salvar arquivo: %w", err)
	}

	// Retorna caminho relativo para armazenar no banco
	return "/" + filepath.ToSlash(caminho), nil
}

// DeletarImagem remove um arquivo de imagem do disco.
// caminhoRelativo: caminho como armazenado no banco (ex: "/uploads/avatars/foto.jpg")
func DeletarImagem(caminhoRelativo string) error {
	if caminhoRelativo == "" {
		return nil
	}

	// Remove a barra inicial se houver
	caminho := strings.TrimPrefix(caminhoRelativo, "/")

	if err := os.Remove(caminho); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("erro ao deletar imagem: %w", err)
	}

	return nil
}
