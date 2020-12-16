/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	Process.go
 */

package Entities

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type Process struct {
	No				int 	`json:"no"`
	InitialAptitude	int 	`json:"initialAptitude"`
	Adr   			string 	`json:"adr"`
	NoPort   		int 	`json:"noPort"`
	aptitude		int
}

// GetAdress returns addr:port of node
func (p Process) GetAdress() string {
	return p.Adr + ":" + strconv.Itoa(p.NoPort)
}

// EnvoiMessage à tous sauf à moi
func (p Process) EnvoiMessage(processes []Process) {
	moi := p.No
	msg := p.getMessage()

	for i, p := range processes {
		if i == moi { continue } // à tous sauf moi
		p.sendMsg(msg)
		// hypothese: le réseau est fiable mais même localhost a de la peine avec
		// la congestion
		time.Sleep(15 * time.Millisecond)
	}
}

// Pong envoyer un message en réponse à un ping
func (p Process) Ping(coordinator Process){
	coordinator.sendMsg("PING " + strconv.Itoa(p.No))
}

// Pong envoyer un message en réponse à un ping
func (p Process) Pong(){
	p.sendMsg("PONG")
}

/* ===============
 * === private ===
 * =============*/

// getMessage a envoyer aux autres processus
func (p Process) getMessage() string {
	return "MESSAGE " + strconv.Itoa(p.No) + " " + strconv.Itoa(p.aptitude)
}

// sendMsg envoie le message msg au processus p
func (p Process) sendMsg(msg string) {
	addr := p.GetAdress()
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprint(conn, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Envoye a " + addr + ":")
	log.Print(msg)
}