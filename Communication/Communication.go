/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	Communication.go
 */

package Communication

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"prr-labo3/Entities"
	"prr-labo3/Utils"
	"strconv"
	"strings"
	"time"
)

const TRANSPORT_PROTOCOL = "udp"
const STARTUP_DURATION = 2 * time.Second

var (
	bullyImpl		Entities.BullyImpl
	t				time.Duration
	electionChannel = make(chan struct{})
	getEluChannel   = make(chan struct{})
	messageChannel  = make(chan string)
)

/* ==============
 * === Public ===
 * =========== */

// Init
func Init(bi Entities.BullyImpl, tt time.Duration){
	bullyImpl = bi
	t = tt
}

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
			msg := msg.Text()
			log.Print("Received: " + msg + " from " + cliAddr.String() + "\n")
			messageChannel <- msg
		}
	}
}

// ReadUserInput boucle infinie lisant les inputs de l'utilisateur
func ReadUserInput() {

	// Déclencher une élection au démarrage du processus
	electionChannel <- struct{}{}

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
			electionChannel <- struct{}{}
		} else if gCmd.Match(userInput) {
			getEluChannel <- struct{}{}
		} else {
			// aptitude
			monApt, err := strconv.Atoi(userInput)
			if err != nil {
				log.Fatal(err)
			}
			bullyImpl.SetMonApt(monApt)
		}
	}
}

// HandleCommunication boucle infinie modifiant les données critiques de l'élection
func HandleCommunication() {

	waitOtherProcessInitialisation()

	electionDuration := 2*t
	log.Print("Election results are known after " + electionDuration.String())

	initialElection := true

	timer := time.NewTimer(electionDuration)
	log.Print("Election initial demaree")

	for {
		select {
		case <- electionChannel:
			fmt.Print("Lancement d'une nouvelle election")
			bullyImpl.Election()
			timer.Stop()
			timer = time.NewTimer(electionDuration)
			break
		case <- getEluChannel:
			log.Print("L'utilisateur veut connaitre l'elu")
			elu := bullyImpl.GetElu()
			msg := "Le processus " + strconv.Itoa(elu) + " est l'elu!"
			fmt.Println(msg)
			log.Print(msg)
			break
		case msg := <-messageChannel:
			processId, apt := handleMessage(msg)

			if !bullyImpl.EnCours() {
				log.Print("Election lancee apres reception de MESSAGE(pid, apt)")
				bullyImpl.Election()
			}
			bullyImpl.SetApt(processId, apt)
			break
			case <- timer.C: // timeout
			log.Print("Timeout! Fin de l'election")
			bullyImpl.Timeout()

			if initialElection {
				initialElection = false
				fmt.Println("Fin de l'election initial")
				elu := bullyImpl.GetElu()

				fmt.Println("L'elu de l'election initiale est le processus: " +
					strconv.Itoa(elu))
			}
			break
		default:
			break
		}
	}
}

/* ===============
 * === private ===
 * =============*/

// waitOtherProcessInitialisation for e2e tests or manual testing
func waitOtherProcessInitialisation() {
	time.Sleep(STARTUP_DURATION)

	// note: this is useful when you want the initial election to be meaningful
}

// handleMessage aptitude
func handleMessage(msg string) (int, int){
	log.Print("Received: " + msg)
	tokens := strings.Split(msg, " ")
	if len(tokens) < 3 {
		log.Fatal("Invalid message received")
	}

	msgType := tokens[0]
	if msgType != "MESSAGE"{
		log.Fatal("Unknown message type: " + msgType)
	}

	processId, err := strconv.Atoi(tokens[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Print("received process id: " + strconv.Itoa(processId))
	aptitude, err := strconv.Atoi(tokens[2])
	if err != nil {
		log.Fatal(err)
	}
	log.Print("received aptitude: " + strconv.Itoa(aptitude))

	return processId, aptitude
}


