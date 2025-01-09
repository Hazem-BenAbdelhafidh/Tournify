package main

import (
	"github.com/Hazem-BenAbdelhafidh/Tournify/api"
)

func main() {
	router := api.SetupRouter()
	router.Run(":8000")
}
