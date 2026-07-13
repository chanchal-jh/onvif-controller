package handlers

import (
	"net/http"

	onvifservice "onvif-camera-controller/internal/onvif"

	"github.com/gin-gonic/gin"
)

type PTZCapabilitiesRequest struct {
	IP       string `json:"ip" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GetPTZCapabilities(c *gin.Context) {

	var req PTZCapabilitiesRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"error":   err.Error(),
			},
		)
		return
	}

	_, err := onvifservice.ConnectCamera(
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

	capabilities, err := onvifservice.GetPTZCapabilities(
		req.IP,
		req.Username,
		req.Password,
	)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success":      true,
			"capabilities": capabilities,
		},
	)
}
