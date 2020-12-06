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
