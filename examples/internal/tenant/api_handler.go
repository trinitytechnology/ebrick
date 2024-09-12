package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TenantHandler interface {
	CreateTenant(ctx *gin.Context)
}

type tenantHandler struct {
	service TenantService
	log     *zap.Logger
}

func NewTenantHandler(service TenantService, log *zap.Logger) TenantHandler {
	return &tenantHandler{
		service: service,
		log:     log,
	}
}

// CreateTenant implements TenantHandler.
func (t *tenantHandler) CreateTenant(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ten := Tenant{Name: request.Name}
	t.log.Info("Creating Tenant", zap.String("name", ten.Name))
	t.service.CreateTenant(c.Request.Context(), ten)
	c.JSON(http.StatusOK, gin.H{"message": "Tenant Created"})
}
