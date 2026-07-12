package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware permite que el frontend (en otro dominio, ej. Vercel) consuma
// esta API (en Railway). origin puede ser "*" o un dominio específico.
func CORSMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
