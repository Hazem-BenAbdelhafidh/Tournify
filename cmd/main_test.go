package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/Hazem-BenAbdelhafidh/Tournify/db"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	db.ConnectToDb()
	err := loadEnv()
	if err != nil {
		fmt.Println("Error loading .env file in test", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Perform any teardown if needed
	os.Exit(code)
}

func loadEnv() error {
	// Explicitly specify the path to the .env file
	cwd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", cwd)
	return godotenv.Load(cwd + "/../.env")
}
