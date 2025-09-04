package room

import "tcp-groupchat/internal/user"

type Room struct {
	RoomId int
	Users  []*user.User
}

func NewRoom(roomId int) *Room {
	return &Room{RoomId: roomId}
}

func (r *Room) AddUser(user *user.User) {
	r.Users = append(r.Users, user)
}
