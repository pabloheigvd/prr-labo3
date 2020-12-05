package main

import (
	"fmt"
	"log"
	"os"
	"prr-labo3/Utils"
	_ "prr-labo3/Utils"
	"strconv"
	"time"
)

const LOG_DIRECTORY_NAME = "logs"
const LOG_FILE_PREFIX = "log"
const SLOWDOWN_DEBUG_DURATION = 8 * time.Second
const EXIT_MSG = "exiting..."

func main(){
	if len(os.Args) < 2 {
		log.Print("Please indicate the process id")
		os.Exit(0)
	}

	var processId, _ = strconv.Atoi(os.Args[1])

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

}