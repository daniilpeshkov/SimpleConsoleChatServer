package main

import "net"

type ClientConn struct {
	conn net.Conn
	name string
}

func newClientConn(conn net.Conn) *ClientConn {
	return &ClientConn{
		conn: conn,
	}
}
