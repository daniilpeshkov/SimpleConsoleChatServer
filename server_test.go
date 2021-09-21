package main

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"
)

const TEST_IP = "127.0.0.1:25565"
const TEST_PORT = "25565"

func TestServer1(t *testing.T) {

	go RunServer(TEST_PORT)
	time.Sleep(time.Second * 1)
	time.Sleep(time.Millisecond * 2)
	fmt.Println("trying to connect")
	conn, _ := net.Dial("tcp", TEST_IP)

	buf := bytes.NewBuffer([]byte{1, 2, 3, 4, 5})

	fmt.Println("Sending:")
	for _, v := range buf.Bytes() {
		fmt.Printf("%3x", v)
	}
	fmt.Println()
	netw := NewNetIO(conn)

	netw.ReadFrom(buf)

	buf = bytes.NewBuffer([]byte{0x15, 0x7d, 0x7e})

	fmt.Println("Sending:")
	for _, v := range buf.Bytes() {
		fmt.Printf("%3x", v)
	}

	fmt.Println()
	netw.ReadFrom(buf)

	time.Sleep(time.Second * 2)
	conn.Close()
	time.Sleep(time.Second * 1)
}
