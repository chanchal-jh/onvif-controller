package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	onvifservice "onvif-camera-controller/internal/onvif"
)

// ConnectRequest represents connect camera request
type ConnectRequest struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ConnectCamera handles camera connection requests
func ConnectCamera(c *gin.Context) {

	var req ConnectRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	cameraInfo, err := onvifservice.ConnectCamera(
		req.IP,
		req.Username,
		req.Password,
	)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"camera":  cameraInfo,
	})
}
