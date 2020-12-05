package Entities

type Configuration struct {
	NbProcess int  `json:"nbProcess"`
	Trace	bool   `json:"trace"`
	Debug   bool   `json:"debug"`
	Processes   []Process `json:"Process"`
}