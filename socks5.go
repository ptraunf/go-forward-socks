package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const Version = 0x05
const Resv = 0x00

type method byte

const (
	NoAuthRequired      method = 0x00
	GSSAPI              method = 0x01
	UserPass            method = 0x02
	NoAcceptableMethods method = 0xFF
)

type cmd byte

const (
	Connect      cmd = 0x01
	Bind         cmd = 0x02
	UDPAssociate cmd = 0x03
)

type addressFamily byte

const (
	IPv4       addressFamily = 0x01
	DomainName addressFamily = 0x03
	IPv6       addressFamily = 0x04
)

type replyCode byte

const (
	Succeeded               replyCode = 0x00
	GenralFailure           replyCode = 0x01
	ConnectionNotAllowed    replyCode = 0x02
	NetworkUnreachable      replyCode = 0x03
	HostUnreachable         replyCode = 0x04
	ConnectionRefused       replyCode = 0x05
	TTLExpired              replyCode = 0x06
	CommandNotSupported     replyCode = 0x07
	AddressTypeNotSupported replyCode = 0x08
)

type socksNegotiation []byte

func (n socksNegotiation) parse() (int, int) {
	return 0, 0
}

func handleBasic(c net.Conn) {
	defer c.Close()
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

func readBytes(r *bufio.Reader, n int) ([]byte, error) {
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return res, err
		} else {
			res[i] = b
		}
	}
	return res, nil
}

func hasByte(a []byte, value byte) bool {
	for _, b := range a {
		if b == value {
			return true
		}
	}
	return false
}

func handleSocks5(c net.Conn) {
	defer c.Close()
	log.Printf("-> SOCKS5 Client Connection:\t%v", c.RemoteAddr())
	connReader := bufio.NewReader(c)
	negotiation, err := readBytes(connReader, 2)
	if err != nil {
		log.Printf("Error Reading Bytes:\n%v", err.Error())
	}

	version := negotiation[0]
	log.Printf("Version: % x", version)
	nMethods := int(negotiation[1])
	log.Printf("nMethods: % x", nMethods)
	methods, err := readBytes(connReader, nMethods)
	if err != nil {
		log.Printf("Error Reading Methods:\n%v", err.Error())
	}
	if hasByte(methods, byte(NoAuthRequired)) {
		log.Printf("Method: 0%2x", methods[0])
	} else {
		log.Printf("Expected method 0x00: NO AUTH REQUIRED")
		return
	}
	methodReply := []byte{Version, byte(NoAuthRequired)}
	log.Printf("Writing bytes: % x ", methodReply)
	c.Write(methodReply)

	request, err := readBytes(connReader, 4)
	if err != nil {
		log.Printf("Error Reading Request:\n%v", err.Error())
	}

	version = request[0]
	cmd := request[1]
	// 3rd byte is RESERVED 0x00 byte
	addrType := request[3]
	var address string
	if addrType == byte(IPv4) {
		log.Println("Address Type: IPv4")
		ipv4, err := readBytes(connReader, 4)
		if err != nil {
			log.Printf("Error Reading IPv4 Addr:\n%v", err.Error())
		}
		address = fmt.Sprintf("%v.%v.%v.%v", ipv4[0], ipv4[1], ipv4[2], ipv4[3])
	} else if addrType == byte(DomainName) {
		log.Println("Address Type: DomainName")
		domainLen, err := connReader.ReadByte()
		if err != nil {
			log.Printf("Error Reading Domain Len:\n%v", err.Error())
		}
		domain, err := readBytes(connReader, int(domainLen))
		ips, err := net.LookupIP(string(domain))
		address = ips[0].String()
	} else {
		log.Printf("Error: Unsupported Address Type: % x", addrType)
		return
	}
	portBytes, err := readBytes(connReader, 2)
	port := binary.BigEndian.Uint16(portBytes)
	address = fmt.Sprintf("%v:%v", address, port)
	log.Printf("Address:\t%v", address)
	var bindAddress string
	if cmd == byte(Connect) {
		remote, err := net.Dial("tcp", address)
		if err != nil {
			log.Printf("Error dialing remote address (%v):\n%v", address, err.Error())
			return
		}
		bindAddress = remote.LocalAddr().String()
	} else {
		log.Printf("Error: Command %v not supported", cmd)
		c.Write(generateFailureReply(CommandNotSupported, addressFamily(addrType)))
		return
	}
	bindAddressParts := strings.Split(bindAddress, ":")
	bindHost := net.ParseIP(bindAddressParts[0])
	bindPort, err := strconv.ParseUint(bindAddressParts[1], 10, 16)
	if err != nil {
		bindPort = 1080
	}
	portBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, uint16(bindPort))

	reply := generateReply(Succeeded, IPv4, bindHost, portBytes)
	c.Write(reply)
}
func exchange(client, remote net.Conn) {
	// buffSize := 4096

	// go func() {
	// 	for {

	// 	}
	// }()
}

func generateReply(code replyCode, addrType addressFamily, address []byte, port []byte) []byte {
	reply := make([]byte, 10)
	reply[0] = Version
	reply[1] = byte(code)
	reply[2] = 0x00
	reply[3] = byte(addrType)

	reply[4] = address[0]
	reply[5] = address[1]
	reply[6] = address[2]
	reply[7] = address[3]
	reply[8] = port[0]
	reply[9] = port[1]
	return reply
	// return []byte{byte(Version), byte(code), 0x00, byte(addrType), address[:], port[:]...}
}

func generateFailureReply(errorCode replyCode, addrType addressFamily) []byte {
	return generateReply(errorCode, addrType, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
}
func handleConnection(c net.Conn) {
	// handleBasic(c)
	handleSocks5(c)
}
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
	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf(err.Error())
			return
		}
		go handleConnection(c)
	}
}
func main() {
	log.Printf("-> GO FWD SOCKS5 ->")
	run()
}
