package data

import "errors"

type Client struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type _clientRepo struct{}

var ClientRepo = _clientRepo{}

func (_clientRepo) Create(name string) (Client, error) {
	query := `INSERT INTO clients (name) VALUES ($1) RETURNING uuid, name`
	var client Client
	err := db.QueryRow(query, name).Scan(&client.Uuid, &client.Name)
	if err != nil {
		return Client{}, err
	}

	return client, nil
}

func (_clientRepo) Get(uuid string) (Client, error) {
	var client Client

	query := `SELECT uuid, name FROM clients WHERE uuid = $1`
	err := db.QueryRow(query, uuid).Scan(&client.Uuid, &client.Name)
	if err != nil {
		return Client{}, err
	}

	return client, nil
}

func (_clientRepo) GetAll() ([]Client, error) {
	var clients []Client

	query := `SELECT uuid, name FROM clients`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var client Client

		if err := rows.Scan(&client.Uuid, &client.Name); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

func (_clientRepo) Delete(uuid string) error {
	query := `DELETE FROM clients WHERE uuid = $1`
	res, err := db.Exec(query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("client not found")
	}

	return nil
}
