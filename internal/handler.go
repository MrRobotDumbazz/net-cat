package internal

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

func Handler(conn net.Conn) {
	var name string

	logo, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Print(err)
	}
	Mm.Lock()
	ClientId += 1
	Mm.Unlock()
	conn.Write(logo)
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	reader := bufio.NewReader(conn)
	defer conn.Close()
	for {
		name, err = reader.ReadString('\n')
		if err != nil {
			conn.Write([]byte("Error in ReadString"))
			conn.Close()
			break
		}
		name = strings.TrimSuffix(name, "\n")
		if CheckUsername(name, conn) {
			break
		}
	}

	Mm.Lock()

	Clients[name] = conn
	Mm.Unlock()

	Messages <- ServiceMsg(" has joined our chat...", name)
	data, err := os.ReadFile("data.txt")
	if err != nil {
		log.Print(err)
	}

	conn.Write(data)
	m := ChatMsg("", name)
	conn.Write([]byte(m.text))
	input := bufio.NewScanner(conn)
	for input.Scan() {
		temp := CheckMsg(input.Text())
		if strings.ReplaceAll(temp, " ", "") != "" {
			Messages <- ChatMsg(temp, name)
		}
		m := ChatMsg("", name)
		conn.Write([]byte(m.text))
	}
	Mm.Lock()
	delete(Clients, m.address)
	ClientId -= 1
	Mm.Unlock()

	Messages <- ServiceMsg(" has left our chat...", name)
	if err := conn.Close(); err != nil {
		log.Print(err)
		return
	}
}
