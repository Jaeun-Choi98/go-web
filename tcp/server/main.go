package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"juchoi/tcp/model/message"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Client struct {
	conn net.Conn
	name string
}

var mu sync.RWMutex
var clients = make(map[*Client]bool)
var messageChan = make(chan string, 10)
var wg sync.WaitGroup

func main() {

	// go func() {
	// 	http.ListenAndServe("0.0.0.0:6060", nil)
	// }()

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Println(err)
		return
	}

	// graceful shut down...
	signChan := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signal.Notify(signChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signChan
		cancel()
		log.Println("shut down server...")
		listener.Close()
		close(messageChan)
		close(signChan)
	}()

	// SendMessage goroutine을 닫기 위한 동기화
	wg.Add(1)
	go SendMessage(ctx)

	timeout, timeoutCancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer timeoutCancel()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			select {
			case <-timeout.Done():
				wg.Wait()
				return
			default:
			}
			continue
		}
		// HandleConnection goroutine을 닫기 위한 동기화
		wg.Add(1)
		go HandleConnection(ctx, conn)
	}
}

func HandleConnection(ctx context.Context, conn net.Conn) {

	client := &Client{conn: conn}
	defer func() {
		conn.Close()
		// 정상 종료 시에 HandleConnection goroutine이 제대로 닫히는지 확인
		log.Println("closed HandleConnection goruntine")
		wg.Done()
	}()
	mu.Lock()
	clients[client] = true
	mu.Unlock()

	handleMessage := func(msg *message.Message) {
		switch msg.Type {
		case "CHECK_TYPE":
			var payload message.CheckPayload
			jsonPayload, err := json.Marshal(msg.Payload)
			if err != nil {
				log.Println(err)
				break
			}
			err = json.Unmarshal(jsonPayload, &payload)
			if err != nil {
				log.Println(err)
				break
			}
			client.name = payload.GetPayload()
			messageChan <- fmt.Sprintf("welcome [%s]\n", client.name)
			/*
				type이 interface{}인 경우, 역직렬화 시 map[string]interface{}로 변환(conversion)된다.
				아래 코드처럼 해도 되지만, 필드가 여러 개일 경우 위의 방법이 더 편리함.
				client.name = msg.Payload.(map[string]interface{})["name"].(string)
			*/

		case "CHAT_MESSAGE_TPYE":
			var payload message.ChatMsgPayload
			jsonPayload, err := json.Marshal(msg.Payload)
			if err != nil {
				log.Println(err)
				break
			}
			err = json.Unmarshal(jsonPayload, &payload)
			if err != nil {
				log.Println(err)
				break
			}
			messageChan <- fmt.Sprintf("[%s]: %s\n", client.name, payload.GetPayload())
		default:
			log.Println("Unknown Message Type")
		}
	}

	/*
		1.To read a specific number of input bytes, you want io.ReadAtLeast or io.ReadFull.
		To read until some arbitrary condition is met you should just loop on the Read call as long as there is no error.
		(Then you may want to error out on on too-large inputs to prevent a bad client from eating server resources.)
		-> e.g. buffer := make([]byte,10) // 10바이트 고정 크기
		 				io.ReadFull(conn, buffer) // 10바이트만 읽음.

		2.If you're implementing a text-based protocol you should consider net/textproto, which puts a bufio.Reader in front of the connection so you can read lines.
		-> 아래 receiveMessage의 경우 bufio.Reader를 사용함.

		3.The context package helps manage timeouts, deadlines, and cancellation, and is especially useful if, for example, you're writing a complex server that's going to do many network operations every request
		-> 아래 주석 코드
	*/
	receiveMessage := func() {
		var msg message.Message
		reader := bufio.NewReader(conn)
		receiveMsgChan := make(chan string, 1)
		defer func() {
			close(receiveMsgChan)
			log.Println("closed receiveMessage func")
		}()

		//3번 방법을 사용한 timeout read 처리
		//timeout, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		//defer timeoutCancel()

		for {
			// if an existing connection was forcibly closed by the remote host, break select block
			closedConnByRemoteHost := make(chan bool, 1)
			defer close(closedConnByRemoteHost)

			go func() {
				line, err := reader.ReadString('\n')
				if err != nil {
					log.Println(err)
					messageChan <- fmt.Sprintf("exit [%s]\n", client.name)
					closedConnByRemoteHost <- true
					return
				}
				jsonMsg := strings.TrimSpace(line)
				receiveMsgChan <- jsonMsg
			}()

			select {
			case <-ctx.Done():
				return
			case <-closedConnByRemoteHost:
				return
			// case <-timeout.Done():
			// 	log.Println("Timeout: No message received")
			// 	return
			case jsonMsg := <-receiveMsgChan:
				err := json.Unmarshal([]byte(jsonMsg), &msg)
				if err != nil {
					log.Println(err)
					continue
				}
				handleMessage(&msg)
			}
		}
	}

	log.Printf("connected new client: %s\n", conn.RemoteAddr().String())
	receiveMessage()
	mu.Lock()
	delete(clients, client)
	mu.Unlock()
	log.Printf("disconnected client: %s\n", conn.RemoteAddr().String())
}

func SendMessage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// 정상 종료 시에 SendMessage goroutine이 제대로 닫히는지 확인.
			log.Println("closed SendMessage goruntine")
			wg.Done()
			return
		case strMsg := <-messageChan:
			mu.RLock()
			for client, exists := range clients {
				if exists {
					client.conn.Write([]byte(strMsg))
				}
			}
			mu.RUnlock()
		}
	}
}
