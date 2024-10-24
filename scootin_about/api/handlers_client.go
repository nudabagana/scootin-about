package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nudabagana/scootin-about/data"
)

type _clientHandler struct{}

var clientHandler = _clientHandler{}

func (_clientHandler) getScooters(c *gin.Context) {
	topLeftLatStr := c.Query("top_left_lat")
	topLeftLongStr := c.Query("top_left_long")
	bottomRightLatStr := c.Query("bottom_right_lat")
	bottomRightLongStr := c.Query("bottom_right_long")
	occupiedStr := c.Query("occupied")

	topLeftLat, err := strconv.ParseFloat(topLeftLatStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid top_left_lat"})
		return
	}

	topLeftLong, err := strconv.ParseFloat(topLeftLongStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid top_left_long"})
		return
	}

	bottomRightLat, err := strconv.ParseFloat(bottomRightLatStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bottom_right_lat"})
		return
	}

	bottomRightLong, err := strconv.ParseFloat(bottomRightLongStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bottom_right_long"})
		return
	}

	var occupied *bool
	if occupiedStr != "" {
		b := (occupiedStr == "true")
		occupied = &b
	}

	scooters, err := data.ScooterRepo.GetScootersInSquare(topLeftLat, topLeftLong, bottomRightLat, bottomRightLong, occupied)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scooters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  scooters,
		"count": len(scooters),
	})
}
