package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUUIDParam(c *gin.Context, param string) (uuid.UUID, bool) {
	p := c.Param(param)
	pUUID, err := uuid.Parse(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse " + param})
		return uuid.Nil, false
	}
	return pUUID, true
}

// Get Tenant by get "tenant_id" param return uuid.Nil if not found or error
func GetTenantUUID(c *gin.Context) uuid.UUID {
	p := c.Param("tenant_id")
	pUUID, err := uuid.Parse(p)
	if err != nil {
		return uuid.Nil
	}
	return pUUID
}
