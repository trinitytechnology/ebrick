package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// setupProbeRoute sets up the probe routes for the application
func setupProbeRoute(router *gin.Engine) {

	// Health Check Endpoint
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Readiness Check Endpoint
	router.GET("/ready", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}
