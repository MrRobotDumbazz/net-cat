package internal

import (
	"fmt"
	"log"
	"os"
	"time"
)

func ChatMsg(msg string, name string) Message {
	msgTime := time.Now().Format("2006-01-02 15:04:05")
	NowTimeMsg := fmt.Sprintf("[%s][%s]: ", msgTime, name)
	return Message{
		text:    "\r" + NowTimeMsg + msg,
		address: name,
	}
}

func ServiceMsg(msg string, name string) Message {
	return Message{
		text:    "\r" + name + msg,
		address: name,
	}
}

func Broadcaster() {
	dataFile, err := os.Create("data.txt")
	if err != nil {
		log.Print(err)
		return
	}
	for {
		msg := <-Messages
		dataFile.WriteString(msg.text + "\n")
		Mm.Lock()
		for user, conn := range Clients {
			if msg.address == user {
				continue
			}
			if _, err := fmt.Fprintln(conn, "\n"+msg.text); err != nil {
				conn.Write([]byte("Error"))
				conn.Close()

			}
			m := ChatMsg("", user)

			conn.Write([]byte(m.text))
		}
		Mm.Unlock()

	}
}
