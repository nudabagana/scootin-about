package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nudabagana/scootin-about/data"
)

type _scooterHandler struct{}

var scooterHandler = _scooterHandler{}

func (_scooterHandler) tripStart(ginCtx *gin.Context) {
	uuid := ginCtx.Param("uuid")

	err := data.ScooterRepo.SetOccupied(uuid, true)
	if err != nil {
		if err.Error() == "scooter not found" {
			ginCtx.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start trip"})
		}
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{"message": "Trip started", "scooter_uuid": uuid})
}

func (_scooterHandler) endTrip(ginCtx *gin.Context) {
	uuid := ginCtx.Param("uuid")

	err := data.ScooterRepo.SetOccupied(uuid, false)
	if err != nil {
		if err.Error() == "scooter not found" {
			ginCtx.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end trip"})
		}
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{"message": "Trip ended", "scooter_uuid": uuid})
}

type ReportLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (_scooterHandler) reportLocation(ginCtx *gin.Context) {
	var req ReportLocationRequest

	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	scooterUuid := ginCtx.Param("uuid")

	location, err := data.LocationRepo.Add(req.Latitude, req.Longitude, scooterUuid)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to report location"})
		return
	}

	ginCtx.JSON(http.StatusCreated, location)
}
