package main

import "testing"

const TEST_IP = "127.0.0.1:25565"
const TEST_PORT = "25565"

func TestLoginClient(t *testing.T) {
	server := NewServer("25565")
	go server.RunServer()

	res := server.loginClient("user1", nil)
	if res != LOGIN_OK {
		t.FailNow()
	}
	res = server.loginClient("user2", nil)
	if res != LOGIN_OK {
		t.FailNow()
	}
	res = server.loginClient("user1", nil)
	if res != LOGIN_ERR {
		t.FailNow()
	}
}

// func TestServer1(t *testing.T) {

// 	go RunServer(TEST_PORT)
// 	time.Sleep(time.Second * 1)
// 	time.Sleep(time.Millisecond * 2)
// 	conn, _ := net.Dial("tcp", TEST_IP)

// 	buf := bytes.NewBuffer([]byte("Чмоха соси хуй"))
// 	netw := NewNetIO(conn)

// 	netw.ReadFrom(buf)

// 	buf = bytes.NewBuffer([]byte("Жуков - пидор"))

// 	netw.ReadFrom(buf)

// 	time.Sleep(time.Second * 2)
// 	conn.Close()
// 	time.Sleep(time.Second * 1)
// }
