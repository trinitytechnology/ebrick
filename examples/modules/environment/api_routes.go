package environment

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(handler EnvironmentHandler, router *gin.Engine) {
	router.POST("/api/envs", handler.CreateEnvironment)
}
