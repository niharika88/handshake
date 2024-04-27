package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

// readAndParseResponse reads and parses the response from the Bitcoin node.
// It verifies if the expected response was received.
func readAndParseResponse(command Command, conn net.Conn) error {
	// Receive and handle response
	resp := make([]byte, 1024)
	n, err := conn.Read(resp)
	if err != nil {
		return fmt.Errorf("error receiving %s response: %s", command, err)
	}
	log.Printf("%s msg received, READ %d bytes", command, n)

	// Parse response message
	cmd, err := parseMessage(resp[:n])
	if err != nil {
		return fmt.Errorf("error parsing %s response: %s", command, err)
	}

	// Verify if expected response message was received
	if cmd != string(command) {
		return fmt.Errorf("unexpected response, expected %s, received %s msg: ", command, cmd)
	}
	return nil
}

// createVersionMessage creates a version message with the given nonce.
func createVersionMessage(nonce int) ([]byte, error) {
	payload := new(bytes.Buffer)

	err := binary.Write(payload, binary.LittleEndian, int32(ProtocolVersion)) // protocol version
	if err != nil {
		return nil, fmt.Errorf("failed to write protocol version: %s", err)
	}
	err = binary.Write(payload, binary.LittleEndian, uint64(Services)) // kind of services supported
	if err != nil {
		return nil, fmt.Errorf("failed to write local services: %s", err)
	}
	err = binary.Write(payload, binary.LittleEndian, uint64(time.Now().Unix())) // unix timestamp of client machine
	if err != nil {
		return nil, fmt.Errorf("failed to write timestamp: %s", err)
	}
	err = binary.Write(payload, binary.LittleEndian, uint64(0)) // services supported by remote node
	if err != nil {
		return nil, fmt.Errorf("failed to write remote node services: %s", err)
	}
	ip := net.ParseIP(RemoteNodeHost).To16() // convert remote node ip4 to ip6
	_, err = payload.Write(ip)               // set ip6 in payload
	if err != nil {
		return nil, fmt.Errorf("failed to write remote node host: %s", err)
	}
	err = binary.Write(payload, binary.BigEndian, uint16(RemoteNodePort)) // remote node port
	if err != nil {
		return nil, fmt.Errorf("failed to write remote node port: %s", err)
	}
	err = binary.Write(payload, binary.LittleEndian, uint64(0)) // list of local host services
	if err != nil {
		return nil, fmt.Errorf("failed to write localhost services: %s", err)
	}
	localIp := net.ParseIP(Localhost).To16() // convert local node ip4 to ip6
	_, err = payload.Write(localIp)          // local host IP, same as remote node for this test
	if err != nil {
		return nil, fmt.Errorf("failed to write localhost ip: %s", err)
	}
	err = binary.Write(payload, binary.BigEndian, uint16(LocalhostPort)) // local host port, same as remote node for this test
	if err != nil {
		return nil, fmt.Errorf("failed to write localhost port: %s", err)
	}
	err = binary.Write(payload, binary.LittleEndian, uint64(nonce)) // Nonce
	if err != nil {
		return nil, fmt.Errorf("failed to write nonce value: %s", err)
	}
	err = payload.WriteByte(byte(len(UserAgent))) // indicate length of the upcoming UserAgent string
	if err != nil {
		return nil, fmt.Errorf("failed to write user agent size: %s", err)
	}
	_, err = payload.WriteString(UserAgent) // UserAgent
	if err != nil {
		return nil, fmt.Errorf("failed to write user agent: %s", err)
	}
	err = binary.Write(payload, binary.LittleEndian, int32(StartHeight)) // We don't have any block to share
	if err != nil {
		return nil, fmt.Errorf("failed to write StartHeight for last block: %s", err)
	}
	err = binary.Write(payload, binary.LittleEndian, uint8(Relay)) // Announce relayed transactions
	if err != nil {
		return nil, fmt.Errorf("failed to write relay boolean: %s", err)
	}

	// Create the version message from magic bytes + command + payload size + checksum
	// Add HEADER
	message := new(bytes.Buffer)

	_, err = message.Write(MagicBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write magic bytes: %s", err)
	}

	_, err = message.WriteString(string(Version))
	if err != nil {
		return nil, fmt.Errorf("failed to write version command string: %s", err)
	}
	_, err = message.Write(make([]byte, 12-len(Version))) // padding on the right with empty bytes.
	if err != nil {
		return nil, fmt.Errorf("failed to add padding to version cmd string: %s", err)
	}

	err = binary.Write(message, binary.LittleEndian, uint32(payload.Len()))
	if err != nil {
		return nil, fmt.Errorf("failed to write payload size: %s", err)
	}

	hash := sha256.Sum256(payload.Bytes())
	hash = sha256.Sum256(hash[:])
	_, err = message.Write(hash[:4])
	if err != nil {
		return nil, fmt.Errorf("failed to write checksum: %s", err)
	}

	// Add PAYLOAD
	_, err = message.Write(payload.Bytes()) // add the actual payload to the version message
	if err != nil {
		return nil, fmt.Errorf("failed to write payload: %s", err)
	}

	return message.Bytes(), nil
}

// createVerackMessage creates a verack message.
func createVerackMessage() ([]byte, error) {
	message := new(bytes.Buffer)

	_, err := message.Write(MagicBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write magic bytes: %s", err)
	}

	_, err = message.WriteString(string(Verack))
	if err != nil {
		return nil, fmt.Errorf("failed to write verack command: %s", err)
	}

	_, err = message.Write(make([]byte, 12-len(Verack))) // padding on the right with empty bytes.
	if err != nil {
		return nil, fmt.Errorf("failed to add verack cmd padding: %s", err)
	}

	// Empty payload for verack message.

	return message.Bytes(), nil
}

// parseMessage parses a Bitcoin message and returns the command.
// also prints out the received response and payload to the stdout.
func parseMessage(data []byte) (string, error) {
	if len(data) < 24 {
		return "", fmt.Errorf("message too short")
	}

	magic := data[:4] // First 4 bytes are magic bytes.
	if !bytes.Equal(magic, MagicBytes) {
		return "", fmt.Errorf("invalid magic bytes")
	}

	command := string(bytes.TrimRight(data[4:16], "\x00")) // Next 12 bytes contain the command.

	s := binary.LittleEndian.Uint32(data[16:20]) // These 4 bytes tell the payload size.
	if len(data[24:]) < int(s) {
		return "", fmt.Errorf("payload length mismatch")
	}

	payload := data[24 : 24+s]

	log.Printf(
		"received message with:\n magicBytes: %v\n command: %s\n payloadLen: %d\n checksum: %v\n payload: %v\n",
		magic,
		command,
		s,
		data[20:24],
		payload,
	) // Empty payload for VERACK message.

	return command, nil
}
