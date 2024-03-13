package main

import (
	"fmt"
	"go-service/app"
	route "go-service/internal"

	"gorm.io/driver/postgres"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	app, err := app.NewApp(db)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	route.Route(r, app)
	r.Run(":8080")
}
