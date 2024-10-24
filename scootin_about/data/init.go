package data

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	host       = "localhost"
	port       = 5432
	user       = "tom_scoot"
	dbname     = "scootdb"
	maxRetries = 5
	retryDelay = 3 * time.Second
)

var db *sql.DB

func Init() error {
	password, found := os.LookupEnv(`PGPASSWORD`)
	if !found {
		log.Fatal("PGPASSWORD not set")
		return errors.New("PGPASSWORD not set")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("Error connecting to database: %v. Retrying in %v...", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to the database!")
			return nil
		}

		log.Printf("Database is not ready yet. Retrying in %v...", retryDelay)
		time.Sleep(retryDelay)
	}

	log.Fatal("Failed to connect to the database after several attempts")
	return err
}

func Stop() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database: %v\n", err)
		} else {
			log.Println("Database connection closed successfully.")
		}
	}
}
