package web

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	setupProbeRoute(router)
	return router
}
