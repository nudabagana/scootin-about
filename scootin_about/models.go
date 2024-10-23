package main

import "time"

type Client struct {
	uuid string
	name string
}

type Scooter struct {
	uuid     string
	name     string
	occupied bool
	active   bool
}

type Location struct {
	id int
	latitude  float64
	longitude float64
	createdAt time.Time
	scooter Scooter
}

type ApiKey struct {
	uuid       string
	key        string
	validUntil time.Time
	clinet     *Client
	scooter    *Scooter
}

