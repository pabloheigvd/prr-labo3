/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	Communication.go
 */

package Communication

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"os"
	"prr-labo3/Entities"
	"prr-labo3/interface"
	"strconv"
	"strings"
)

const TRANSPORT_PROTOCOL = "udp"

var (
	userIo          = bufio.NewReader(os.Stdin)
	electionChannel = make(chan struct{})
	getEluChannel   = make(chan struct{})
	messageChannel  = make(chan string)
)

/* ==============
 * === Public ===
 * =========== */

// ListenToRemoteMessage from other process
func ListenToRemoteMessage(moiP Entities.Process){
	conn, err := net.ListenPacket(TRANSPORT_PROTOCOL, moiP.GetAdress())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// from udp slide
	buf := make([]byte, 1024)
	for {
		n, cliAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		msg := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for msg.Scan() {
			msg := msg.Text() + " from " + cliAddr.String() + "\n"
			log.Print("Received: " + msg)
			messageChannel <- msg
		}
	}
}

// ReadUserInput to userInput
func ReadUserInput(bully Interface.Bully) {

	// Déclencher une élection au démarrage du processus
	bully.Election()

	eCmd := Entities.ElectionCmd{}
	gCmd := Entities.GetEluCmd{}

	for {
		userInput, err := getUserInput() // bloquant
		if err != nil {
			log.Fatal(err)
		}

		if eCmd.Match(userInput) {
			electionChannel <- struct{}{}
		}

		if gCmd.Match(userInput) {
			getEluChannel <- struct{}{}
		}
		//bully.Election()
		//_ = bully.GetElu()

		// select sur Election et GetElu?

		// FIXME move?
		//fmt.Print(elu)
		//select {
		//	case <- time.After(electionMaxDuration):
		//		log.Println("timeout")
		//		// Algorithm.Timeout()
		//		// election.EnCours = false
		//		break
		//}
	}
}

/* ===============
 * === private ===
 * =============*/

// getUserInput bloquant
func getUserInput() (string, error) {
	return userIo.ReadString('\n')
}

// handleCommunication
func handleCommunication() {
	for {
		select {
			case <- electionChannel:
				log.Print("Lancement d'une nouvelle election")
				break
			case <- getEluChannel:
				log.Print("L'utilisateur veut connaitre l'elu")
				break
			case msg := <-messageChannel:
				_, _ = handleMessage(msg)
				// TODO maj election
				break
			default:
				break
		}
	}
}

// handleMessage aptitude
func handleMessage(msg string) (int, int){
	log.Print("Received: " + msg)
	tokens := strings.Split(msg, " ")
	if len(tokens) < 2 {
		log.Fatal("Invalid message received")
	}

	processId, err := strconv.Atoi(tokens[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Print("received process id: " + strconv.Itoa(processId))
	aptitude, err := strconv.Atoi(tokens[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Print("received aptitude: " + strconv.Itoa(aptitude))

	return processId, aptitude
}


