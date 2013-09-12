package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/mark-rushakoff/go_tftpd/serverconfig"
)

var host string
var port int

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Host to use for server")
	flag.IntVar(&port, "port", 69, "Port to use for server")
}

func main() {
	flag.Parse()

	bindAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, strconv.Itoa(port)))
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
		PacketConn:     udpConn,
		DefaultTimeout: time.Second,
		TryLimit:       2,
	}

	serverConfig.Serve()
}
