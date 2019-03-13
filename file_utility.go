package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

/* At start of the connector check if a batch file is present and not empty
if it's the case it means that a preivous load didn't go well when working in mode persistent and we need to reload*/
func checkIfAlreadyPresentAndLoad(filePath string, batch int, gpssClient *gpssClient) {

	log.Printf("Checking if batch file already present")
	// Check if the file already exist and is not empty
	fi, err := os.Stat(filePath)
	if err != nil {
		// file not exist yet, good no previous batch to load
		log.Printf("Batch file not present")
		return
	}

	size := fi.Size()

	if size != 0 {
		log.Printf("Batch file present and size > 0, I need to consume previous batches: Open the file")

		// it was a batch from a previous execution we need to load back the info, open the file
		file, _ := os.Open(filePath)
		buffer := make([]string, batch)
		reader := bufio.NewReader(file)

		log.Printf("Connecting to database to catch the number of fields of the table")
		// Check of how many columns the table is made
		gpssClient.ConnectToGreenplumDatabase()
		lenRow := len(gpssClient.DescribeTable().Columns)
		gpssClient.DisconnectToGreenplumDatabase()

		// Load the file and decompose
		var line string
		i := 0
		log.Printf("Reading file")
		for {

			row := ""
			for j := 0; j < lenRow; j++ {
				line, err = reader.ReadString('\n')
				row = row + line
			}

			if err == nil {
				row = row[:len(row)-1]
				buffer[i] = row
				i++
			}

			if i == batch-1 || err != nil {
				// end of a batch ask gpss to write
				log.Printf("Sending request to gpss... Writing on the database")
				gpssClient.ConnectToGreenplumDatabase()
				gpssClient.WriteToGreenplum(buffer)
				gpssClient.DisconnectToGreenplumDatabase()
				i = 0

			}
			if err != nil {
				// nothing else to read
				break
			}
		}
		// We sent what remained to gp we can now truncate the file
		log.Printf("Batch completed I'm truncated the existing file")
		os.Truncate(filePath, 0)

	}

}

type AppConfigProperties map[string]string

/* Read property file */
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
