package main

import (
	"fmt"
	"log"
	"os"
	"prr-labo3/Entities"
	"prr-labo3/Utils"
	"strconv"
	"time"
)

const LOG_DIRECTORY_NAME = "logs"
const LOG_FILE_PREFIX = "log"
const SLOWDOWN_DEBUG_DURATION = 8 * time.Second
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
		errMsg := "You need to start a configuration between 0 and " +
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

	fmt.Print(election)

	// traitement d'une réception à la fois
	// go routine, for select pour ne traiter qu'un seul msg à la fois

	// Prise de connaissance des aptitudes des autres processus


	// demarrage des elections
	/*
	boucle sans fin sur les réceptions de
	- [enCours = faux] élection: demande élection
	- [enCours = faux] getElu // attente fin élection
	- MESSAGE(i,apti) // réception aptitude
	- Timeout // fin d élection
	fin boucle
	 */
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