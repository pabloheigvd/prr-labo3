/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	Parsing.go
 */

package Utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"prr-labo3/Entities"
)

const CONFIG_FILE_PATH = "./config.json"

func Parsing() (Entities.Configuration, error){
	log.Print("Opening configuration file: " + CONFIG_FILE_PATH)
	jsonFile, err := os.Open(CONFIG_FILE_PATH)
	if err != nil {
		log.Print("Error while opening " + CONFIG_FILE_PATH)
		log.Fatal(err)
	}

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var processes Entities.Configuration
	json.Unmarshal(byteValue, &processes)
	err = jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Unmarshalled json:")
	log.Print(processes)

	return processes, nil
}

func ParseUserInput(userInput string) {
	userInput = trimBN(userInput)
	switch userInput {
		case "e":
		break
	}
}

// trimBN enlève '\n'
func trimBN(input string) string {
	return input[:len(input)-1]
}



