package main

import (
	"log"
	"net"
	"sync"
	"time"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
)

const (
	INITIAL_CLIENTS_RESERVED_SIZE = 100
)

const (
	TagSys      = 1
	TagMessage  = 2
	TagFileName = 3
	TagFile     = 4
	TagTime     = 5
	TagName     = 6
)

const (
	// request:  should contain TagName with name
	// response: TagSys [SysLoginResponse, LoginStatus]
	SysLoginRequest   = 1
	LOGIN_OK          = 1
	NAME_USED         = 2
	NAME_WRONG_FORMAT = 3

	//request: 	None
	//response: TagSys [SysUserLoginNotiffication, type], TagName and TagTime
	SysUserLoginNotiffication = 3
	USER_CONNECTED            = 1
	USER_DISCONECTED          = 2

	//request: TagMessage TagText
	//response to sender: TagSys [SysMessage, message status], TagTime
	//response to others: TagSys [SysMessage], TagText, TagName, TagTime
	SysMessage           = 4
	MESSAGE_SENT         = 1
	MESSAGE_WRONG_FORMAT = 2
)

type Server struct {
	ln   net.Listener
	port string

	clients     []*Client
	clientsLock sync.Mutex

	usedNames     map[string]struct{}
	usedNamesLock sync.Mutex
	msgChan       chan AddressedMessage
}

type AddrType int

var (
	Broadcast = AddrType(1)
	OnlyTo    = AddrType(2)
	Except    = AddrType(3)
)

type AddressedMessage struct {
	msg      *simpleTcpMessage.Message
	client   *Client
	addrType AddrType
}

type loginErrCode byte

func NewServer(port string) *Server {
	return &Server{
		clients:   make([]*Client, 0, INITIAL_CLIENTS_RESERVED_SIZE),
		port:      port,
		msgChan:   make(chan AddressedMessage, 100),
		usedNames: make(map[string]struct{}),
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
			log.Panicln(err.Error())
		}
		client := NewClient(conn)
		server.clientsLock.Lock()
		server.clients = append(server.clients, client)
		server.clientsLock.Unlock()
		go server.serveClient(client)
	}
}

func (server *Server) msgSendGorutine() {
	for {
		addrMsg := <-server.msgChan
		timeNow := time.Now()
		timeB, _ := timeNow.MarshalBinary()
		addrMsg.msg.AppendField(TagTime, timeB)

		if addrMsg.addrType == OnlyTo {
			addrMsg.client.io.SendMessage(addrMsg.msg)
		} else {
			server.clientsLock.Lock()
			for _, v := range server.clients {
				if v.logined && (addrMsg.addrType == Broadcast || (v != addrMsg.client)) {
					v.io.SendMessage(addrMsg.msg)
				}
			}
			server.clientsLock.Unlock()
		}
	}
}
