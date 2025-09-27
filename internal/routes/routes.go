package routes

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/modules/auth"
	"github.com/rms-diego/image-processor/internal/modules/image"
)

func Init(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "server is running 🚀"})
	})

	r.GET("/docs", func(c *gin.Context) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL:  "docs/swagger.yml",
			DarkMode: true,
		})

		if err != nil {
			c.Error(err)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	})

	auth.RouteInit(r.Group("/auth"))
	image.RouteInit(r.Group("/images"))
}
