package middleware

import (
	"net/http"
	"strings"

	"github.com/example/ag-directory-manager/internal/services"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserID      = "userID"
	ContextUsername    = "username"
	ContextPermissions = "permissions"
	ContextOUScopes    = "ouScopes"
)

func AuthRequired(auth *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if authz == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
			return
		}
		parts := strings.SplitN(authz, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}
		claims, err := auth.ParseAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			return
		}
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextPermissions, claims.Permissions)
		c.Set(ContextOUScopes, claims.OUScopes)
		c.Next()
	}
}
