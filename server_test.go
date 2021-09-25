package main

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

const TEST_IP = "127.0.0.1:25565"
const TEST_PORT = "25565"

func TestLoginClient(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)
	login := func(name string) {
		conn, err := net.Dial("tcp", TEST_IP)
		if err != nil {
			t.Log(err.Error())
		}
		clientConn := newClientConn(conn)
		msg := NewMessage()
		msg.appendField(TypeName, []byte(name))
		clientConn.SendMessage(msg)
		for {

		}
	}
	go login("User1")
	go login("User2")
	go login("User1")
	server := NewServer(PORT)
	go server.RunServer()
	wg.Wait()
}

func TestSend(t *testing.T) {
	wg := sync.WaitGroup{}
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		fmt.Println(err.Error())
	}
	wg.Add(1)
	go func() {
		conn, err := net.Dial("tcp", TEST_IP)
		if err != nil {
			fmt.Println(err.Error())
		}
		clientConn := newClientConn(conn)
		msg := NewMessage()

		msg.appendField(TypeName, []byte("Жуков"))
		msg.appendField(TypeText, []byte("Я чмоха"))
		clientConn.SendMessage(msg)

		wg.Done()
	}()

	conn, _ := ln.Accept()
	clientConn := newClientConn(conn)

	msg, _ := clientConn.RecieveMessage()

	t.Log(msg)
	wg.Wait()
}
