package main

import (
	"fmt"
	"net"
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
	TypeSys      = 5
)

const (
	LoginCode = 1
)

type Server struct {
	ln   net.Listener
	port string

	clients     map[string]*simpleTcpMessage.ClientConn
	clientsLock sync.Mutex

	msgChan chan *simpleTcpMessage.Message
}

type loginErrCode byte

const (
	LOGIN_OK  = 1
	NAME_USED = 2
)

//checks if a client with name exists. If not returns LOGIN_OK else returns LOGIN_ERR
func (server *Server) loginClient(name string, clientConn *simpleTcpMessage.ClientConn) loginErrCode {
	server.clientsLock.Lock()
	defer server.clientsLock.Unlock()

	if _, nameExists := server.clients[name]; nameExists {
		return NAME_USED
	}
	server.clients[name] = clientConn
	return LOGIN_OK
}

func NewServer(port string) *Server {
	return &Server{
		clients: make(map[string]*simpleTcpMessage.ClientConn, INITIAL_CLIENTS_RESERVED_SIZE),
		port:    port,
		msgChan: make(chan *simpleTcpMessage.Message, 100),
	}
}

func (server *Server) RunServer() error {
	var err error
	server.ln, err = net.Listen("tcp", ":"+server.port)
	if err != nil {
		return err
	}
	go server.msgSendGorutine()
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
	var name string
	for {
		msg, err := clientConn.RecieveMessage()
		if err != nil {
			return
		}

		logCommand, ok := msg.GetField(TypeSys)
		if !ok {
			break
		}
		if len(logCommand) > 2 && logCommand[0] == LoginCode {

			name = string(logCommand[1:])
			res := server.loginClient(name, clientConn)

			msg := simpleTcpMessage.NewMessage()
			msg.AppendField(TypeSys, []byte{LoginCode, byte(res)})
			clientConn.SendMessage(msg)

			//todo send everyone that user connected

			if res == NAME_USED {
				fmt.Println("refused login another " + string(name))
				continue
			} else if res == LOGIN_OK {
				fmt.Println(string(name) + " logged in")
				break
			}
		}

	}
	for {
		msg, err := clientConn.RecieveMessage()

		if err != nil {
			//todo send everyone that user disconected
			fmt.Println("User disconneted " + name)
			break
		}
		msg.RemoveFieldIfExist(TypeName)
		msg.AppendField(TypeName, []byte(name))
		text, _ := msg.GetField(TypeText)
		fmt.Printf("%s sent: %s\n", name, string(text))
		server.msgChan <- msg
	}

}

func (server *Server) msgSendGorutine() {
	for {
		select {
		case msg := <-server.msgChan:
			server.clientsLock.Lock()
			for _, v := range server.clients {
				v.SendMessage(msg)
			}
			server.clientsLock.Unlock()
		default:
			continue
		}
	}
}
