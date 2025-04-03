package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"juchoi/tcp/model/message"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

/*
can make console ui using 'tview'
*/

func main() {

	conn, err := net.DialTimeout("tcp", "localhost:5000", time.Second*5)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Your Name: ")
	line, _ := reader.ReadString('\n')
	checkMsg := message.Message{
		Type: "CHECK_TYPE",
		Payload: &message.CheckPayload{
			Name: strings.TrimSpace(line),
		},
	}
	err = SendMessage(conn, checkMsg)
	if err != nil {
		log.Println(err)
	}

	go ReceiveMessage(conn)

	for {
		line, _ := reader.ReadString('\n')
		chatMsg := message.Message{
			Type: "CHAT_MESSAGE_TPYE",
			Payload: &message.ChatMsgPayload{
				Message: strings.TrimSpace(line),
			},
		}
		err := SendMessage(conn, chatMsg)
		if err != nil {
			log.Println(err)
		}
	}
}

func ReceiveMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		line, _ := reader.ReadString('\n')
		fmt.Print(line)

	}
}

func SendMessage(conn net.Conn, msg message.Message) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	jsonMsg = append(jsonMsg, '\n')
	_, err = conn.Write(jsonMsg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
