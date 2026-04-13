package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidationFailed(c *gin.Context, fields map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":  "validation failed",
		"fields": fields,
	})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}
