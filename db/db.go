package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/Hazem-BenAbdelhafidh/Tournify/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDb() *gorm.DB {
	var dsn, dbUser, password, port, database string

	if testing.Testing() {
		dbUser = "hazem2"
		password = "test1234"
		port = "5433"
		database = "tournify_test"
		dsn = fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable", dbUser, password, database, port)
	} else {
		dbUser = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		port = os.Getenv("DB_PORT")
		database = os.Getenv("DB_DATABASE")
		dsn = fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable", dbUser, password, database, port)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&entities.User{}, &entities.Tournament{})
	if err != nil {
		fmt.Println("error while migrating : ", err.Error())
	}

	fmt.Println("connected to db")

	return db
}
