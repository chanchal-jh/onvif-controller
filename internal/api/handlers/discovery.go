package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	onvifservice "onvif-camera-controller/internal/onvif"
)

// DiscoverCameras handles ONVIF camera discovery
func DiscoverCameras(c *gin.Context) {

	// Optional query parameter:
	// Example:
	// /api/v1/discovery?interface=eth0
	networkInterface := c.Query("interface")

	var (
		cameras []onvifservice.CameraDevice
		err     error
	)

	// Discover using specific interface if provided
	if networkInterface != "" {

		cameras, err = onvifservice.DiscoverCamerasWithInterface(
			networkInterface,
		)

	} else {

		cameras, err = onvifservice.DiscoverCameras()
	}

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(cameras),
		"devices": cameras,
	})
}

// ListInterfaces returns available network interfaces
func ListInterfaces(c *gin.Context) {

	interfaces, err := onvifservice.ListInterfaces()
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"interfaces": interfaces,
	})
}
