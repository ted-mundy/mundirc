package net

import "net"

type Server struct {
	clients []Client
}

type Client struct {
	id   uint16
	conn net.Conn
	nick string
}
