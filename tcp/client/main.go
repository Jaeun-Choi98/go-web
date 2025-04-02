package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter message: ")
		text, _ := reader.ReadString('\n')
		conn.Write([]byte(text))

		serverResponse, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Server response: ", serverResponse)
	}
}
