package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type method int

const (
	NoAuthRequired      method = 0x00
	GSSAPI              method = 0x01
	UserPass            method = 0x02
	NoAcceptableMethods method = 0xFF
)

type cmd int

const (
	Connect      cmd = 0x01
	Bind         cmd = 0x02
	UDPAssociate cmd = 0x03
)

type addressFamily int

const (
	IPv4       addressFamily = 0x01
	DomainName addressFamily = 0x03
	IPv6       addressFamily = 0x04
)

type socksNegotiation []byte

func (n socksNegotiation) parse() (int, int) {
	return 0, 0
}
func handleClient()
func run() {
	port := 8080
	address := fmt.Sprintf(":%v", port)
	log.Printf("Listening at %v", address)
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Printf(err.Error())
			return
		}
		if strings.TrimSpace(string(netData)) == "STOP" {
			log.Printf("Exiting")
			return
		}

		log.Printf("-> %v", string(netData))
		t := time.Now()
		response := "Responded at " + t.Format(time.RFC3339) + "\n"
		log.Printf("<- %v", response)
		c.Write([]byte(response))
	}

}
func main() {
	log.Printf("-> GO FWD SOCKS5 ->")
	run()
}
