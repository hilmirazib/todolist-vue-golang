package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from panics and returns a JSON error response.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("panic recovered: %v", recovered)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "internal server error",
			"message": "an unexpected error occurred",
		})
	})
}
