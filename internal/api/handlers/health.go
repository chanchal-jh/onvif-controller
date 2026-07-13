package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns API health status
func HealthCheck(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ONVIF Camera API is running",
	})
}
