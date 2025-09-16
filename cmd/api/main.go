package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/database"
	"github.com/rms-diego/image-processor/internal/middleware"
	"github.com/rms-diego/image-processor/internal/routes"
	"github.com/rms-diego/image-processor/pkg/config"
	"github.com/rms-diego/image-processor/pkg/gateway"
)

func main() {
	if err := config.InitServerCfg(); err != nil {
		panic(err)
	}

	if err := config.InitAwsCfg(); err != nil {
		panic(err)
	}

	if err := database.Init(); err != nil {
		panic(err)
	}

	if err := gateway.InitSQS(); err != nil {
		panic(err)
	}

	app := gin.Default()
	app.Use(middleware.ErrorHandler())

	routes.Init(app.Group("/"))
	app.Run(":" + config.ServerEnv.PORT)
}
