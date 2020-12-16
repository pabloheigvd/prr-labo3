package Communication

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"prr-labo3/Entities"
)

// ListenToRemoteMessage from other process
func ListenToRemoteMessage(moiP Entities.Process){
	conn, err := net.ListenPacket(TRANSPORT_PROTOCOL, moiP.GetAdress())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// from udp slide
	buf := make([]byte, 1024)
	for {
		log.Print("Waiting for remote message...")
		n, cliAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		msg := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for msg.Scan() {
			msg := msg.Text()
			log.Print("Received: " + msg + " from " + cliAddr.String() + "\n")
			go func() {messageChannel <- msg}()
		}
		log.Print("Finished parsing remote msg")
	}
}