package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/config"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}

	// r.GET("/", func(c *gin.Context) {
	// 	c.String(200, "Hello, World!")
	// })
	r := gin.Default()
	r.Run(":" + config.Env.PORT)
}
