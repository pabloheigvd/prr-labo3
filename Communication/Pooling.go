package Communication

import (
	"log"
	"strconv"
	"time"
)

var endTimer = make(chan struct{})

/* ==============
 * === Public ===
 *=============*/

// StartPing to coordinator every interval duration
func StartPing(interval time.Duration){
	log.Println("Ping Time")
	pinging := true
	for pinging {
		pinging = false
		log.Print("Pinging set to " + strconv.FormatBool(pinging))

		if shouldIPing() {
			pinging = true
			log.Print("Pinging set to " + strconv.FormatBool(pinging))

			bullyImpl.SetIsCoordinatorAlive(false)
			log.Print(bullyImpl.IsCoordinatorAlive())
			bullyImpl.GetMoi().Ping(bullyImpl.GetCoordinator())

			timer := time.NewTimer(interval)
			select {
			case <- endTimer:
				log.Print("Pinging was ended")
			case <- timer.C:
				if shouldIPing() && !bullyImpl.IsCoordinatorAlive() {
					pinging = false
					log.Print("Pinging set to " + strconv.FormatBool(pinging))
					log.Print("Coordinator is not alive, launching a new Election")
					go func(){ electionChannel <- struct{}{} }()
				}
			}
		}
	}
}

// StopPinging
func StopPinging() {
	if !initialElection {
		log.Print("No pinging desired")
		go func() { endTimer <- struct{}{} }()
	}
}

/* ===============
 * === private ===
 * =============*/

// shouldIPing
func shouldIPing() bool {
	return !initialElection &&
		!bullyImpl.EnCours() &&
		!bullyImpl.IsCoordinator()
}

// handlePing
func handlePing(pId string) {
	pingingProcessId, err := strconv.Atoi(pId)
	if err != nil {
		log.Fatal(err)
	}
	pingingProcess := bullyImpl.GetProcess(pingingProcessId)

	pingingProcess.Pong() // assumes this process is the coordinator
}

// handlePong
func handlePong() {
	if !bullyImpl.IsCoordinatorAlive() {
		log.Print("Coordinator responded in time")
		bullyImpl.SetIsCoordinatorAlive(true)
	} else {
		log.Print("Message was received too late")
	}
}
