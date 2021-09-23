package main

import "fmt"

const (
	PORT = "25565"
)

func main() {

	server, err := NewServer(PORT)
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
	}

}
