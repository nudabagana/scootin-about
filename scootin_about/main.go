package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port, found := os.LookupEnv(`PORT`)
	if !found {
		port = "8090"
	}

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	router.GET(`/health`, HealthGET)
	router.Run(":" + port)
}

func HealthGET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "1.0.0",
	})
}
