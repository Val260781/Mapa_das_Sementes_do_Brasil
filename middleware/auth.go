package middleware

import (
	"strings"

	"mapa-sementes-brasil/config"
	"mapa-sementes-brasil/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, 401, "Token não informado")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error(c, 401, "Formato de token inválido")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.App.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			utils.Error(c, 401, "Token inválido ou expirado")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.Error(c, 401, "Token inválido")
			c.Abort()
			return
		}

		c.Set("usuario_id", uint(claims["usuario_id"].(float64)))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		roleStr := role.(string)

		for _, r := range roles {
			if r == roleStr {
				c.Next()
				return
			}
		}

		utils.Error(c, 403, "Acesso não autorizado para este perfil")
		c.Abort()
	}
}
