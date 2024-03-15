package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthMock() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", "login_user_id")
		c.Set("organization_id", "login_organization_id")
		c.Next()
	}
}
