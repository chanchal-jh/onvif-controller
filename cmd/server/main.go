package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"onvif-camera-controller/internal/api"
	"onvif-camera-controller/internal/config"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	cfg := config.LoadConfig()

	router := api.SetupRoutes()

	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ONVIF Camera API running on :%s\n", cfg.Port)

	err = router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal(err)
	}
}
