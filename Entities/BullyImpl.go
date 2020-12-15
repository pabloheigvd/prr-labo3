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
	// note: * est essentiel pour que EnCours passe a true
	b.Demarre()
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
	log.Print("Process " + strconv.Itoa(processId) + " aptitude was set to " + strconv.Itoa(apt))
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

// Demarre une election
func (b *BullyImpl) Demarre() {
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

///* ===============
// * === private ===
// * =============*/

func (b *BullyImpl) waitUntilElectionIsOver(){
	// attente active
	cycles := 0
	for b.election.EnCours {
		cycles++
		time.Sleep(SLEEP_CYCLE_DURATION) // ne pas tuer la machine
		log.Print("wait cycle " + strconv.Itoa(cycles))
	}
	waitingTime := time.Duration(cycles) * SLEEP_CYCLE_DURATION
	log.Print("Wait time: " + waitingTime.String())
}

// Message màj de l'aptitude du processus appelant en interne
func (b *BullyImpl) message(processId int, aptitude int){
	if !b.election.EnCours{
		b.Demarre()
	}

	b.election.Apts[processId] = aptitude
}

// getElu renvoie l'elu avec la règle de départage indiquée
func (b *BullyImpl) getElu() int {
	missingProcess := []int{}

	coordinator := 0 // <=> elu
	for i, p := range b.processes {
		apt := b.election.Apts[i]
		log.Print("Parsing process " + strconv.Itoa(i) + " with apt " +
			strconv.Itoa(apt))

		// no process have 0 aptitude
		if apt == 0 {
			missingProcess = append(missingProcess, p.No)
		}

		if b.election.Apts[i] > b.election.Apts[coordinator] {
			coordinator = i
		}
		/*
		 * règle de départage:
		 * "En cas d’égalité d’aptitudes, celui ayant
		 * le plus petit numéro sera élu"
		 * note: satifsfaite implicitement dans l'ordre de traitement des
		 * processus
		 */
	}

	if len(missingProcess) == 0 {
		log.Print("No missed the election")
	} else {
		log.Print("Processes who did not participate in the election:")
		for _, pId := range missingProcess {
			log.Print("- process " + strconv.Itoa(pId))
		}
	}

	log.Print("L'elu est le processus " + strconv.Itoa(coordinator))
	return b.processes[coordinator].No
}
