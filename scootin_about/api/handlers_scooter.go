package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nudabagana/scootin-about/data"
)

type _scooterHandler struct{}

var scooterHandler = _scooterHandler{}

// tripStart godoc
// @Summary Start a trip for a scooter
// @Description Set the scooter's status to occupied (start trip).
// @Tags scooter
// @Accept json
// @Produce json
// @Param uuid path string true "Scooter UUID"
// @Success 200 {object} interface{}
// @Failure 404 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /scooter/{uuid}/start-trip [put]
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

// endTrip godoc
// @Summary End a trip for a scooter
// @Description Set the scooter's status to not occupied (end trip).
// @Tags scooter
// @Accept json
// @Produce json
// @Param uuid path string true "Scooter UUID"
// @Success 200 {object} interface{}
// @Failure 404 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /scooter/{uuid}/end-trip [put]
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

// reportLocation godoc
// @Summary Report the location of a scooter
// @Description Updates the location of the specified scooter with latitude and longitude.
// @Tags scooter
// @Accept json
// @Produce json
// @Param uuid path string true "Scooter UUID"
// @Param location body ReportLocationRequest true "Location data"
// @Success 201 {object} data.Location
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /scooter/{uuid}/report-location [post]
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
