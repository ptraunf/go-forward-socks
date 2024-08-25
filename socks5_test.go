package main

import (
	"errors"
	"log"
	"net"
	"testing"
	"time"
)

const CLIENT_ADDR = "192.168.1.111"

type mockConn struct {
	deadline      time.Time
	readDeadline  time.Time
	writeDeadline time.Time
	receive       chan []byte
	send          chan []byte
	closed        bool
}

func newMockConn() *mockConn {
	s := make(chan []byte, 1)
	r := make(chan []byte, 1)
	timeout := time.Time{}
	return &mockConn{
		deadline:      timeout,
		readDeadline:  timeout,
		writeDeadline: timeout,
		receive:       r,
		send:          s,
		closed:        false,
	}
}

func (mc mockConn) Read(b []byte) (n int, err error) {
	if !mc.closed {
		bytes, ok := <-mc.receive
		if !ok {
			return 0, errors.New("mockConn: Could not read from receive channel!")
		}
		n = copy(b, bytes)
		return n, nil
	} else {
		return 0, errors.New("mockConn: Attempted to read from closed connection!")
	}
}
func (mc mockConn) Write(b []byte) (n int, err error) {
	// n = copy(mc.send, b)
	if !mc.closed {
		mc.send <- b
		return len(b), nil
	} else {
		return 0, errors.New("Attempted to read from closed mockConn!")
	}
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (mc mockConn) Close() error {
	if mc.closed {
		return errors.New("mockConn: Attempted to close already-closed connection.")
	}
	mc.closed = true
	return nil
}

type addr struct {
	network string
	str     string
}

func (a addr) Network() string {
	return a.network
}
func (a addr) String() string {
	return a.str
}

// LocalAddr returns the local network address, if known.
func (mc mockConn) LocalAddr() net.Addr {
	return addr{
		network: "tcp",
		str:     "127.0.0.1:53388",
	}
}

// RemoteAddr returns the remote network address, if known.
func (mc mockConn) RemoteAddr() net.Addr {
	return addr{
		network: "tcp",
		str:     "127.0.0.1:8080",
	}

}

func (mc mockConn) SetDeadline(t time.Time) error {
	mc.readDeadline = t
	mc.writeDeadline = t
	return nil
}

func (mc mockConn) SetReadDeadline(t time.Time) error {
	mc.readDeadline = t
	return nil
}

func (mc mockConn) SetWriteDeadline(t time.Time) error {
	mc.writeDeadline = t
	return nil
}

// type byteIndex int
// const (
// 	versionIdx byteIndex iota,
// 	replyCodeIdx,
// 	resvIdx,
// 	addressTypeIdx,

// )
func TestGenerateReply(t *testing.T) {
	var fields = map[int]string{
		0: "Version",
		1: "ReplyCode",
		2: "Reserved",
		3: "AddressType",
		4: "Address First Byte",
		5: "Address Second Byte",
		6: "Address Third Byte",
		7: "Address Fourth Byte",
		8: "Port First Byte",
		9: "Port Second Byte",
	}
	var code replyCode = Succeeded
	var addrType addressFamily = IPv4
	var address = []byte{byte(127), 0x00, 0x00, byte(1)}
	var port = []byte{0x00, byte(80)}
	expected := []byte{
		Version,
		byte(code),
		0x00,
		byte(addrType),
		address[0],
		address[1],
		address[2],
		address[3],
		port[0],
		port[1],
	}
	actual := generateReply(code, addrType, address, port)

	if actual[0] != Version {

	}
	// if actual[1] !=
	for idx, name := range fields {
		if expected[idx] != actual[idx] {
			t.Errorf("\nExpected %v:\t%v\nActual %v:\t%v", name, expected[idx], name, actual[idx])
		}
	}
}
func TestGenerateFailedReply(t *testing.T) {

}

func TestHandleSocks5(t *testing.T) {
	conn := newMockConn()
	conn.receive <- []byte{Version, byte(Connect), Resv, byte(IPv4), 0x0, byte(80)}
	// handleSocks5(conn)
	go func(sent chan []byte) {
		for bytes := range sent {
			log.Printf("Bytes sent to client: % x", bytes)
		}
	}(conn.send)
}
