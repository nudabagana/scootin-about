package simulator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nudabagana/scootin-about/data"
)

const (
	authHeader = "Bearer EXAMPLE"
)

type UserSimulator struct {
	UserID       string
	ScooterID    string
	InTrip       bool
	TripStart    time.Time
	CurrentLat   float64
	CurrentLong  float64
	StopChannel  chan bool
	UpdateTicker *time.Ticker
	router       *gin.Engine
}

func NewUserSimulator(userID string, router *gin.Engine) *UserSimulator {
	return &UserSimulator{
		UserID:       userID,
		CurrentLat:   0,
		CurrentLong:  0,
		StopChannel:  make(chan bool),
		UpdateTicker: time.NewTicker(1 * time.Second),
		router:       router,
	}
}

func (u *UserSimulator) Start() {
	defer u.UpdateTicker.Stop()

	for {
		select {
		case <-u.StopChannel:
			fmt.Printf("Stopping user %s\n", u.UserID)
			return
		case <-u.UpdateTicker.C:
			u.chooseAction()
		}
	}
}

func (u *UserSimulator) Stop() {
	if u.InTrip {
		u.endTrip()
	}
	u.StopChannel <- true
	close(u.StopChannel)
}

func (u *UserSimulator) chooseAction() {
	if !u.InTrip && rand.Float32() < 0.3 {
		u.startTrip()
	} else if u.InTrip && time.Since(u.TripStart) > 10*time.Second && rand.Float32() < 0.5 {
		u.endTrip()
	}

	if u.InTrip && int(time.Since(u.TripStart).Seconds())%3 == 0 {
		u.reportLocation()
	}
}

func (u *UserSimulator) startTrip() {
	req, _ := http.NewRequest("GET", "/client/scooters?top_left_lat=45.47&top_left_long=-76.05&bottom_right_lat=44.95&bottom_right_long=-75.30&occupied=false", nil)
	req.Header.Set("Authorization", authHeader)
	w := httptest.NewRecorder()
	u.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("Failed to retrieve scooters for user %s\n", u.UserID)
		return
	}

	var response struct {
		Data  []data.ScooterWithLocation `json:"data"`
		Count int                        `json:"count"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response.Data) == 0 {
		fmt.Printf("No available scooters for user %s\n", u.UserID)
		return
	}

	u.ScooterID = response.Data[0].Uuid
	u.CurrentLat = response.Data[0].Latitude
	u.CurrentLong = response.Data[0].Longitude

	startTripURL := fmt.Sprintf("/scooter/%s/start-trip", u.ScooterID)
	tripReq, _ := http.NewRequest("PUT", startTripURL, bytes.NewBuffer([]byte{}))
	tripReq.Header.Set("Authorization", authHeader)

	tripW := httptest.NewRecorder()
	u.router.ServeHTTP(tripW, tripReq)

	if tripW.Code == http.StatusOK {
		u.InTrip = true
		u.TripStart = time.Now()
		fmt.Printf("User %s started a trip on scooter %s\n", u.UserID, u.ScooterID)
	} else {
		fmt.Printf("Failed to start trip for user %s on scooter %s\n", u.UserID, u.ScooterID)
	}
}

func (u *UserSimulator) endTrip() {
	if !u.InTrip || u.ScooterID == "" {
		fmt.Printf("No active trip to end for user %s\n", u.UserID)
		return
	}

	endTripURL := fmt.Sprintf("/scooter/%s/end-trip", u.ScooterID)
	req, _ := http.NewRequest("PUT", endTripURL, bytes.NewBuffer([]byte{}))
	req.Header.Set("Authorization", authHeader)

	w := httptest.NewRecorder()
	u.router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		fmt.Printf("User %s ended the trip on scooter %s\n", u.UserID, u.ScooterID)
		u.InTrip = false
		u.ScooterID = ""
	} else {
		fmt.Printf("Failed to end trip for user %s on scooter %s\n", u.UserID, u.ScooterID)
	}
}

func (u *UserSimulator) reportLocation() {
	if !u.InTrip || u.ScooterID == "" {
		fmt.Printf("User %s is not in a trip; no location to report.\n", u.UserID)
		return
	}

	u.CurrentLat += 0.001
	u.CurrentLong += 0.001

	locationData := map[string]float64{
		"latitude":  u.CurrentLat,
		"longitude": u.CurrentLong,
	}
	payload, _ := json.Marshal(locationData)

	reportLocationURL := fmt.Sprintf("/scooter/%s/report-location", u.ScooterID)
	req, _ := http.NewRequest("POST", reportLocationURL, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	u.router.ServeHTTP(w, req)

	if w.Code == http.StatusCreated {
		fmt.Printf("User %s reported location: (%f, %f)\n", u.UserID, u.CurrentLat, u.CurrentLong)
	} else {
		fmt.Printf("Failed to report location for user %s on scooter %s\n", u.UserID, u.ScooterID)
	}
}

var (
	simulatorManager = struct {
		simulators []*UserSimulator
		mu         sync.Mutex
	}{}
)

func SimulateUser(userID string, router *gin.Engine) {
	userSim := NewUserSimulator(userID, router)

	simulatorManager.mu.Lock()
	simulatorManager.simulators = append(simulatorManager.simulators, userSim)
	simulatorManager.mu.Unlock()

	go userSim.Start()
}

func StopAllSimulations() {
	simulatorManager.mu.Lock()
	defer simulatorManager.mu.Unlock()

	for _, sim := range simulatorManager.simulators {
		sim.Stop()
	}

	simulatorManager.simulators = nil
}
