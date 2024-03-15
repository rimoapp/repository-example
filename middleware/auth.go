package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthMock() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginUserID := c.Query("login_user_id")
		if loginUserID == "" {
			loginUserID = "login_user_id"
		}
		c.Set("user_id", loginUserID)
		c.Next()
	}
}
