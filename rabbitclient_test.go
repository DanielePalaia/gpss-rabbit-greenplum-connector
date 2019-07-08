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

// Functional test of the product
func TestRabbit(t *testing.T) {

	// Read rabbit and gpss info
	gpssClient := connectToGpss()
	log.Printf("Connected to the grpc server")
	prop, _ := ReadPropertiesFile("./properties.ini")

	batch, _ := strconv.Atoi(prop["batch"])

	// Connecting to rabbit and insert 3 elements
	rabbitserver, queueName := connectToRabbitAndSend(batch)

	// Call the main process
	rabbit := makeRabbitClient(rabbitserver, queueName, batch, gpssClient, 0, "")
	rabbit.connect()
	rabbit.consume(false)

	// Testing the result
	count := connectToPostgres(t)

	//Check if the batch elements have been inserted
	if count != batch {
		t.Errorf("handler returned wrong status code: got %v want %v",
			count, batch)
	}

}
