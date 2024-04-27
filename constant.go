package main

import "math/rand"

const (
	ProtocolVersion = 70015           // a set protocol version
	Services        = 0               // leaving as zero since we are testing but it can be the list of services our nodes offers.
	UserAgent       = "/Niharika:24/" // a nice user agent name
	RemoteNodeHost  = "127.0.0.1"     // Change the address and port if necessary, ideally this should be part of a config file.
	RemoteNodePort  = 8333
	Localhost       = "127.0.0.1"
	LocalhostPort   = 8333
	StartHeight     = 0
	Relay           = 1
)

type Command string

const (
	Version Command = "version"
	Verack  Command = "verack"
)

var (
	Nonce      = rand.Int() // Nonce is set to a random integer
	MagicBytes = []byte{0xf9, 0xbe, 0xb4, 0xd9}
)
