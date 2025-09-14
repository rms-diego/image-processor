package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/config"
	"github.com/rms-diego/image-processor/internal/database"
	"github.com/rms-diego/image-processor/internal/routes"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}

	if err := database.Init(); err != nil {
		panic(err)
	}

	app := gin.Default()

	routes.Init(app.Group("/"))
	app.Run(":" + config.Env.PORT)
}
