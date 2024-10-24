package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nudabagana/scootin-about/data"
)

type _testingHandler struct{}

var testingHandler = _testingHandler{}

func (_testingHandler) getAllScooters(ginCtx *gin.Context) {
	scooters, err := data.ScooterRepo.GetAll()

	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scooters"})
		return
	}
	ginCtx.JSON(http.StatusOK, scooters)
}

func (h *_testingHandler) getScooter(ginCtx *gin.Context) {
	uuid := ginCtx.Param("uuid")

	scooter, err := data.ScooterRepo.Get(uuid)
	if err != nil {
		ginCtx.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		return
	}

	ginCtx.JSON(http.StatusOK, scooter)
}

type CreateScooterRequest struct {
	Name     string `json:"name"`
	Occupied bool   `json:"occupied"`
}

func (h *_testingHandler) createScooter(c *gin.Context) {
	var req CreateScooterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	scooter, err := data.ScooterRepo.CreateScooter(req.Name, req.Occupied)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create scooter"})
		return
	}

	c.JSON(http.StatusCreated, scooter)
}

func (h *_testingHandler) deleteScooter(ginCtx *gin.Context) {
	uuid := ginCtx.Param("uuid")

	err := data.ScooterRepo.Delete(uuid)
	if err != nil {
		ginCtx.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		return
	}

	ginCtx.JSON(http.StatusNoContent, nil)
}

func (_testingHandler) getAllLocations(ginCtx *gin.Context) {
	locations, err := data.LocationRepo.GetAll()

	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch locations"})
		return
	}
	ginCtx.JSON(http.StatusOK, locations)
}

func (_testingHandler) GetAllClients(ginCtx *gin.Context) {
	clients, err := data.ClientRepo.GetAll()
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve clients"})
		return
	}
	ginCtx.JSON(http.StatusOK, clients)
}

func (_testingHandler) GetClient(ginCtx *gin.Context) {
	uuid := ginCtx.Param("uuid")
	client, err := data.ClientRepo.Get(uuid)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve client"})
		return
	}
	ginCtx.JSON(http.StatusOK, client)
}

type CreateClientRequest struct {
	Name string `json:"name"`
}

func (_testingHandler) CreateClient(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	client, err := data.ClientRepo.Create(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create client"})
		return
	}
	c.JSON(http.StatusCreated, client)
}

func (_testingHandler) DeleteClient(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := data.ClientRepo.Delete(uuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete client"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
