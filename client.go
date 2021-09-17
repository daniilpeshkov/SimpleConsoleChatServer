package main

import (
	"bufio"
	"net"
)

type Client struct {
	ReadWriter *bufio.ReadWriter
}

func NewClient(conn net.Conn) Client {
	return Client{bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))}
}
