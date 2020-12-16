package Communication

import (
	"log"
	"strconv"
	"time"
)

// handlePing
func handlePing(pId string) {
	processId, err := strconv.Atoi(pId)
	if err != nil {
		log.Fatal(err)
	}
	process := bullyImpl.GetProcess(processId)

	process.Pong()
}

// handlePong
func handlePong() {
	if !bullyImpl.IsCoordinatorAlive() {
		log.Print("Coordinator responded in time")
		bullyImpl.SetIsCoordinatorAlive(true)
	}
}

// sendPing to coordinator every interval duration
func sendPing(interval time.Duration){
	// TODO fix pong not received properly

	for {
		if !initialElection && !bullyImpl.EnCours() && !bullyImpl.IsCoordinator() {
			bullyImpl.SetIsCoordinatorAlive(false)
			bullyImpl.GetMoi().Ping(bullyImpl.GetCoordinator())
		}

		timer := time.NewTimer(interval)
		<- timer.C
		if !initialElection &&
			!bullyImpl.EnCours() &&
			!bullyImpl.IsCoordinator() &&
			!bullyImpl.IsCoordinatorAlive() {
			log.Print("Coordinator is not alive, launching a new Election")
			go func(){ electionChannel <- struct{}{} }()
		}
	}
}