package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	port = ":50051"
)

func main() {
	// Get service configuration from FS
	config, configErr := GetServiceConfig(os.Getenv("BLOCK_CONFIG_FILE"))
	if configErr != nil {
		panic(configErr)
	}

	// Handle service discovery
	service, err := HandleInit(config)
	if err != nil {
		panic(err)
	}

	// Setup connection to Postgres database
	connStr := fmt.Sprintf("dbname=%s host=%s port=%d password=%s sslmode=disable user=%s", config.Postgres.DBName, config.Postgres.Host, config.Postgres.Port, config.Postgres.Password, config.Postgres.User)
	db, connErr := sql.Open("postgres", connStr)
	if connErr != nil {
		panic(connErr)
	}

	r := mux.NewRouter()
	r.Path("/pv").Queries("device", "{device}").HandlerFunc(GetPvsByDeviceHandler(db, service)).Methods("GET")

	// Log successful listen
	log.Printf("Started block service with identifier: %s", service.ID)

	// Logs the error if ListenAndServe fails.
	log.Fatal(http.ListenAndServe(port, r))
}
