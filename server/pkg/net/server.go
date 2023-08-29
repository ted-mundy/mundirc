package net

import (
	"fmt"
	"net"
	"strings"
)

var server = Server{}

var CLIENT_CONNECT = byte(0x0)
var CLIENT_DISCONNECT = byte(0x1)

func Listen(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	fmt.Printf("listening on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("+connection (%s)\n", conn.RemoteAddr().String())

	nickValid := false
	var nick string
	var client Client

	ack := byte(0x6)

	for !nickValid {
		conn.Write([]byte{ack})
		conn.Write([]byte("Enter your nick: "))
		buf := make([]byte, 16)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("error reading:", err.Error())
			return
		}

		nick = string(buf[:n])
		nick = strings.Trim(nick, "\n")
		nick = strings.Trim(nick, "\r")
		nickValid = ValidNick(&server, nick)
	}

	id := uint16(len(server.clients))
	client = Client{id, conn, nick}
	server.clients = append(server.clients, client)
	defer handleDisconnect(id)

	fmt.Printf("client %s connected with nick %s\n", conn.RemoteAddr().String(), nick)

	go announceConnection(client)

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("error reading:", err.Error())
			return
		}

		handleMessage(client, string(buf[:n]))
	}
}

func handleDisconnect(id uint16) {
	// client := server.clients[index]
	var client Client
	var index int
	for i, _client := range server.clients {
		if _client.id == id {
			client = _client
			index = i
		}
	}
	defer fmt.Printf("-connection (%s)\n", client.PrettyPrint())

	client.conn.Close()

	// Remove the client from the clients, so we aren't referencing them anymore in
	// potential calls (such as checking for duplicate nicks)
	server.clients = append(server.clients[:index], server.clients[index+1:]...)
	output := []byte{CLIENT_DISCONNECT}
	output = append(output, []byte(client.PrettyPrint())...)
	for _, _client := range server.clients {
		_client.conn.Write(output)
	}
}

func announceConnection(client Client) {
	output := []byte{CLIENT_CONNECT}
	output = append(output, []byte(client.PrettyPrint())...)
	for _, _client := range server.clients {
		_client.conn.Write(output)
	}
}

func handleMessage(client Client, message string) {
	message = strings.Trim(message, "\n")
	message = strings.Trim(message, "\r")
	fmt.Printf("<%s>: %s\n", client.PrettyPrint(), message)

	for _, _client := range server.clients {
		_client.conn.Write([]byte(fmt.Sprintf("<%s>: %s\n", client.PrettyPrint(), message)))
	}
}
