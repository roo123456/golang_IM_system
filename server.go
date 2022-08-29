package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//online user list
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//channel to broadcast message
	Message chan string
}

// NewServer create a server interface
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

//the method to broadcast message
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg

}

// ListenMessager a goroutine to listen channel that storages broadcast message , it will send message to all online users when message comes
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//broadcast message to all online users by putting message to single user's message channel
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// Handler what handler do actually
func (this *Server) Handler(conn net.Conn) {

	user := NewUser(conn)

	//user hop online,put user to onlineMap
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	//broadcast a message that current user hop online
	this.Broadcast(user, " has already hoped online")

	//current handler blocked
	select {}
}

// Start the interface of starting server
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listen err : ", err)
	}
	//close listen socket
	defer listener.Close()

	//start a goroutine to listen message
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err : ", err)
		}
		//do handler
		go this.Handler(conn)
	}
}
