package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AUTH_TOKEN = "Bearer EXAMPLE"
)

func Init(router *gin.Engine) {
	scooterGroup := router.Group("/scooter")
	scooterGroup.Use(AuthMiddleware())
	scooterGroup.PUT("/:uuid/start-trip", scooterHandler.tripStart)
	scooterGroup.PUT("/:uuid/end-trip", scooterHandler.endTrip)
	scooterGroup.POST("/:uuid/report-location", scooterHandler.reportLocation)

	clientGroup := router.Group("/client")
	clientGroup.Use(AuthMiddleware())
	clientGroup.GET("/scooters", clientHandler.getScooters)

	testingGroup := router.Group("/testing")
	testingGroup.GET("/scooters", testingHandler.getAllScooters)
	testingGroup.GET("/scooters/:uuid", testingHandler.getScooter)
	testingGroup.POST("/scooters", testingHandler.createScooter)
	testingGroup.DELETE("/scooters/:uuid", testingHandler.deleteScooter)

	testingGroup.GET("/clients", testingHandler.GetAllClients)
	testingGroup.GET("/clients/:uuid", testingHandler.GetClient)
	testingGroup.POST("/clients", testingHandler.CreateClient)
	testingGroup.DELETE("/clients/:uuid", testingHandler.DeleteClient)

	testingGroup.GET("/locations", testingHandler.getAllLocations)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		if tokenString != AUTH_TOKEN {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
