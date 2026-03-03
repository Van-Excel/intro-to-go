package main

import (
	"fmt"
	"net"
)

func main() {

	tcplistener, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8090")
	if err != nil {
		fmt.Println("tcp server not found")

	}
	servertcplistener, err := net.ListenTCP("tcp", tcplistener)
	if err != nil {
		fmt.Println("server not listening")

	}
	serverconn, err := servertcplistener.Accept()

	serverbuf := make([]byte, 1024)

	count, err := serverconn.Read(serverbuf)
	fmt.Println("server received this message:", string([]byte(serverbuf[:count])))
	defer serverconn.Close()
	if err != nil {
		fmt.Println("server couldnt read")
	}

	serverconn.Write([]byte("HTTP/1.0 200 OK \r\n" +
		"ETag: -9985996 \r\n" +
		"Last-Modified: Thu, 25 Mar 2010 17:51:10 GMT \r\n" +
		"Content-Length: 18074 \r\n" +
		"Connection: close \r\n" +
		"Date: Sat, 28 Aug 2010 00:43:48 GMT \r\n" +
		"Server: lighttpd/1.4.23 \r\n" +
		"\r\n"))
}
