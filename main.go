package main

import (
	"go-service/app"
	route "go-service/internal"
	"go-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	logger := logger.NewLogger()

	app, err := app.NewApp(logger)
	if err != nil {
		logger.LogError(err.Error(), nil)
		return
	}

	route.Route(r, app)
	r.Run(":8080")
}
