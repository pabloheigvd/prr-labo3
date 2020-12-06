/*
 * Work: 	PRR-labo3
 * Author: 	Pablo Mercado
 * File: 	Configuration.go
 */

package Entities

type Configuration struct {
	NbProcess 	int  		`json:"nbProcess"`
	Trace		bool   		`json:"trace"`
	Debug   	bool   		`json:"debug"`
	Processes   []Process 	`json:"processes"`
}