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
		wg.Add(1)
		conn, err := net.Dial("tcp", TEST_IP)
		if err != nil {
			t.Log(err.Error())
		}
		clientConn := simpleTcpMessage.NewClientConn(conn)
		msg := simpleTcpMessage.NewMessage()
		msg.AppendField(TypeName, []byte(name))
		clientConn.SendMessage(msg)
		time.Sleep(time.Second * 3)
		wg.Done()
	}
	go login("User1")
	go login("User2")
	go login("User1")
	server := NewServer(PORT)
	go server.RunServer()
	wg.Wait()
}
