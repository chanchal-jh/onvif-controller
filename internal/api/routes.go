package api

import (
	"github.com/gin-gonic/gin"

	"onvif-camera-controller/internal/api/handlers"
	"onvif-camera-controller/internal/config"
)

// SetupRoutes configures API routes
func SetupRoutes() *gin.Engine {

	cfg := config.LoadConfig()

	router := gin.Default()

	// Health endpoint without auth
	router.GET("/health", handlers.HealthCheck)

	// Protected API group
	auth := router.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		cfg.APIUsername: cfg.APIPassword,
	}))

	{
		auth.GET("/discovery", handlers.DiscoverCameras)
		auth.GET("/interfaces", handlers.ListInterfaces)

		auth.POST("/camera/connect", handlers.ConnectCamera)
		auth.POST("/camera/streams", handlers.GetStreams)
		auth.POST("/camera/ptz-capabilities", handlers.GetPTZCapabilities)
		auth.POST("/action", handlers.HandlePTZAction)
	}

	return router
}
