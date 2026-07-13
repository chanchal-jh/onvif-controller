package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	onvifservice "onvif-camera-controller/internal/onvif"
)

// StreamRequest represents stream request
type StreamRequest struct {
	IP       string `json:"ip" binding:"required,ip"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GetStreams handles stream requests
func GetStreams(c *gin.Context) {

	var req StreamRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		log.Printf("ERROR GetStreams: invalid request from client=%s: %v", c.ClientIP(), err)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	log.Printf("INFO GetStreams: request received for camera=%s", req.IP)

	streams, err := onvifservice.GetStreamProfiles(
		req.IP,
		req.Username,
		req.Password,
	)
	if err != nil {

		log.Printf("ERROR GetStreams: camera=%s: %v", req.IP, err)

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

	log.Printf("INFO GetStreams: successfully retrieved %d stream profile(s) for camera=%s", len(streams), req.IP)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"profiles": streams,
	})
}
