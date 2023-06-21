package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"TCPChat/internal"
)

func main() {
	var port string
	if len(os.Args) == 1 {
		port = "8989"
	} else if len(os.Args) == 2 && internal.IsValidInput(string(os.Args[1])) {
		port = os.Args[1]
	} else {
		log.Println("[USAGE]: ./TCPChat $port")
		return
	}
	TCPserver(port)
}

func TCPserver(port string) {
	addr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		log.Print(err)
		return
	}
	listener, err := net.ListenTCP("tcp", addr)
	defer listener.Close()
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println("Listening on the port :", port)
	go internal.Broadcaster()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			conn.Write([]byte("client has left the chat..."))
			if err := conn.Close(); err != nil {
				log.Print(err)
				return
			}
			continue
		}
		internal.Mm.Lock()
		
		if internal.ClientId > 9 {
			conn.Write([]byte("\nserver is full. Try again later..."))
			if err := conn.Close(); err != nil {
				log.Print(err)
			}
			internal.Mm.Unlock()
			continue
		}
		internal.Mm.Unlock()
		go internal.Handler(conn)
	}
}
