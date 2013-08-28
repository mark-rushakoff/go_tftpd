package main

import (
	"log"
	"net"
	"time"

	"github.com/mark-rushakoff/go_tftpd/server_config"
)

func main() {
	bindAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6969")
	if err != nil {
		panic(err.Error())
	}

	udpConn, err := net.ListenUDP("udp", bindAddr)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Listening on %v\n", udpConn.LocalAddr())

	serverConfig := server_config.ServerConfig{
		BlockSize:  512,
		PacketConn: udpConn,
		Timeout:    5 * time.Second,
		TryLimit:   5,
	}

	serverConfig.Serve()
}
