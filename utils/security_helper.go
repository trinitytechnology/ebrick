package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func ExtractClaimsFromContext(c *gin.Context) (map[string]any, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, errors.New("Claims not found in context")
	}
	claimsMap, ok := claims.(map[string]any)
	if !ok {
		return nil, errors.New("Claims have wrong type")
	}
	return claimsMap, nil
}
