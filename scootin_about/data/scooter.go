package data

import (
	"database/sql"
	"errors"
)

type Scooter struct {
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	Occupied bool   `json:"occupied"`
}

type _scooterRepo struct{}

var ScooterRepo = _scooterRepo{}

func (_scooterRepo) GetAll() ([]Scooter, error) {
	var scooters []Scooter

	query := `SELECT uuid, name, occupied FROM scooters`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var scooter Scooter

		if err := rows.Scan(&scooter.Uuid, &scooter.Name, &scooter.Occupied); err != nil {
			return nil, err
		}
		scooters = append(scooters, scooter)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scooters, nil
}

func (_scooterRepo) Get(uuid string) (Scooter, error) {
	var scooter Scooter

	query := `SELECT uuid, name, occupied FROM scooters WHERE uuid = $1`

	err := db.QueryRow(query, uuid).Scan(&scooter.Uuid, &scooter.Name, &scooter.Occupied)
	if err != nil {
		if err == sql.ErrNoRows {
			return scooter, errors.New("scooter not found")
		}
		return scooter, err
	}

	return scooter, nil
}

func (_scooterRepo) Delete(uuid string) error {
	query := `DELETE FROM scooters WHERE uuid = $1`

	result, err := db.Exec(query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("scooter not found")
	}

	return nil
}

func (_scooterRepo) CreateScooter(name string, occupied bool) (Scooter, error) {
	var newScooter Scooter

	query := `INSERT INTO scooters (name, occupied) VALUES ($1, $2) RETURNING uuid, name, occupied`

	err := db.QueryRow(query, name, occupied).Scan(&newScooter.Uuid, &newScooter.Name, &newScooter.Occupied)
	if err != nil {
		return newScooter, err
	}

	return newScooter, nil
}

func (_scooterRepo) SetOccupied(uuid string, occupied bool) error {
	query := `UPDATE scooters SET occupied = $1 WHERE uuid = $2`

	result, err := db.Exec(query, occupied, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("scooter not found")
	}

	return nil
}

type ScooterWithLocation struct {
	Scooter
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (_scooterRepo) GetScootersInSquare(topLeftLat, topLeftLong, bottomRightLat, bottomRightLong float64, occupied *bool) ([]ScooterWithLocation, error) {
	var scooters = []ScooterWithLocation{}

	query := `
		SELECT s.uuid, s.name, s.occupied, l.latitude, l.longitude
		FROM scooters s
		JOIN LATERAL (
			SELECT latitude, longitude
			FROM locations
			WHERE scooter_uuid = s.uuid
			ORDER BY created_at DESC
			LIMIT 1
		) l ON TRUE
		WHERE l.latitude BETWEEN $1 AND $2
		AND l.longitude BETWEEN $3 AND $4
	`

	var args []interface{}
	args = append(args, bottomRightLat, topLeftLat, topLeftLong, bottomRightLong)

	if occupied != nil {
		query += " AND s.occupied = $5"
		args = append(args, *occupied)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var scooter ScooterWithLocation

		if err := rows.Scan(&scooter.Uuid, &scooter.Name, &scooter.Occupied, &scooter.Latitude, &scooter.Longitude); err != nil {
			return nil, err
		}

		scooters = append(scooters, scooter)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scooters, nil
}
