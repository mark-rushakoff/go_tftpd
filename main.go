package main

import (
	"log"
	"net"
	"os"
	"os/signal"
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

	// handle ctrl-c
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		for sig := range c {
			log.Printf("Received %v, exiting", sig)
			os.Exit(0)
		}
	}()

	serverConfig := server_config.ServerConfig{
		PacketConn: udpConn,
		Timeout:    1 * time.Second,
		TryLimit:   2,
	}

	serverConfig.Serve()
}
