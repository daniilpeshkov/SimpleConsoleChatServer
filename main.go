package main

import "fmt"

const (
	PORT = "25565"
)

func main() {
	server := NewServer(PORT)
	err := server.RunServer()
	if err != nil {
		fmt.Println(err.Error())
	}
}
