package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	fmt.Println("Server running on port 8000")
	router.Run(":8000")
}
