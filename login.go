package main

import (
	"fmt"
	"strings"
)

const (
	MaxNameLength = 20
)

func nameCorrect(name string) bool {
	return len(name) > 0 && len(name) <= MaxNameLength && !strings.ContainsAny(name, "\n\x00") &&
		!strings.HasPrefix(name, " ") && !strings.HasSuffix(name, " ")
}

//checks if a client with name exists. If not returns LOGIN_OK else returns LOGIN_ERR
func (server *Server) loginClient(client *Client) loginErrCode {
	server.usedNamesLock.Lock()
	defer server.usedNamesLock.Unlock()

	if nameCorrect(client.name) {
		if _, nameExists := server.usedNames[client.name]; nameExists {
			return NAME_USED
		}
		server.usedNames[client.name] = struct{}{}
		client.logined = true

		return LOGIN_OK
	}
	return NAME_WRONG_FORMAT
}

func (server *Server) logoutClient(client *Client) {
	server.clientsLock.Lock()
	server.usedNamesLock.Lock()
	defer server.clientsLock.Unlock()
	defer server.usedNamesLock.Unlock()

	for i, v := range server.clients {
		if v == client {
			server.clients[i] = server.clients[len(server.clients)-1]
			server.clients = server.clients[:len(server.clients)-1]
			break
		}
	}
	fmt.Printf("after disconnect %d\n", len(server.clients))
	delete(server.usedNames, client.name)
}
