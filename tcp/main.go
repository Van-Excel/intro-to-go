package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	userBuffer := make([]byte, 4096)

	n, err := conn.Read(userBuffer)
	if err != nil {
		fmt.Println("socket buffer could not be read")
	}
	fmt.Println("Received:")
	fmt.Println(string(userBuffer[:n]))
	body := "Thank you for subscribing"
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
		"\r\n" +
		body
	conn.Write([]byte(response))

}

func main() {

	fmt.Println("Building a server")

	listener, err := net.Listen("tcp", ":8080") // calls sock.Socket(), bind()
	if err != nil {
		fmt.Println("error during connection")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConnection((conn))
	}

}
