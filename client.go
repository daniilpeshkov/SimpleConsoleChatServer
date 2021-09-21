package main

import (
	"net"
)

type Client struct {
	NetIO *NetIO
}

func NewClient(conn net.Conn) *Client {
	return &Client{NewNetIO(conn)}
}
