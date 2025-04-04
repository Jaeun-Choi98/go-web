package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"juchoi/tcp/model/message"
	"log"
	"net"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

/*
can make console ui using 'tview'
*/

func main() {

	// go func() {
	// 	http.ListenAndServe("0.0.0.0:6060", nil)
	// }()

	pctx, pcancel := context.WithCancel(context.Background())
	defer pcancel()

	// require modifying, bad code
	bufioReaderReadStringBreak := make(chan struct{}, 1)
	defer close(bufioReaderReadStringBreak)

	waitForSendMessage := func(conn net.Conn, ctx context.Context) {
		reader := bufio.NewReader(os.Stdin)
		sendMsgChan := make(chan message.Message, 5)
		defer func() {
			close(sendMsgChan)
			log.Println("closed waitForSendMessage")
			// require modifying, bad code
			bufioReaderReadStringBreak <- struct{}{}
		}()

		for {
			go func() {
				line, err := reader.ReadString('\n')
				if err != nil {
					log.Println(err)
					return
				}
				chatMsg := message.Message{
					Type: "CHAT_MESSAGE_TPYE",
					Payload: &message.ChatMsgPayload{
						Message: strings.TrimSpace(line),
					},
				}
				// require modifying, bad code
				select {
				case <-bufioReaderReadStringBreak:
					log.Println("closed bufioReaderReadString block")
					return
				default:
				}
				sendMsgChan <- chatMsg
			}()

			select {
			case <-ctx.Done():
				return
			case data := <-sendMsgChan:
				err := SendMessage(conn, data)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}

	connect := func() {
		reader := bufio.NewReader(os.Stdin)
		ctx, cancel := context.WithCancel(pctx)
		defer func() {
			cancel()
			log.Println("closed connect")
		}()
		conn, err := net.DialTimeout("tcp", "localhost:5000", time.Second*5)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
		// require modifying, bad code
		fmt.Print("If it's recovery connection, push Enter Before write your name\nYour Name: ")
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
		go ReceiveMessage(conn, cancel)
		waitForSendMessage(conn, ctx)
	}

	recoverTimeInterval := 2 * time.Second

	//graceful shut down...
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		pcancel()
	}()

	connect()
	// Recovery 시도
	recoveryCnt := 0
	for {
		if recoveryCnt >= 5 {
			log.Println("Failed to recover connect")
			pcancel()
			break
		}
		select {
		case <-time.After(recoverTimeInterval):
			connect()
			recoveryCnt++
		case <-pctx.Done():
			return
		}
	}
}

func ReceiveMessage(conn net.Conn, cancel context.CancelFunc) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			cancel()
			return
		}
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
