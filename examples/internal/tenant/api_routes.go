package tenant

import (
	"github.com/gin-gonic/gin"
)

func setupApiRoutes(handler TenantHandler, router *gin.Engine) {
	router.POST("/api/tenants", handler.CreateTenant)
}
