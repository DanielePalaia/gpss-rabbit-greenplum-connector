package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "test"
)

func connectToGpss() *gpssClient {

	/* Reading properties from ./properties.ini */
	prop, _ := ReadPropertiesFile("./properties.ini")
	port, _ := strconv.Atoi(prop["GreenplumPort"])

	log.Printf("Properties read: Connecting to the Grpc server specified")

	/* Connect to the grpc server specified */
	gpssClient := MakeGpssClient(prop["GpssAddress"], prop["GreenplumAddress"], int32(port), prop["GreenplumUser"], prop["GreenplumPassword"], prop["Database"], prop["SchemaName"], prop["TableName"])
	gpssClient.ConnectToGrpcServer()

	return gpssClient

}

func connectToRabbitAndSend(batch int) (string, string) {
	prop, _ := ReadPropertiesFile("./properties.ini")

	ch := connect(prop["rabbit"])
	send(ch, prop["queue"], batch)

	return prop["rabbit"], prop["queue"]

}

func connectToPostgres(t *testing.T) int {

	prop, _ := ReadPropertiesFile("./properties.ini")
	host := prop["GreenplumAddress"]
	user := prop["GreenplumUser"]
	//passwd := prop["GreenplumPassword"]
	port, _ := strconv.Atoi(prop["GreenplumPort"])
	dbName := prop["Database"]
	tableName := prop["TableName"]
	//schemaName := prop["SchemaName"]

	var db *sql.DB
	var err error

	// Connecting to the database
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbName)

	if db, err = sql.Open("postgres", dbinfo); err != nil {
		t.Errorf("Error connecting to the database")
	}

	count := 0

	if err := db.QueryRow("SELECT count(*) from " + tableName + ";").Scan(&count); err != nil {
		t.Errorf("Error querying the database %v", err)
	}

	return count

}
