package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

const TEST_IP = "127.0.0.1:25565"
const TEST_PORT = 25565

func TestServer1(t *testing.T) {

	go RunServer(TEST_PORT)
	time.Sleep(time.Second * 1)
	time.Sleep(time.Millisecond * 2)
	fmt.Println("trying to connect")
	conn, _ := net.Dial("tcp", TEST_IP)
	msg := []byte{1, 2, 0x15, 0x7D, 0x7E, 4}
	fmt.Println("Sending:")

	for _, v := range msg {
		fmt.Printf("%3x", v)
	}
	fmt.Println()
	msg = ConvertToNet(msg)

	r := io.TeeReader(bytes.NewReader(msg), conn)
	io.ReadAll(r)

	msg = []byte{1, 2, 3, 4, 5}
	msg = ConvertToNet(msg)
	r = io.TeeReader(bytes.NewReader(msg), conn)
	io.ReadAll(r)

	time.Sleep(time.Second * 2)
	conn.Close()
	time.Sleep(time.Second * 1)
}
