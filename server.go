package main

import (
	"fmt"
	"net"
	"runtime"
	"sync"
)

const (
	INITIAL_CLIENTS_RESERVED_SIZE = 100
)

type Server struct {
	ln   net.Listener
	port string

	clients     map[string]*ClientConn
	clientsLock sync.Mutex
}

type loginErrCode int

const (
	LOGIN_OK  = iota
	LOGIN_ERR = iota
)

//checks if a client with name exists. If not returns LOGIN_OK else returns LOGIN_ERR
func (server *Server) loginClient(name string, clientConn *ClientConn) loginErrCode {
	server.clientsLock.Lock()
	defer server.clientsLock.Unlock()

	if _, nameExists := server.clients[name]; nameExists {
		return LOGIN_ERR
	}

	server.clients[name] = clientConn
	server.clients[name].name = name
	return LOGIN_OK
}

func NewServer(port string) *Server {
	return &Server{
		clients: make(map[string]*ClientConn, INITIAL_CLIENTS_RESERVED_SIZE),
		port:    port,
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

		clientConn := newClientConn(conn)

		go server.serveClient(clientConn)
	}
}

func (server *Server) serveClient(clientConn *ClientConn) {
	for {

		msg, err := clientConn.RecieveMessage()
		if err != nil {
			return
		}
		res := server.loginClient(string(msg.fields[TypeName]), clientConn)
		if res == LOGIN_ERR {
			fmt.Println("refused login another " + string(msg.fields[TypeName]))
			clientConn.conn.Close()
			return
		} else {
			fmt.Println(clientConn.name + " logged in")
			break
		}
	}
	for {
		runtime.Gosched()
	}

}
