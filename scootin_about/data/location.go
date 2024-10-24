package data

import (
	"database/sql"
	"errors"
	"time"
)

type Location struct {
	Id          int       `json:"id"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	CreatedAt   time.Time `json:"created_at"`
	ScooterUuid string    `json:"scooter_uuid"`
}

type _locationRepo struct{}

var LocationRepo = _locationRepo{}

func (_locationRepo) GetAll() ([]Location, error) {
	var locations []Location

	query := `
		SELECT id, latitude, longitude, created_at, scooter_uuid
		FROM locations
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var location Location

		if err := rows.Scan(&location.Id, &location.Latitude, &location.Longitude, &location.CreatedAt, &location.ScooterUuid); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}

func (_locationRepo) Add(latitude float64, longitude float64, scooterUuid string) (Location, error) {
	var location Location

	query := `
		INSERT INTO locations (latitude, longitude, scooter_uuid)
		VALUES ($1, $2, $3)
		RETURNING id, latitude, longitude, created_at, scooter_uuid
	`

	err := db.QueryRow(query, latitude, longitude, scooterUuid).Scan(&location.Id, &location.Latitude, &location.Longitude, &location.CreatedAt, &location.ScooterUuid)
	if err != nil {
		return location, err
	}

	return location, nil
}

func (_locationRepo) GetByScooterUuid(scooterUuid string) (Location, error) {
	var location Location

	query := `
		SELECT id, latitude, longitude, created_at, scooter_uuid
		FROM locations
		WHERE scooter_uuid = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := db.QueryRow(query, scooterUuid).Scan(&location.Id, &location.Latitude, &location.Longitude, &location.CreatedAt, &location.ScooterUuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return location, errors.New("no location found for the specified scooter UUID")
		}
		return location, err
	}

	return location, nil
}
