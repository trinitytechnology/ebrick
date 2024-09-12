package environment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EnvironmentHandler interface {
	CreateEnvironment(ctx *gin.Context)
}

type envHandler struct {
	service EnvironmentService
	log     *zap.Logger
}

func NewEnvironmentHandler(service EnvironmentService, log *zap.Logger) EnvironmentHandler {
	return &envHandler{
		service: service,
		log:     log,
	}
}

// CreateEnvironment implements EnvironmentHandler.
func (h *envHandler) CreateEnvironment(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ten := Environment{Name: request.Name}
	h.log.Info("Creating Environment", zap.String("name", ten.Name))
	h.service.CreateEnvironment(c.Request.Context(), ten)
	c.JSON(http.StatusOK, gin.H{"message": "Environment Created"})
}
