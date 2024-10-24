package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nudabagana/scootin-about/data"
	"github.com/stretchr/testify/assert"
)

func TestE2E_Get_Scooters_Success_Case(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()
	Init(router)
	data.Init()
	authHeader := "Bearer EXAMPLE"

	scooters := []struct {
		name     string
		occupied bool
	}{
		{"Scooter A", false},
		{"Scooter B", true},
	}

	var createdScooters []struct {
		Uuid string `json:"uuid"`
	}

	for _, scooter := range scooters {
		payload, _ := json.Marshal(CreateScooterRequest{
			Name:     scooter.name,
			Occupied: scooter.occupied,
		})
		req, _ := http.NewRequest("POST", "/testing/scooters", bytes.NewBuffer(payload))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Uuid string `json:"uuid"`
		}
		json.Unmarshal(w.Body.Bytes(), &response)
		println(response.Uuid)
		createdScooters = append(createdScooters, response)
	}

	for _, scooter := range createdScooters {
		locationPayload, _ := json.Marshal(ReportLocationRequest{
			Latitude:  25.4215,
			Longitude: -75.6972,
		})
		req, _ := http.NewRequest("POST", fmt.Sprintf("/scooter/%s/report-location", scooter.Uuid), bytes.NewBuffer(locationPayload))
		req.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	req, _ := http.NewRequest("GET", "/client/scooters?top_left_lat=25.5&top_left_long=-75.7&bottom_right_lat=25.4&bottom_right_long=-75.6&occupied=false", nil)
	req.Header.Set("Authorization", authHeader)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Step 4: Validate the response
	var response struct {
		Data  []data.ScooterWithLocation `json:"data"`
		Count int                        `json:"count"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Validate that the count matches the expected number of scooters
	assert.Equal(t, 1, response.Count) // Check we only got one scooter back (Scooter A)
	assert.Equal(t, createdScooters[0].Uuid, response.Data[0].Uuid)

	// Step 5: Clean up (Delete scooters)
	for _, scooter := range createdScooters {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/testing/scooters/%s", scooter.Uuid), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	}
	data.Stop()
}
