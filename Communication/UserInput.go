package Communication

import (
	"log"
	"prr-labo3/Entities"
	"prr-labo3/Utils"
	"strconv"
)

// ReadUserInput boucle infinie lisant les inputs de l'utilisateur
func ReadUserInput() {
	eCmd := Entities.ElectionCmd{}
	gCmd := Entities.GetEluCmd{}

	for {
		userInput, err := Utils.GetUserInput() // bloquant
		if err != nil {
			log.Fatal(err)
		}

		userInput = Utils.ParseUserInput(userInput)
		log.Print()

		if eCmd.Match(userInput) {
			log.Print("User has inputted the election command")
			go func() {electionChannel <- struct{}{}}()
		} else if gCmd.Match(userInput) {
			log.Print("User has inputted the getElu command")
			go func() {getEluChannel <- struct{}{}}()
		} else {
			log.Print("User is setting up his aptitude")
			monApt, err := strconv.Atoi(userInput)
			if err != nil {
				log.Fatal(err)
			}
			bullyImpl.SetMonApt(monApt)

			if !bullyImpl.EnCours() {
				log.Print("No election was running, let's start one")
				go func() {electionChannel <- struct{}{}}()
			} else {
				log.Print("An election is currently being held. " +
					"An election will be started (if necessary) by the end of the current election")
			}
			/*
			 * note: si l'élection est déjà en cours, alors à la fin de l'élection, on vérifie
			 * s'il est judicieux de lancer une élection en se basant sur la valeur actuelle de
			 * monApt
			 */
		}
	}
}