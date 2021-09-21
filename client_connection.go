package main

import (
	"net"
)

type ClientConn struct {
	NetIO      *NetIO
	IsLoggedIn bool
	name       string
}

func InitClientConn(conn net.Conn) *ClientConn {
	clientConn := &ClientConn{
		NetIO:      NewNetIO(conn),
		IsLoggedIn: false,
	}
	return clientConn
}
