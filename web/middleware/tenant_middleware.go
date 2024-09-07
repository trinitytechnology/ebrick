package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ValidateTenantID validates the tenant ID in the request
func ValidateTenantID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tenantId := ctx.Param("tenant_id")

		if tenantId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID is required"})
			ctx.Abort()
			return
		}

		// Validate UUID format
		_, err := uuid.Parse(tenantId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tenant ID"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
