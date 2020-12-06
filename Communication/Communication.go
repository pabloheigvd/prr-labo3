package Communication

import (
	"log"
)

/* ==============
 * === Public ===
 * =========== */

// Listen
func Listen(){
	// TODO plrs channel à écouter et message à traiter

	election := make(chan[] string)

	for {
		select {
			case msg := <- election:
				log.Print(msg)
				break
			default:
			break
		}
	}
}


