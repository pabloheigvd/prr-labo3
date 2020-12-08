/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	Parsing.go
 */

package Utils

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"prr-labo3/Entities"
)

const CONFIG_FILE_PATH = "./config.json"

var userIo = bufio.NewReader(os.Stdin)

// Parsing of the configuration file
func Parsing() (Entities.Configuration){
	log.Print("Opening configuration file: " + CONFIG_FILE_PATH)
	jsonFile, err := os.Open(CONFIG_FILE_PATH)
	if err != nil {
		log.Print("Error while opening " + CONFIG_FILE_PATH)
		log.Fatal(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var configuration Entities.Configuration
	json.Unmarshal(byteValue, &configuration)
	err = jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Configuration (raw):")
	log.Print(configuration)

	return configuration
}

// GetUserInput bloquant
func GetUserInput() (string, error) {
	return userIo.ReadString('\n')
}

// ParseUserInput clean user input by trimming
func ParseUserInput(userInput string) string {
	return trimBN(userInput)
}

// trimBN enlève '\n'
func trimBN(input string) string {
	return input[:len(input)-1]
}



