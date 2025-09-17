package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/modules/auth"
	"github.com/rms-diego/image-processor/internal/modules/image"
)

func Init(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "server is running 🚀"})
	})

	auth.RouteInit(r.Group("/auth"))
	image.RouteInit(r.Group("/images"))
}
