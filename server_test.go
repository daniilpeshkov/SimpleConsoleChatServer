package main

const TEST_IP = "127.0.0.1:25565"
const TEST_PORT = "25565"

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
