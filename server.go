package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
)

func serveClient(client Client) {
	nr := NetReader{client.ReadWriter.Reader}
	buf := bytes.NewBuffer(make([]byte, 0, 1000))
	for {
		buf.Reset()
		_, err := buf.ReadFrom(nr)
		if errors.Is(err, EOP{}) {
			fmt.Println("Recieved:")
			for _, v := range buf.Bytes() {
				fmt.Printf("%3x", v)
			}
			fmt.Println()
		} else {
			fmt.Println("Connection closed")
			break
		}
	}
}

func RunServer(port uint) {
	fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", ":"+strconv.Itoa(int(port)))

	for {

		conn, err := ln.Accept()
		var opError *net.OpError

		if err != nil {
			if errors.As(err, &opError) {
				fmt.Println(opError.Err.Error())
			} else {
				fmt.Println(err.Error())
				continue
			}
		}
		fmt.Println("New Connection!")
		client := NewClient(conn)
		go serveClient(client)
	}
}
