package main

import (
	"github.com/Mikael88/go-mygram/config"
	"github.com/Mikael88/go-mygram/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load environment variables")
	}

	config.InitDB()

	r := gin.Default()

	routes.SetupRoutes(r)

	r.Run()
}
