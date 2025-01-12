package main

import (
	"fmt"
	"os"

	_ "github.com/Hazem-BenAbdelhafidh/Tournify/docs"

	"github.com/Hazem-BenAbdelhafidh/Tournify/api"
	"github.com/Hazem-BenAbdelhafidh/Tournify/config"
)

// @title	Tournify REST API
func main() {
	config.Config()
	port := os.Getenv("PORT")
	router := api.SetupRouter()
	router.Run(fmt.Sprintf(":%s", port))
}
