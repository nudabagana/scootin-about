package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nudabagana/scootin-about/api"
	"github.com/nudabagana/scootin-about/data"
	_ "github.com/nudabagana/scootin-about/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	api.Init(router)
	err := data.Init()
	if err != nil {
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	router.GET(`/health`, HealthGET)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	<-quit
	data.Stop()
}

func HealthGET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "1.0.0",
	})
}
