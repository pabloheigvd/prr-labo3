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
const TRANSMISSION_MAX_DURATION = 1500 * time.Millisecond

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
	var configuration = Utils.Parsing()

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

	// FIXME remove test

	//b := Entities.BullyImpl{}
	//b.InitBully(election, configuration.Processes)
	//
	//fmt.Print(b.EnCours())
	//b.Election()
	//fmt.Print(b.EnCours())
	//time.Sleep(time.Second * 5)
	//b.Timeout()
	//fmt.Print(b.EnCours())

	// -----------

	bully := Entities.BullyImpl{}
	bully.InitBully(election, configuration.Processes)
	Communication.Init(bully, TRANSMISSION_MAX_DURATION)

	// traitement d'une réception à la fois
	// go routine, for select pour ne traiter qu'un seul msg à la fois
	go Communication.ListenToRemoteMessage(moiP)

	go Communication.ReadUserInput()

	Communication.HandleCommunication()
}