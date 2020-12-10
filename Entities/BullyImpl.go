/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	BullyImpl.go
 */

package Entities

import (
	"log"
	"strconv"
	"time"
)

type BullyImpl struct {
	election  Election
	processes []Process
}

const SLEEP_CYCLE_DURATION = 50 * time.Millisecond

/* ==============
 * === Public ===
 * ============*/
// rappel: Méthodes définies uniquement sur des types déclarés dans le même package

// InitBully constructeur pour ne pas exposer les champs
func (b *BullyImpl) InitBully (el Election, ps []Process){
	nbProcesses := len(ps)
	nbDefinedAptitude := len(el.Apts)
	if nbProcesses != nbDefinedAptitude {
		log.Print("Nb of processes: " + strconv.Itoa(nbProcesses))
		log.Print("Nb of defined aptitudes: " + strconv.Itoa(nbDefinedAptitude))
		log.Fatal("Please, define the aptitudes of all processes")
	}

	b.election = el
	b.processes = ps
}

// Implementation implicite de l'interface Bully
// Election lancement d'une nouvelle election
func (b *BullyImpl) Election() {
	b.waitUntilElectionIsOver()

	// note: * est essentiel pour que EnCours passe a true
	b.demarre()
}

// GetElu appel bloquant renvoyant l'élu
func (b BullyImpl) GetElu() int {
	b.waitUntilElectionIsOver()
	return b.getElu()
}

// SetMonApt no effect on current election
func (b *BullyImpl) SetMonApt(monApt int) {
	// TODO verification de l'input
	b.election.MonApt = monApt
	log.Print("MonApt = " + strconv.Itoa(monApt))
}

// SetApt used while in election
func (b *BullyImpl) SetApt(processId int, apt int) {
	if processId < 0 || processId >= b.election.N {
		log.Fatal("Process doesn't exist")
	}

	// TODO validation apt

	b.election.Apts[processId] = apt
	b.processes[processId].aptitude = apt
	log.Print("Process " + strconv.Itoa(processId) + " aptitude is " + strconv.Itoa(apt))
}

// EnCours vrai s'il y a une élection en cours
func (b BullyImpl) EnCours() bool {
	return b.election.EnCours
}

// timeout actions à prendre en fin d'une élection
func (b *BullyImpl) Timeout() {
	b.election.Elu = b.getElu()
	b.election.EnCours = false
	log.Print("L'election est terminee")
}

///* ===============
// * === private ===
// * =============*/

func (b *BullyImpl) waitUntilElectionIsOver(){
	// attente active
	cycles := 0
	for b.election.EnCours {
		cycles++
		time.Sleep(SLEEP_CYCLE_DURATION) // ne pas tuer la machine
	}
	waitingTime := time.Duration(cycles) * SLEEP_CYCLE_DURATION
	log.Print("Wait time: " + waitingTime.String())
}

// Message màj de l'aptitude du processus appelant en interne
func (b *BullyImpl) message(processId int, aptitude int){
	if !b.election.EnCours{
		b.demarre()
	}

	b.election.Apts[processId] = aptitude
}

// getElu renvoie l'elu avec la règle de départage indiquée
func (b *BullyImpl) getElu() int {
	coordinator := 0 // <=> elu
	for i, _ := range b.processes {
		log.Print("Parsing process " + strconv.Itoa(i) + " with apt " +
			strconv.Itoa(b.election.Apts[i]))
		if b.election.Apts[i] > b.election.Apts[coordinator] {
			coordinator = i
		} else if b.election.Apts[i] == b.election.Apts[coordinator] {
			/*
			 règle de départage:
			 "En cas d’égalité d’aptitudes, celui ayant
			 le plus petit numéro sera élu"
			*/
			if b.processes[i].No < b.processes[coordinator].No {
				log.Print("Departage!")
				coordinator = i
			}
		}
	}

	log.Print("L'elu est le processus " + strconv.Itoa(coordinator))
	return b.processes[coordinator].No
}

// demarre une election
func (b *BullyImpl) demarre() {
	b.election.EnCours = true
	log.Print("Election en cours")
	for i, _ := range b.election.Apts {
		b.election.Apts[i] = 0
	}

	b.SetApt(b.election.Moi, b.election.MonApt)
	log.Print("Mon apt pour cette election: " +
		strconv.Itoa(b.election.Apts[b.election.Moi]))

	moi := b.election.Moi
	moiP := b.processes[moi]
	moiP.EnvoiMessage(b.processes)

	// timer enclenché dans Communication.go
}
