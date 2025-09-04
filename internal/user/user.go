package user

import "net"

type User struct {
	Username string
	Conn     net.Conn
}

func NewUser(username string, Conn net.Conn) *User {
	return &User{
		Username: username,
		Conn:     Conn,
	}
}
