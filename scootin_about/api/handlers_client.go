package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nudabagana/scootin-about/data"
)

type _clientHandler struct{}

var clientHandler = _clientHandler{}

type ScootersResponse struct {
	Data  []data.ScooterWithLocation `json:"data"`
	Count int                        `json:"count"`
}

// getScooters godoc
// @Summary Retrieve scooters in a specified area
// @Description Fetches scooters within a defined rectangular area using latitude and longitude boundaries.
// @Tags client
// @Accept json
// @Produce json
// @Param top_left_lat query float64 true "Top left latitude"
// @Param top_left_long query float64 true "Top left longitude"
// @Param bottom_right_lat query float64 true "Bottom right latitude"
// @Param bottom_right_long query float64 true "Bottom right longitude"
// @Param occupied query bool false "Filter by occupancy status"
// @Success 200 {object} ScootersResponse
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /client/scooters [get]
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

	c.JSON(http.StatusOK, ScootersResponse{
		Data:  scooters,
		Count: len(scooters),
	})
}
