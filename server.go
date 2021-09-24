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

	clients     map[string]*Client
	clientsLock sync.Mutex
}

type loginErrCode int

const (
	LOGIN_OK  = iota
	LOGIN_ERR = iota
)

//checks if a client with name exists. If not returns LOGIN_OK else returns LOGIN_ERR
func (server *Server) loginClient(name string, client *Client) loginErrCode {
	server.clientsLock.Lock()
	defer server.clientsLock.Unlock()

	if _, nameExists := server.clients[name]; nameExists {
		return LOGIN_ERR
	}

	server.clients[name] = client
	return LOGIN_OK
}

func NewServer(port string) *Server {
	return &Server{
		clients: make(map[string]*Client, INITIAL_CLIENTS_RESERVED_SIZE),
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
			netIO: NewNetIO(conn),
		}

		go server.serveClient(client)
	}
}

func (server *Server) serveClient(client *Client) {
	//buf := bytes.NewBuffer(make([]byte, 0, 1000))

}
