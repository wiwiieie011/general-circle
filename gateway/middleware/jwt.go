package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"gateway/jwtutil"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			return
		}

		claims, err := jwtutil.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		c.Request.Header.Set(
			"X-User-Id",
			strconv.FormatUint(uint64(claims.UserID), 10),
		)
		c.Request.Header.Set(
			"X-User-Role",
			claims.Role,
		)

		c.Next()
	}
}
