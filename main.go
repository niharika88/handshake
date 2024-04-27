package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	log.Println("Initiating connection to bitcoin node ...")
	// Connect to a Bitcoin Core node
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", RemoteNodeHost, RemoteNodePort))
	if err != nil {
		log.Fatalln("Error connecting to bitcoin node: ", err)
	}
	defer conn.Close()

	// --------------- Step 1 ---------------
	// Send version message
	vMsg, err := createVersionMessage(Nonce)
	if err != nil {
		log.Fatalln("Error creating version message: ", err)
	}
	n, err := conn.Write(vMsg)
	if err != nil {
		log.Fatalln("Error sending version message: ", err)
	}
	log.Printf("Version msg sent, WROTE %d bytes", n)

	// --------------- Step 2 ---------------
	// Handle 'version' message
	err = readAndParseResponse(Version, conn)
	if err != nil {
		log.Fatalln(err)
	}

	// --------------- Step 3 ---------------
	// Handle 'verack' message
	err = readAndParseResponse(Verack, conn)
	if err != nil {
		log.Fatalln(err)
	}

	// --------------- Step 4 ---------------
	// Send verack message to complete the handshake
	ackMsg, err := createVerackMessage()
	if err != nil {
		log.Fatalln("Error creating verack message: ", err)
	}
	n, err = conn.Write(ackMsg)
	if err != nil {
		log.Fatalln("Error sending verack message:", err)
	}
	log.Printf("Verack msg sent, WROTE %d bytes", n)

	log.Println("Handshake is successful!!")
}
