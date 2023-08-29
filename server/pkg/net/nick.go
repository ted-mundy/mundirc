package net

import (
	"fmt"
	"strings"
)

func ValidNick(server *Server, nick string) bool {
	_nick := strings.ToLower(nick)
	for _, client := range server.clients {
		if strings.ToLower(client.nick) == _nick {
			return false
		}
	}

	return true
}

func (client *Client) PrettyPrint() string {
	address := strings.Split(client.conn.RemoteAddr().String(), ":")[0]
	return fmt.Sprintf("%s!%s", client.nick, address)
}
