package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

func serveClient(client *ClientConn) {
	buf := bytes.NewBuffer(make([]byte, 0, 1000))
	for {
		buf.Reset()
		_, err := buf.ReadFrom(client.NetIO)
		if errors.Is(err, EOP{}) {
			fmt.Printf("[%s]: %s\n", client.name, string(buf.Bytes()))

			client.NetIO.ReadFrom(buf)
			fmt.Println()
		} else {
			fmt.Println("Connection closed")
			break
		}
	}
}

func RunServer(port string) {
	fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", ":"+port)

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

		client := InitClientConn(conn)
		buf := make([]byte, 100)

		name_len, _ := client.NetIO.Read(buf)

		client.name = string(buf[:name_len])

		fmt.Printf("<%s joined>\n", client.name)
		go serveClient(client)

	}
}
