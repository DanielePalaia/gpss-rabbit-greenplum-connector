package main

import (
	"log"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
)

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
