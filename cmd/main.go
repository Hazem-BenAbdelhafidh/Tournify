package main

import (
	"fmt"
	"os"

	"github.com/Hazem-BenAbdelhafidh/Tournify/api"
	"github.com/Hazem-BenAbdelhafidh/Tournify/config"
)

func main() {
	config.Config()
	port := os.Getenv("PORT")
	router := api.SetupRouter()
	router.Run(fmt.Sprintf(":%s", port))
}
