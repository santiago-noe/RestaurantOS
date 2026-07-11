package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"restaurantos/internal/auth"
)

const claimsKey = "claims"

func JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ValidateJWT(tokenStr, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}
		c.Set(claimsKey, claims)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get(claimsKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sin autenticación"})
			return
		}
		claims := raw.(*auth.Claims)
		for _, r := range roles {
			if claims.Rol == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acceso denegado"})
	}
}

func GetClaims(c *gin.Context) *auth.Claims {
	raw, _ := c.Get(claimsKey)
	claims, _ := raw.(*auth.Claims)
	return claims
}
