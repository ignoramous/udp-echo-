package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	port int = 5000
)

func init() {
	if v := os.Getenv("ECHO_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("can't parse ECHO_PORT: %s", err)
		}

		port = p
	}
}

func main() {
	udp, err := net.ListenPacket("udp", fmt.Sprintf("fly-global-services:%d", port))
	if err != nil {
		log.Fatalf("can't listen on %d/udp: %s", port, err)
	}

	handleUDP(udp)
}

func handleUDP(c net.PacketConn) {
	packet := make([]byte, 2000)

	for {
		n, addr, err := c.ReadFrom(packet)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			log.Printf("error reading on %d/udp: %s", port, err)
			continue
		}

		c.WriteTo(packet[:n], addr)
	}
}
