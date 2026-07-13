package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	onvifservice "onvif-camera-controller/internal/onvif"
)

// HandlePTZAction handles PTZ action requests
func HandlePTZAction(c *gin.Context) {

	var req onvifservice.ActionRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		log.Printf("ERROR HandlePTZAction: invalid request from client=%s: %v", c.ClientIP(), err)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	log.Printf("INFO HandlePTZAction: action=%s camera=%s", req.Action, req.IP)

	err := onvifservice.HandleAction(req)
	if err != nil {

		log.Printf("ERROR HandlePTZAction: action=%s camera=%s: %v", req.Action, req.IP, err)

		if strings.Contains(err.Error(), "ter:NotAuthorized") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authentication failed:invalid username or password",
			})
			return
		}

		c.JSON(onvifservice.GetHTTPStatus(err), gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	log.Printf("INFO HandlePTZAction: action=%s camera=%s completed successfully", req.Action, req.IP)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Action completed successfully",
	})
}
