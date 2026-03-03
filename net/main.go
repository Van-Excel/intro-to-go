package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Sockets programming")

	addr, err := net.ResolveIPAddr("ip4", "google.com")
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println("resolved address:", addr)

	ipaddr := net.ParseIP("127.1.0.0")
	fmt.Println("parsed IP address:", ipaddr)

	add, err := net.LookupHost("google.com")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("array of addresses for google.com:", add)

	add2, err := net.LookupCNAME("google.com")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("cannonical address for google.com:", add2)

	// building a tcp client

	tcpddress, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8090")
	if err != nil {
		fmt.Println("could initialise tcp address struct")
	}
	newconn, err := net.DialTCP("tcp", nil, tcpddress)
	if err != nil {
		fmt.Println("couldn't connect to server")
	}

	_, err = newconn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	if err != nil {
		fmt.Println("couldn't send message to server")
	}

	buf := make([]byte, 1024)
	count, err := newconn.Read(buf)
	if count > 0 {
		fmt.Println("received message:", string(buf[:count]))

	}
	if count == 0 {
		fmt.Println("reading done")
	}
	if err != nil {
		fmt.Println("reading error:", err)
	}

	//learn to loop reads, close connections

	// build server

}
