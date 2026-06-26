package middleware

import (
	"net/http"

	"github.com/example/ag-directory-manager/internal/permissions"
	"github.com/gin-gonic/gin"
)

func RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawPermissions, exists := c.Get(ContextPermissions)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing permission context"})
			return
		}
		permissionList, ok := rawPermissions.([]string)
		if !ok || !permissions.HasPermission(permissionList, code) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		c.Next()
	}
}
