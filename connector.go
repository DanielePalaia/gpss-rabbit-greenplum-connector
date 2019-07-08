package main

/* Product main file*/
import (
	"log"
	"strconv"
)

func main() {

	log.Printf("Starting the connector and reading properties in the properties.ini file")
	/* Reading properties from ./properties.ini */
	prop, _ := ReadPropertiesFile("./properties.ini")
	port, _ := strconv.Atoi(prop["GreenplumPort"])
	batch, _ := strconv.Atoi(prop["batch"])
	// 1 if batch is persistent, 0 otherwise
	mode, _ := strconv.Atoi(prop["mode"])
	// where entries are temporary stored in batch (if persistence activated)
	batchFile := "./batch"

	log.Printf("Properties read: Connecting to the Grpc server specified")

	/* Connect to the grpc server specified */
	gpssClient := MakeGpssClient(prop["GpssAddress"], prop["GreenplumAddress"], int32(port), prop["GreenplumUser"], prop["GreenplumPassword"], prop["Database"], prop["SchemaName"], prop["TableName"])
	gpssClient.ConnectToGrpcServer()

	log.Printf("Connected to the grpc server")

	// If persistence is set then we need to verify if a past crash caused a batch file to remain filled
	if mode == 1 {
		log.Printf("Mode persistency activated, checking if batch file already exist")
		checkIfAlreadyPresentAndLoad("./batch", batch, gpssClient)
	}

	log.Printf("Connecting to rabbit and consuming")
	/* Generate teh rabbit connection */
	rabbit := makeRabbitClient(prop["rabbit"], prop["queue"], batch, gpssClient, mode, batchFile)
	rabbit.connect()
	rabbit.consume(true)

}
