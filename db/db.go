package db

import (
	"fmt"
	"testing"

	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/tournament"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDb() *gorm.DB {
	var dsn string
	if testing.Testing() {
		dsn = "host=localhost user=hazem2 password=test1234 dbname=tournify_test port=5433 sslmode=disable"
	} else {
		dsn = "host=localhost user=hazem password=test123 dbname=tournify port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&user.User{}, &tournament.Tournament{})
	if err != nil {
		fmt.Println("error while migrating : ", err.Error())
	}

	fmt.Println("connected to db")

	return db
}
