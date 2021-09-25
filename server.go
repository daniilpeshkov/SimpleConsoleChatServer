package main

import (
	"fmt"
	"net"
	"runtime"
	"sync"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
)

const (
	INITIAL_CLIENTS_RESERVED_SIZE = 100
)

const (
	TypeText     = 0
	TypeFileName = 1
	TypeFile     = 2
	TypeDate     = 3
	TypeName     = 4
)

type Server struct {
	ln   net.Listener
	port string

	clients     map[string]*simpleTcpMessage.ClientConn
	clientsLock sync.Mutex
}

type loginErrCode int

const (
	LOGIN_OK  = iota
	LOGIN_ERR = iota
)

//checks if a client with name exists. If not returns LOGIN_OK else returns LOGIN_ERR
func (server *Server) loginClient(name string, clientConn *simpleTcpMessage.ClientConn) loginErrCode {
	server.clientsLock.Lock()
	defer server.clientsLock.Unlock()

	if _, nameExists := server.clients[name]; nameExists {
		return LOGIN_ERR
	}
	server.clients[name] = clientConn
	return LOGIN_OK
}

func NewServer(port string) *Server {
	return &Server{
		clients: make(map[string]*simpleTcpMessage.ClientConn, INITIAL_CLIENTS_RESERVED_SIZE),
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

		clientConn := simpleTcpMessage.NewClientConn(conn)

		go server.serveClient(clientConn)
	}
}

func (server *Server) serveClient(clientConn *simpleTcpMessage.ClientConn) {
	for {

		msg, err := clientConn.RecieveMessage()
		if err != nil {
			return
		}

		name, ok := msg.GetField(TypeName)
		if ok {
			res := server.loginClient(string(name), clientConn)
			if res == LOGIN_ERR {
				fmt.Println("refused login another " + string(name))

				clientConn.Close()
				return
			} else {
				fmt.Println(string(name) + " logged in")
				break
			}
		}
	}
	for {
		runtime.Gosched()
	}

}
