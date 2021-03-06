/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	BullyImpl.go
 */

package Communication

import (
	"fmt"
	"log"
	"prr-labo3/Entities"
	"strconv"
	"time"
)

type BullyImpl struct {
	election           Entities.Election
	processes          []Entities.Process
	isCoordinatorAlive bool
	participating      bool
}

const SLEEP_CYCLE_DURATION = 50 * time.Millisecond

/* ==============
 * === Public ===
 * ============*/
// rappel: Méthodes définies uniquement sur des types déclarés dans le même package

// InitBully constructeur pour ne pas exposer les champs
func (b *BullyImpl) InitBully (el Entities.Election, ps []Entities.Process){
	nbProcesses := len(ps)
	nbDefinedAptitude := len(el.Apts)
	if nbProcesses != nbDefinedAptitude {
		log.Print("Nb of processes: " + strconv.Itoa(nbProcesses))
		log.Print("Nb of defined aptitudes: " + strconv.Itoa(nbDefinedAptitude))
		log.Fatal("Please, define the aptitudes of all processes")
	}

	b.election = el
	b.processes = ps
	b.participating = false
}

// Implementation implicite de l'interface Bully
// Election lancement d'une nouvelle election
func (b *BullyImpl) Election() {
	// note: * est essentiel pour que EnCours passe a true
	StopPinging()
	// Note: ne pas traîter les messages du canal de messages et les consommer
	// ensuite
	b.WaitUntilElectionIsOver()

	fmt.Println("Lancement d'une nouvelle election")
	b.Demarre()
}

// IHaveChangedAptitude = vrai si l'aptitude actuelle est différente de
// celle utilisée lors de la dernière élection
func (b BullyImpl) IHaveChangedAptitude() bool {
	return b.election.MonApt != b.election.Apts[b.election.Moi]
}

// SetMonApt no effect on current election
func (b *BullyImpl) SetMonApt(monApt int) {
	if monApt <= 0 {
		log.Fatal("Aptitude is positive")
	}

	b.election.MonApt = monApt
	log.Print("MonApt = " + strconv.Itoa(monApt))
	fmt.Println("L'aptitude de ce processus sera de " + strconv.Itoa(monApt) +
		" pour la prochaine election.")
}

// SetApt used while in election
func (b *BullyImpl) SetApt(processId int, apt int) {
	if processId < 0 || processId >= b.election.N {
		log.Fatal("Process doesn't exist")
	}

	if apt <= 0 {
		log.Fatal("Aptitude is positive")
	}

	b.election.Apts[processId] = apt
	b.processes[processId].Aptitude = apt
	log.Print("Process " + strconv.Itoa(processId) +
		" aptitude was set to " + strconv.Itoa(apt))
}

// EnCours vrai s'il y a une élection en cours
func (b BullyImpl) EnCours() bool {
	return b.election.EnCours
}

// Timeout actions à prendre en fin d'une élection
func (b *BullyImpl) Timeout() {
	b.election.EnCours = false
	b.setElu()
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

	SetTimeout()
}

// GetElu retourne le processus élu lors de la dernière élection
func (b BullyImpl) GetElu() int {
	b.WaitUntilElectionIsOver()
	return b.election.Elu
}

// WaitUntilElectionIsOver bloquer jusqu'à ce que l'élection se termine
func (b *BullyImpl) WaitUntilElectionIsOver(){
	// attente active
	// TODO est-ce qu'il ne vaut mieux pas attendre sur un channel?
	cycles := 0
	oldParticipating := b.participating
	for !b.participating || b.election.EnCours {
		cycles++
		time.Sleep(SLEEP_CYCLE_DURATION) // ne pas tuer la machine

		if !oldParticipating && oldParticipating != b.participating {
			log.Print("Process can now participate")
			b.election.EnCours = false
		}
	}
	waitingTime := time.Duration(cycles) * SLEEP_CYCLE_DURATION
	log.Print("Wait time: " + waitingTime.String())
}

// IsCoordinatorAlive = vrai si l'elu est toujours present
func (b BullyImpl) IsCoordinatorAlive() bool {
	return b.isCoordinatorAlive
}

// IsCoordinatorAlive = vrai si l'elu est toujours present
func (b *BullyImpl) SetIsCoordinatorAlive(s bool) {
	b.isCoordinatorAlive = s
	log.Print("isCoordinatorAlive set to " + strconv.FormatBool(s))
}

// GetCoordinator retourne le processus elu
func (b BullyImpl) GetCoordinator() Entities.Process {
	return b.GetProcess(b.election.Elu)
}

// IsCoordinator = vrai si Moi est l'elu
func (b BullyImpl) IsCoordinator() bool {
	return b.GetMoi() == b.GetCoordinator()
}

// GetMoi retourne le processus qui est lié à moi
func (b BullyImpl) GetMoi() Entities.Process {
	return b.GetProcess(b.election.Moi)
}

// GetProcess
func (b *BullyImpl) GetProcess(id int) Entities.Process {
	return b.processes[id]
}

// RestrictParticipation for one election
func (b *BullyImpl) RestrictParticipation() {
	b.participating = false
	// limite -> election lancee presque au démarrage, ne réponds pas avant 3 t
	time.AfterFunc(3*b.election.T, func() {
		b.participating = true
		log.Print("End of Election participation restriction")
	},
	)
}

// IsParticipating = vrai si le processus peut participer
func (b BullyImpl) IsParticipating() bool {
	return b.participating
}

///* ===============
// * === private ===
// * =============*/

// Message màj de l'aptitude du processus appelant en interne
func (b *BullyImpl) message(processId int, aptitude int){
	if !b.election.EnCours{
		b.Demarre()
	}

	b.election.Apts[processId] = aptitude
}

// GetElu renvoie l'elu avec la règle de départage indiquée
func (b *BullyImpl) setElu() {
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

	b.election.Elu = coordinator
	log.Print("L'elu est le processus " + strconv.Itoa(coordinator))
}
