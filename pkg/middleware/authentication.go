package middleware

import (
	"errors"
	"os"

	"apiDentalClinic/pkg/web"
	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	requiredToken := os.Getenv("TOKEN")

	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			c.Abort()
			return
		}

		if token != requiredToken {
			web.Failure(c, 401, errors.New("Invalid Token"))
			c.Abort()
			return
		}

		c.Next()
	}

}
