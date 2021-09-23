package main

import (
	"net"
	"sync"
)

const (
	INITIAL_CLIENTS_RESERVED_SIZE = 100
)

type Server struct {
	ln   net.Listener
	port string

	clients     []*Client
	clientsLock sync.Locker

	usedNames      map[string]struct{}
	usedNameRWlock sync.RWMutex
}

func NewServer(port string) *Server {
	return &Server{
		clients:     make([]*Client, INITIAL_CLIENTS_RESERVED_SIZE),
		clientsLock: &sync.Mutex{},
		usedNames:   make(map[string]struct{}, INITIAL_CLIENTS_RESERVED_SIZE),
	}
}

func (server *Server) RunServer() error {
	var err error
	server.ln, err = net.Listen("tcp", ":"+server.port)
	if err != nil {
		return err
	}

	for {
		conn, err := server.ln.Accept()
		if err != nil {
			return err
		}

		client := &Client{
			netIO:    NewNetIO(conn),
			LoggedIn: false,
		}

		server.clientsLock.Lock()
		server.clients = append(server.clients, client)
		server.clientsLock.Unlock()

		go server.serveClient(client)
	}
}

func (server *Server) serveClient(client *Client) {
	//buf := bytes.NewBuffer(make([]byte, 0, 1000))

	for {

	}
}
