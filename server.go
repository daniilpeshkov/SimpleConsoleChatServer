package main

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
	"time"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
)

const (
	INITIAL_CLIENTS_RESERVED_SIZE = 100
)

const (
	TagText     = 1
	TagFileName = 2
	TagFile     = 3
	TagDate     = 4
	TagName     = 5
	TagSys      = 6
)

const (
	SysLoginRequest        = 1
	MinSysLoginRequestSize = 2

	SysLoginResponse = 2

	SysUserLoginNotiffication = 3
	USER_CONNECTED            = 1
	USER_DISCONECTED          = 2
)

const (
	LOGIN_OK  = 1
	NAME_USED = 2
)

type Server struct {
	ln   net.Listener
	port string

	clients     map[string]*simpleTcpMessage.ClientConn
	clientsLock sync.Mutex

	msgChan chan *simpleTcpMessage.Message
}

type loginErrCode byte

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

		logCommand, ok := msg.GetField(TagSys)
		if !ok {
			break
		}

		if len(logCommand) >= MinSysLoginRequestSize && logCommand[0] == SysLoginRequest {

			name = string(logCommand[1:])
			loginRes := server.loginClient(name, clientConn)

			//send response
			msg := simpleTcpMessage.NewMessage()
			msg.AppendField(TagSys, []byte{SysLoginResponse, byte(loginRes)})
			clientConn.SendMessage(msg)

			if loginRes == NAME_USED {
				log.Default().Printf("%s: login refused    [name= %s; ip= %s]\n", time.Now().Format(time.UnixDate), string(name), "NOT IMPLEMENTED")
				continue
			} else if loginRes == LOGIN_OK {
				log.Default().Printf("%s: user connected [name= %s; ip= %s]\n", time.Now().Format(time.UnixDate), string(name), "NOT IMPLEMENTED")

				//tell others obout new user
				msg = simpleTcpMessage.NewMessage()
				msg.AppendField(TagSys, append([]byte{SysUserLoginNotiffication, USER_CONNECTED}, []byte(name)...))
				server.msgChan <- msg

				break
			}
		}

	}
	for {
		msg, err := clientConn.RecieveMessage()

		if err != nil {
			log.Default().Printf("%s: user disconnected [name= %s; ip= %s]\n", time.Now().Format(time.UnixDate), string(name), "NOT IMPLEMENTED")

			//tell others about disconneted user
			msg = simpleTcpMessage.NewMessage()
			msg.AppendField(TagSys, append([]byte{SysUserLoginNotiffication, USER_DISCONECTED}, []byte(name)...))
			server.msgChan <- msg
			break
		}
		msg.RemoveFieldIfExist(TagName)
		msg.AppendField(TagName, []byte(name))
		text, _ := msg.GetField(TagText)

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
			runtime.Gosched()
		}
	}
}
