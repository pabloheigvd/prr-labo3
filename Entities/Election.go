/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	Election.go
 */

package Entities

import "time"

type Election struct {
	N			int
	Moi   		int
	MonApt   	int
	Apts		[]int
	T 			time.Duration
	EnCours		bool
	Elu			int
}

// GetProcess qui est moi
func (e Election) GetProcess(processes []Process) Process {
	return processes[e.Moi]
}
