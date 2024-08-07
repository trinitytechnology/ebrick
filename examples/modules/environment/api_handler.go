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
}

func NewEnvironmentHandler(service EnvironmentService) EnvironmentHandler {
	return &envHandler{
		service: service,
	}
}

// CreateEnvironment implements EnvironmentHandler.
func (t *envHandler) CreateEnvironment(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ten := Environment{Name: request.Name}
	log.Info("Creating Environment", zap.String("name", ten.Name))
	t.service.CreateEnvironment(c.Request.Context(), ten)
	c.JSON(http.StatusOK, gin.H{"message": "Environment Created"})
}
