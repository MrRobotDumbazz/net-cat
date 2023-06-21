package internal

import (
	"net"
	"sync"
)

var (
	Clients  = make(map[string]net.Conn)
	ClientId = 0
	Messages = make(chan Message)
	Mm       sync.Mutex
)

type Message struct {
	text    string
	address string
}

func CheckUsername(name string, conn net.Conn) bool {
	for _, checksym := range name {
		if checksym < ' ' || checksym > '~' {
			_, err := conn.Write([]byte("Not printable symbols. Please write username again. \n[ENTER YOUR NAME]:"))
			if err != nil {
				return false
			}
			return false
		}
	}
	flag := true

	for _, v := range name {
		if v != ' ' {
			flag = false
			break
		}
	}
	if flag {
		_, err := conn.Write([]byte("You writing only space. Please write username again. \n[ENTER YOUR NAME]:"))
		if err != nil {
			return false
		}
		return false
	}
	for user := range Clients {
		if user == name {
			_, err := conn.Write([]byte("Name is taken... \n[ENTER YOUR NAME]:"))
			if err != nil {
				return false
			}
			return false
		}
	}
	if len(name) == 0 {
		_, err := conn.Write([]byte("Enter your name... \n[ENTER YOUR NAME]:"))
		if err != nil {
			return false
		}
		return false
	}
	return true
}

func IsValidInput(port string) bool {
	for _, r := range port {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func CheckMsg(s string) string {
	temp := ""
	for _, check := range s {
		if check > 31 {
			temp += string(check)
		}
	}
	return temp
}
