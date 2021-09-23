package main

import (
	"net"
	"sync"
)

const (
	INITIAL_CLIENTS_SLICE_SIZE = 100
)

type Server struct {
	ln net.Listener

	clients     []*Client
	clientsLock sync.Locker
}

func NewServer(port string) (*Server, error) {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return nil, err
	}
	server := new(Server)
	server.ln = ln
	server.clients = make([]*Client, INITIAL_CLIENTS_SLICE_SIZE)

	server.clientsLock = &sync.Mutex{}

	return server, nil
}

func (server *Server) serveClient(client *Client) {
	//	buf := bytes.NewBuffer(make([]byte, 0, 1000))

}

func (server *Server) RunServer() {
	for {
		conn, err := server.ln.Accept()

		if err != nil {
			return
		}

		client := &Client{
			netIO:    NewNetIO(conn),
			LoggedIn: false,
		}

		server.clientsLock.Lock()
		server.clients = append(server.clients, client)
		server.clientsLock.Unlock()

	}
}
