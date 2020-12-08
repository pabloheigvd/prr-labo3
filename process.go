/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	process.go
 */

package main

import (
	"fmt"
	"log"
	"os"
	"prr-labo3/Communication"
	"prr-labo3/Entities"
	"prr-labo3/Utils"
	"strconv"
	"time"
)

const LOG_DIRECTORY_NAME = "logs"
const LOG_FILE_PREFIX = "log"
const SLOWDOWN_DEBUG_DURATION = 5 * time.Second
const EXIT_MSG = "exiting..."
const TRANSMISSION_MAX_DURATION = 2 * time.Second

func main(){
	if len(os.Args) < 2 {
		log.Print("Please indicate the process id")
		os.Exit(0)
	}

	// Chaque processus devra pouvoir être démarré avec son numéro en paramètre
	var processId, err = strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	log.Print(processId)
	var configuration, _ = Utils.Parsing()

	if !configuration.Trace {
		/**
		 * Redirecting log output to log folder
		 * note: don't try to extract method from it, didn't work for some reason
		 * note: if directory with os.ModeDir permission is created already
		 *       you cannot cd into the dir
		 */
		_ = os.Mkdir(LOG_DIRECTORY_NAME, 0700)
		err := os.Chdir(LOG_DIRECTORY_NAME)
		if err != nil {
			// note: if directory with os.ModeDir permission is created already
			// you cannot cd into the dir
			log.Fatal("Couldn't cd into logs")
		}

		f, err := os.OpenFile(LOG_FILE_PREFIX+ strconv.Itoa(processId),
			os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		err = os.Chdir("..")
		if err != nil {
			log.Fatal("could not cd out of " + LOG_DIRECTORY_NAME)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	/**
	 * Verifying configuration argument
	 */
	if configuration.NbProcess <= processId {
		errMsg := "Error: process id is not between 0 and " +
			strconv.Itoa(configuration.NbProcess)
		fmt.Println(errMsg)
		fmt.Println(EXIT_MSG)
		log.Println(errMsg)
		log.Println(EXIT_MSG)
		os.Exit(0)
	}

	apts := make([]int, configuration.NbProcess) // array declaration
	for i, process := range configuration.Processes {
		apts[i] = process.InitialAptitude
	}

	election := Entities.Election{
		N: configuration.NbProcess,
		Moi: processId,
		MonApt: configuration.Processes[processId].InitialAptitude,
		Apts: apts,
		T: TRANSMISSION_MAX_DURATION,
		EnCours: false,
		Elu: -1, // no one is elected at first
	}

	log.Print("Etat initial:")
	log.Print(election)

	moiP := election.GetProcess(configuration.Processes)

	// traitement d'une réception à la fois
	// go routine, for select pour ne traiter qu'un seul msg à la fois
	go Communication.ListenToRemoteMessage(moiP)

	// demarrage des elections
	// Lorsqu’il démarre, il ne participe pas à une éventuelle élection
	// qui aurait pu débuter avant son démarrage. Il lance ensuite une élection.
	/*
	boucle sans fin sur les réceptions de
	- [enCours = faux] élection: demande élection
	- [enCours = faux] getElu // attente fin élection
	- MESSAGE(i,apti) // réception aptitude
	- timeout // fin d élection
	fin boucle
	 */
	// timeouts: https://gobyexample.com/timeouts
	// Une élection dure au maximum 3T.
	//electionMaxDuration := 3 * election.T
	//Collecte les aptitudes des processus pendant une durée de 2T
	//aptitudeCollectionMaxDuration := 2 * election.T

	bully := Entities.BullyImpl{}
	bully.InitBully(election, configuration.Processes)

	Communication.ReadUserInput(bully)
}

/* FIXME remove
for {
		fmt.Print("Please intialize my aptitude: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		log.Print("read: " + input)
		if err != nil {
			log.Print(err)
			continue
		} else {
			input = input[:len(input) - 1] // shave '\n'
			election.MonApt, err = strconv.Atoi(input)
			log.Print(election)
			if err != nil {
				log.Print(err)
				continue
			}
			break
		}
	}
 */