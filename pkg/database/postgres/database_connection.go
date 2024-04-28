package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDb(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // turn off logger
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successly connected to Postgres!")

	return db, nil
}
