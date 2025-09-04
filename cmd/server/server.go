package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"tcp-groupchat/internal/room"
	"tcp-groupchat/internal/user"
)

var mutex = &sync.Mutex{}

var rooms []*room.Room

func randInt() int {
	return rand.IntN(10000)
}

func checkValidRoom(roomId int) (*room.Room, error) {
	for _, r := range rooms {
		if r.RoomId == roomId {
			return r, nil
		}
	}
	return nil, fmt.Errorf("invalid room id")
}

func broadcastRoomUsers(msg string, r *room.Room, sender *user.User) {
	for _, roomUser := range r.Users {
		if sender.Conn != roomUser.Conn {
			roomUser.Conn.Write([]byte(fmt.Sprintf("%s : %s", sender.Username, msg)))
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	conn.Write([]byte("Enter your username : "))
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error in Reading the string", err)
		return
	}
	username = strings.TrimSpace(username)
	currentUser := user.NewUser(username, conn)

	conn.Write([]byte("Choose an option : \n1.Create a Room \n2. Join an Existing Room\n"))
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error in Reading the string", err)
		return
	}
	input = strings.TrimSpace(input)
	option, err := strconv.Atoi(input)
	if err != nil {
		conn.Write([]byte("Invalid input. Please enter 1 or 2.\n"))
		return
	}
	switch option {
	case 1:
		room := room.NewRoom(randInt())
		room.AddUser(currentUser)
		mutex.Lock()
		rooms = append(rooms, room)
		mutex.Unlock()
		conn.Write([]byte(fmt.Sprintf("Your room id is : %d\n", room.RoomId)))
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				conn.Write([]byte("Couldn't read the message\n"))
				broadcastRoomUsers("has left the room\n", room, currentUser)
				return
			}
			broadcastRoomUsers(input, room, currentUser)
		}
	case 2:
		conn.Write([]byte("Enter the room id\n"))
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error in Reading the string", err)
			return
		}
		input = strings.TrimSpace(input)
		roomid, err := strconv.Atoi(input)
		if err != nil {
			conn.Write([]byte("Invalid input. Please enter a valid room id\n"))
			return
		}
		room, err := checkValidRoom(roomid)
		if err != nil {
			conn.Write([]byte("Invalid input. Please enter a valid room id\n"))
			return
		}
		room.AddUser(currentUser)
		broadcastRoomUsers(fmt.Sprintf("%s has joined the room\n", currentUser.Username), room, currentUser)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				conn.Write([]byte("Couldn't read the message\n"))
				broadcastRoomUsers("has left the room\n", room, currentUser)
				return
			}
			broadcastRoomUsers(input, room, currentUser)
		}
	}

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "42069"
	}
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server listening on : ", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error in accepting connection : ", err)
		}
		fmt.Println("Connection established with : ", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
