package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/utils/exception"
	jwtutils "github.com/rms-diego/image-processor/internal/utils/jwt"
)

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.Error(exception.New("token is required", http.StatusUnauthorized))
			c.Abort()
			return
		}

		user, err := jwtutils.NewJwtUtils().ValidateAndDecodeToken(token)
		if err != nil {
			c.Error(exception.New("token expired", http.StatusUnauthorized))
			c.Abort()
			return
		}

		c.Set("user", *user)
		c.Next()
	}
}
