package main

import (
	"log"
	"net"
	"time"

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

func (server *Server) serveClient(client *Client) {
	for {
		msg, err := client.io.RecieveMessage()

		if client.logined {
			if err != nil { // if error when loggined
				server.logoutClient(client)
				log.Default().Printf("%s: user disconnected [name= %s; ip= %s]\n", time.Now().Format(time.UnixDate), string(client.name), "NOT IMPLEMENTED")
				log.Default().Printf("Online user count: %d\n", len(server.clients))

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
				log.Default().Printf("%s: user message [name= %s; message= %s ip= %s]\n", time.Now().Format(time.UnixDate), string(client.name), text, "NOT IMPLEMENTED")

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
				client.io.SendMessage(msg)

				if loginRes == NAME_USED || loginRes == NAME_WRONG_FORMAT {
					log.Default().Printf("%s: login refused    [name= %s; ip= %s]\n", time.Now().Format(time.UnixDate), string(client.name), "NOT IMPLEMENTED")
					continue
				} else if loginRes == LOGIN_OK {
					log.Default().Printf("%s: user connected [name= %s; ip= %s]\n", time.Now().Format(time.UnixDate), string(client.name), "NOT IMPLEMENTED")

					//tell others obout new user
					msg = simpleTcpMessage.NewMessage()
					msg.AppendField(TagSys, []byte{SysUserLoginNotiffication, USER_CONNECTED})
					msg.AppendField(TagName, []byte(client.name))
					server.msgChan <- AddressedMessage{msg, nil, Broadcast}

					break
				}
			}
		}

	}

}
