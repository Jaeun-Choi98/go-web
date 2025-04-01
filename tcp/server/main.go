package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Client struct {
	conn net.Conn
	id   string
}

var mu sync.RWMutex
var clients map[*Client]bool

func main() {
	/*
		explain diff below the code
		addr, _ := net.ResolveTCPAddr("tcp", ":5000")
		listenner, _ := net.ListenTCP("tcp", addr)
	*/
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Println(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	/*
		1.To read a specific number of input bytes, you want io.ReadAtLeast or io.ReadFull.
		To read until some arbitrary condition is met you should just loop on the Read call as long as there is no error.
		(Then you may want to error out on on too-large inputs to prevent a bad client from eating server resources.)

		2.If you're implementing a text-based protocol you should consider net/textproto, which puts a bufio.Reader in front of the connection so you can read lines.

		3.The context package helps manage timeouts, deadlines, and cancellation, and is especially useful if, for example, you're writing a complex server that's going to do many network operations every request
	*/
	receivMessage := func(conn net.Conn) {
		var buf []byte
		for {
			rt, err := conn.Read(buf)

		}
	}

	client := &Client{conn: conn}
	defer conn.Close()
	mu.Lock()
	clients[client] = true
	mu.Unlock()
	fmt.Printf("connected new client: %s\n", conn.RemoteAddr().String())
	receivMessage(conn)
	mu.Lock()
	delete(clients, client)
	mu.Unlock()
}
