package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/database"
	authRepository "github.com/rms-diego/image-processor/internal/modules/auth/repository"
	"github.com/rms-diego/image-processor/internal/utils/exception"
	jwtutils "github.com/rms-diego/image-processor/internal/utils/jwt"
)

var r = authRepository.NewRepository(database.Db)

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.Error(exception.New("token is required", http.StatusUnauthorized))
			c.Abort()
			return
		}

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
