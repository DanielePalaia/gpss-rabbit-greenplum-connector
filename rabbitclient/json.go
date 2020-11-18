package main

import (
	"io/ioutil"
)

// GetJSON Returns a json string from file
func GetJSON() string {
	json, _ := ioutil.ReadFile("test.json")
	return string(json)
}
