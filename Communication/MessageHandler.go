package Communication

import (
	"log"
	"strconv"
	"strings"
	"time"
)

/* ===============
 * === private ===
 * =============*/

// handleRemoteMessage aptitude
func handleRemoteMessage(msg string){
	log.Print("handling: " + msg)
	tokens := strings.Split(msg, " ")

	msgType := tokens[0]
	switch msgType {
	case "MESSAGE":
		if len(tokens) < 3 {
			log.Fatal("Invalid MESSAGE received")
		}
		handleMessage(tokens[1], tokens[2])
		break
	case "PING":
		if len(tokens) < 2 {
			log.Fatal("Invalid PING received")
		}
		handlePing(tokens[1])
		break
	case "PONG": handlePong()
	default: log.Fatal("Unknown message type: " + msgType)
	}
}

// handleMessage
func handleMessage(pId string, apt string) {
	processId, err := strconv.Atoi(pId)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("received process id: " + strconv.Itoa(processId))

	aptitude, err := strconv.Atoi(apt)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("received aptitude: " + strconv.Itoa(aptitude))

	if !bullyImpl.EnCours() {
		log.Print("Election lancee apres reception de MESSAGE(pid, apt)")
		stopPinging()
		bullyImpl.Election()
		electionDuration := 2*t
		time.AfterFunc(electionDuration, func (){ timeoutChannel <- struct{}{} })
	}
	bullyImpl.SetApt(processId, aptitude)
}