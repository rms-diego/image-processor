package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/database"
	"github.com/rms-diego/image-processor/internal/middleware"
	"github.com/rms-diego/image-processor/internal/routes"
	"github.com/rms-diego/image-processor/pkg/config"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}

	if err := database.Init(); err != nil {
		panic(err)
	}

	app := gin.Default()
	app.Use(middleware.ErrorHandler())

	routes.Init(app.Group("/"))
	app.Run(":" + config.Env.PORT)
}
