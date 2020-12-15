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
	timeoutChannel  = make(chan struct{})
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
		log.Print("Waiting for remote message...")
		n, cliAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		msg := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for msg.Scan() {
			msg := msg.Text()
			log.Print("Received: " + msg + " from " + cliAddr.String() + "\n")
			go func() {messageChannel <- msg}()
		}
		log.Print("Finished parsing remote msg")
	}
}

// ReadUserInput boucle infinie lisant les inputs de l'utilisateur
func ReadUserInput() {

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
			log.Print("User has inputted the election command")
			go func() {electionChannel <- struct{}{}}()
		} else if gCmd.Match(userInput) {
			log.Print("User has inputted the getElu command")
			go func() {getEluChannel <- struct{}{}}()
		} else {
			log.Print("User is trying to set his aptitude")
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
	/*
	 * By default sends and receives block until both the sender and receiver are ready
	 * source: https://gobyexample.com/channels
	 */
	go func() { electionChannel <- struct{}{} }()

	// TODO Election fait une func un timer avec afterfunc plutot que declarer
	//  afterFunc plrs fois

	for {
		log.Print("Waiting for new message...")
		select {
			case <- electionChannel:
				// TODO wait here (bloquant mais où)
				fmt.Print("Lancement d'une nouvelle election")
				bullyImpl.Election()
				time.AfterFunc(electionDuration, func (){ timeoutChannel <- struct{}{} })
				break
			case <- getEluChannel:
				// TODO wait here (bloquant mais où)
				log.Print("L'utilisateur veut connaitre l'elu")
				elu := bullyImpl.GetElu()
				msg := "Le processus " + strconv.Itoa(elu) + " est l'elu!"
				fmt.Println(msg)
				log.Print(msg)
				break
			case msg := <- messageChannel:
				processId, apt := handleMessage(msg)

				if !bullyImpl.EnCours() {
					log.Print("Election lancee apres reception de MESSAGE(pid, apt)")
					bullyImpl.Election()
					time.AfterFunc(electionDuration, func (){ timeoutChannel <- struct{}{} })
				}
				bullyImpl.SetApt(processId, apt)
				break
			case <- timeoutChannel: // timeout
				// "When the Timer expires, the current time will be sent on C"
				// source: https://golang.org/pkg/time/#Timer
				log.Print("Timeout! Fin de l'election")
				bullyImpl.Timeout()

				if initialElection {
					initialElection = false
					fmt.Println("Fin de l'election initial")
					elu := bullyImpl.GetElu()

					fmt.Println("L'elu de l'election initiale est le processus: " +
						strconv.Itoa(elu))
				}

				/*
				 * On peut interactivement changer depuis la console de chaque processus
				 * l’aptitude de celui-ci. Dans ce cas, une nouvelle élection sera déclenchée
				 * pour en tenir compte mais seulement à l’issue d’une éventuelle élection
				 * qui aurait pu démarrer avant ce changement.
				 */
				if bullyImpl.IHaveChangedAptitude() {
					msg := "L'aptitude ayant changee depuis la derniere election precedente, une nouvelle election va etre lancee"
					fmt.Println(msg)
					log.Print(msg)
					go func() { electionChannel <- struct{}{} }()
				}
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


