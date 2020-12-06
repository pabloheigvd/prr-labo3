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
func (b BullyImpl) Election() {
	b.demarre()
}

// GetElu appel bloquant renvoyant l'élu
func (b *BullyImpl) GetElu() Process {
	// attente active
	cycles := 0
	sleepCycleDuration := 50 * time.Millisecond
	for b.election.EnCours {
		cycles++
		time.Sleep(sleepCycleDuration) // ne pas tuer la machine
	}
	waitingTime := time.Duration(cycles) * sleepCycleDuration
	log.Print("L'elu est connu apres " + waitingTime.String())
	return b.getCoordinator()
}

///* ===============
// * === private ===
// * =============*/

// Message màj de l'aptitude du processus appelant en interne
func (b *BullyImpl) message(processId int, aptitude int){
	if !b.election.EnCours{
		b.demarre()
	}

	b.election.Apts[processId] = aptitude
}

// timeout actions à prendre en fin d'une élection
func (b BullyImpl) timeout() {
	b.GetElu()
	b.election.EnCours = false
}

// getCoordinator renvoie l'elu avec la règle de départage indiquée
func (b *BullyImpl) getCoordinator()  Process {
	coordinator := 0 // <=> elu
	for i, _ := range b.processes {
		if b.election.Apts[i] > b.election.Apts[coordinator] {
			coordinator = i
		} else if b.election.Apts[i] == b.election.Apts[coordinator] {
			/*
			 règle de départage:
			 "En cas d’égalité d’aptitudes, celui ayant
			 le plus petit numéro sera élu"
			*/
			if b.processes[i].No < b.processes[coordinator].No {
				coordinator = i
			}
		}
	}

	log.Print("L'élu est le processus " + strconv.Itoa(coordinator))
	return b.processes[coordinator]
}

// demarre une election
func (b *BullyImpl) demarre() {
	log.Print("Election demarree")
	b.election.EnCours = true
	for i, _ := range b.election.Apts {
		b.election.Apts[i] = 0
	}
	b.election.Apts[b.election.Moi] = b.election.MonApt

	moi := b.election.Moi
	moiP := b.processes[moi]
	moiP.EnvoiMessage(b.processes)

	// TODO enclenche timer???
}
