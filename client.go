package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	prop, _ := ReadPropertiesFile("/Users/dpalaia/GO/src/gpssclient/properties.ini")
	port, _ := strconv.Atoi(prop["GreenplumPort"])
	gpssClient := MakeGpssClient(prop["GpssAddress"], prop["GreenplumAddress"], int32(port), prop["GreenplumUser"], "", prop["Database"], prop["SchemaName"], prop["TableName"])
	gpssClient.ConnectToGrpcServer()
	//gpssClient.ConnectToGreenplumDatabase()
	//gpssClient.WriteToGreenplum(prop["columnfile"], prop["datafile"])
	//gpssClient.CloseRequest()
	batch, _ := strconv.Atoi(prop["batch"])
	rabbit := makeRabbitClient(prop["rabbit"], prop["queue"], batch, gpssClient)
	rabbit.connect()
	rabbit.consume()

}

type AppConfigProperties map[string]string

func ReadPropertiesFile(filename string) (AppConfigProperties, error) {
	config := AppConfigProperties{}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return config, nil
}
