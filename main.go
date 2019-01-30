package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	lAddr, err := net.ResolveUDPAddr("udp", ":20004")
	checkError(err)
	go udpreceive(lAddr)
	sAddr, err := net.ResolveUDPAddr("udp", "10.100.23.242:20004")
	checkError(err)
	go udptransmit(sAddr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", "10.100.23.242:33546")
	checkError(err)
	go tcpclient(tcpAddr)
	tcpAddr, err = net.ResolveTCPAddr("tcp", ":8000")
	go tcpServer(tcpAddr)
	quit := make(chan int)
	<- quit
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func udpreceive(lAddr *net.UDPAddr) {
	conn, err := net.ListenUDP("udp", lAddr)
	checkError(err)

	defer func() {
		err := conn.Close()
		checkError(err)
	}()

	buf := make([]byte, 1024)
	for {
		n,addr,err := conn.ReadFromUDP(buf)
		fmt.Printf("UDP Listener: %s from %s\n", string(buf[0:n]), addr)
		checkError(err)
	}
}

func udptransmit(sAddr *net.UDPAddr) {
	conn, err := net.ListenPacket("udp", ":0")
	checkError(err)
	defer func() {
		err := conn.Close()
		checkError(err)
	}()
	for {
		_, err = conn.WriteTo([]byte("Hello world"), sAddr)
		checkError(err)
		time.Sleep(1*time.Second)
	}
}

func tcpclient(tcpAddr *net.TCPAddr) {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer func() {
		err := conn.Close()
		checkError(err)
	}()
	conn.LocalAddr().String()
	_, err = conn.Write([]byte("Connect to: 10.100.23.185:8000\000"))
	checkError(err)
	for {
		_, err = conn.Write([]byte("Hello!\000"))
		checkError(err)
		fmt.Println("TCP Client: ", readTCP(conn))
		time.Sleep(1*time.Second)
	}
}

func tcpServer(tcpAddr *net.TCPAddr) {
	conn, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer func() {
		err := conn.Close()
		checkError(err)
	}()
	for {
		conn, err := conn.AcceptTCP()
		checkError(err)
		go func() {
			for {
				fmt.Println("TCP Server: ", readTCP(conn))
				_, err = conn.Write([]byte("World!\000"))
				checkError(err)
				time.Sleep(1 * time.Second)
			}
		}()
	}
}

func readTCP(conn *net.TCPConn) string {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	checkError(err)
	return string(buf)
}