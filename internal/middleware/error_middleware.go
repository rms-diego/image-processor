package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/utils/exception"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		last := c.Errors.Last().Err
		if ex, ok := last.(*exception.AppError); ok {

			c.AbortWithStatusJSON(
				ex.Code,
				gin.H{"message": ex.Message, "code": ex.Code},
			)

			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"message": last.Error(), "code": http.StatusInternalServerError},
		)
	}
}
