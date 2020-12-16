/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	BullyChannelLoop.go
 */

package Communication

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const TRANSPORT_PROTOCOL = "udp"

var (
	bullyImpl       BullyImpl
	t               time.Duration
	electionChannel = make(chan struct{})
	getEluChannel   = make(chan struct{})
	messageChannel  = make(chan string)
	timeoutChannel  = make(chan struct{})
	initialElection = true
)

/* ==============
 * === Public ===
 * =========== */

// Init
func Init(bi BullyImpl, tt time.Duration){
	bullyImpl = bi
	t = tt
}

// HandleCommunication boucle infinie modifiant les données critiques de l'élection
func HandleCommunication() {
	electionDuration := 2*t
	pingResponseMaxDelay := electionDuration
	log.Print("Election results are known after " + electionDuration.String())

	bullyImpl.RestrictParticipation()
	/*
	 * By default sends and receives block until both the sender and receiver are ready
	 * source: https://gobyexample.com/channels
	 */
	go func() { electionChannel <- struct{}{} }()

	for {
		select {
			case <- electionChannel:
				log.Print("election channel received a msg")
				bullyImpl.Election()
			case <- getEluChannel:
				log.Print("getElu channel received a msg")
				go func() {
					// note: ne doit pas empêcher un timeout d'arriver
					elu := bullyImpl.GetElu()
					msg := "Le processus " + strconv.Itoa(elu) + " est l'elu!"
					fmt.Println(msg)
					log.Print(msg)
				}()
			case msg := <- messageChannel:
				log.Print("msg channel received a msg")
				if bullyImpl.IsParticipating() {
					handleRemoteMessage(msg)
				} else {
					log.Print("Message was rejected because process does not participate in the election")
				}
			case <- timeoutChannel: // timeout
				log.Print("timeout channel received a msg")
				// "When the Timer expires, the current time will be sent on C"
				// source: https://golang.org/pkg/time/#Timer
				log.Print("Timeout! Fin de l'election")
				bullyImpl.Timeout()

				if initialElection {
					initialElection = false
					fmt.Println("Fin de l'election initial")
					elu := bullyImpl.GetCoordinator().No
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
					msg := "L'aptitude ayant changee depuis la derniere election precedente, " +
						"une nouvelle election va etre lancee"
					fmt.Println(msg)
					log.Print(msg)
					go func() { electionChannel <- struct{}{} }()
				} else {
					go StartPing(pingResponseMaxDelay)
				}
		}
	}
}

// SetTimeout
func SetTimeout() {
	time.AfterFunc(2*t, func (){ timeoutChannel <- struct{}{} })
}





