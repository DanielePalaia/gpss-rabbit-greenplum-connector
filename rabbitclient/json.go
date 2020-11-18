package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

// GetJSON Returns a json string from file
func GetJSON() string {
	json, _ := ioutil.ReadFile("test.json")
	return string(json)
}

// GetJSONSingleLine Returns a json string from file
func GetJSONSingleLine() string {

	json := ""
	file, _ := os.Open("test.json")

	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		readline := s.Text()
		json += readline

		// other code what work with parsed line...
	}
	fmt.Printf("json: %s", json)

	return json
}
