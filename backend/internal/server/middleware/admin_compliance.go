package middleware

import (
	"github.com/gin-gonic/gin"
)

func AdminComplianceGuard(settingService interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func isAdminComplianceBypassPath(path string) bool {
	return false
}
