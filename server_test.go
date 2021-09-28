package main

import (
	"net"
	"sync"
	"testing"
	"time"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
)

const TEST_IP = "127.0.0.1:25565"
const TEST_PORT = "25565"

func TestLoginClient(t *testing.T) {

	wg := sync.WaitGroup{}
	login := func(name string) {
		conn, err := net.Dial("tcp", TEST_IP)
		if err != nil {
			t.Log(err.Error())
		}
		clientConn := simpleTcpMessage.NewClientConn(conn)
		msg := simpleTcpMessage.NewMessage()
		msg.AppendField(TagSys, append([]byte{SysLoginRequest}, []byte(name)...))
		clientConn.SendMessage(msg)
		msg, _ = clientConn.RecieveMessage()
		t.Logf("Login response: %v\n", msg)

		msg, _ = clientConn.RecieveMessage()
		t.Logf("Message: %v\n", msg)

		time.Sleep(time.Second * 3)
		wg.Done()
	}

	wg.Add(2)
	go login("User1")
	go login("User2")

	server := NewServer(PORT)
	go server.RunServer()
	time.Sleep(time.Second * 2)
	msg := simpleTcpMessage.NewMessage()
	msg.AppendField(TagText, []byte("Hello"))
	server.msgChan <- msg
	wg.Wait()
}
