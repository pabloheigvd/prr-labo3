package Communication

import (
	"log"
	"strconv"
	"time"
)

var endTimer = make(chan struct{})

// sendPing to coordinator every interval duration
func sendPing(interval time.Duration){
	pinging := true
	log.Print("Pinging set to " + strconv.FormatBool(pinging))
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

// shouldIPing
func shouldIPing() bool {
	return !initialElection &&
		!bullyImpl.EnCours() &&
		!bullyImpl.IsCoordinator()
}

// stopPinging
func stopPinging() {
	if !initialElection {
		log.Print("No pinging desired")
		go func() { endTimer <- struct{}{} }()
	}
}

// handlePing
func handlePing(pId string) {
	pingingProcessId, err := strconv.Atoi(pId)
	if err != nil {
		log.Fatal(err)
	}
	pingingProcess := bullyImpl.GetProcess(pingingProcessId)

	pingingProcess.Pong()
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
