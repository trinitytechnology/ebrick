package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/linkifysoft/ebrick/config"
	"github.com/linkifysoft/ebrick/logger"
	"go.uber.org/zap"
)

var verifier *oidc.IDTokenVerifier
var provider *oidc.Provider

func InitOIDC() {
	cfg := config.GetConfig().Oidc
	logger := logger.DefaultLogger
	if cfg.Enable {
		logger.Info("Setting OIDC...", zap.String("issuer", cfg.Issuer))
		ctx := context.Background()
		var err error
		provider, err = oidc.NewProvider(ctx, cfg.Issuer)
		if err != nil {
			panic("Failed to get provider: " + err.Error())
		}
		verifier = provider.Verifier(&oidc.Config{
			ClientID:          cfg.ClientId,
			SkipClientIDCheck: true,
			SkipExpiryCheck:   false,
			SkipIssuerCheck:   false,
		})
	}
}

func OIDCAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig().Oidc
		if cfg.Enable {
			ctx := context.Background()
			tokenString := c.GetHeader("Authorization")
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
				c.Abort()
				return
			}
			tokenString = strings.Trim(tokenString, "Bearer")
			idToken, err := verifier.Verify(ctx, tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}
			var claims map[string]any
			if err := idToken.Claims(&claims); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
				c.Abort()
				return
			}
			c.Set("claims", claims)
		}
		c.Next()
	}
}
