package main

import (
	"log"
	"net"
	"os"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
)

type Client struct {
	io      *simpleTcpMessage.ClientConn
	logined bool
	ipAddr  string
	name    string
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		io:      simpleTcpMessage.NewClientConn(conn),
		logined: false,
		ipAddr:  "",
		name:    "",
	}
}

var logger = log.New(os.Stdout, "", log.Ltime|log.Ldate)

func (server *Server) serveClient(client *Client) {
	for {
		msg, err := client.io.RecieveMessage()
		if client.logined {
			if err != nil { // if error when loggined
				server.logoutClient(client)
				logger.Printf("user disconnected [%s]\n", string(client.name))
				//tell others about disconneted user
				msg = simpleTcpMessage.NewMessage()
				msg.AppendField(TagSys, []byte{SysUserLoginNotiffication, USER_DISCONECTED})
				msg.AppendField(TagName, []byte(client.name))
				server.msgChan <- AddressedMessage{msg, nil, Broadcast}
				break
			}
			if isMessageRequest(msg) {

				rspMsg := simpleTcpMessage.NewMessage()
				rspMsg.AppendField(TagSys, []byte{SysMessage, MESSAGE_SENT})
				//todo time
				server.msgChan <- AddressedMessage{rspMsg, client, OnlyTo}

				msg.RemoveFieldIfExist(TagName)
				msg.AppendField(TagName, []byte(client.name))
				text, _ := msg.GetField(TagMessage)
				logger.Printf("[message] %s:  %s\n", string(client.name), text)

				server.msgChan <- AddressedMessage{msg, client, Except}
			}
		} else {
			if err != nil {
				return // if error before loging nothing to do
			}
			if isLoginRequest(msg) {
				nameBytes, _ := msg.GetField(TagName)
				client.name = string(nameBytes)
				loginRes := server.loginClient(client)

				msg := simpleTcpMessage.NewMessage()
				msg.AppendField(TagSys, []byte{SysLoginRequest, byte(loginRes)})
				server.msgChan <- AddressedMessage{msg, client, OnlyTo}

				if loginRes == NAME_USED || loginRes == NAME_WRONG_FORMAT {
					logger.Printf("login refused [%s]\n", string(client.name))
					continue
				} else if loginRes == LOGIN_OK {
					logger.Printf("user connected [%s]\n", string(client.name))

					//tell others obout new user
					msg = simpleTcpMessage.NewMessage()
					msg.AppendField(TagSys, []byte{SysUserLoginNotiffication, USER_CONNECTED})
					msg.AppendField(TagName, []byte(client.name))
					server.msgChan <- AddressedMessage{msg, client, Except}
				}
			}
		}

	}

}
