package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/database"
	authRepository "github.com/rms-diego/image-processor/internal/modules/auth/auth_repository"
	"github.com/rms-diego/image-processor/internal/utils/exception"
	jwtutils "github.com/rms-diego/image-processor/internal/utils/jwt"
)

func AuthHandler() gin.HandlerFunc {
	r := authRepository.NewRepository(database.DB)

	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		if bearerToken == "" {
			c.Error(exception.New("token is required", http.StatusUnauthorized))
			c.Abort()
			return
		}

		token := strings.Split(bearerToken, " ")[1]
		jwtUtils := jwtutils.NewJwtUtils()
		tokenDecoded, err := jwtUtils.ValidateAndDecodeToken(token)

		if err != nil {
			c.Error(exception.New("token malformed or expired", http.StatusUnauthorized))
			c.Abort()
			return
		}

		u, err := r.GetUserByID(tokenDecoded.ID)

		if err != nil {
			c.Error(exception.New(err.Error(), http.StatusUnauthorized))
			c.Abort()
			return
		}

		if u == nil {
			c.Error(exception.New("user not found", http.StatusNotFound))
			c.Abort()
			return
		}

		c.Set("user", *tokenDecoded)
		c.Next()
	}
}
