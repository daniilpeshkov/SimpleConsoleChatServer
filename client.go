package main

type Client struct {
	netIO *NetIO
	name  string

	loggedIn bool
}
