package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	lAddr, err := net.ResolveUDPAddr("udp", ":20004")
	checkError(err)
	go receive(lAddr)
	sAddr, err := net.ResolveUDPAddr("udp", "10.100.23.242:20004")
	checkError(err)
	go transmit(sAddr)
	quit := make(chan int)
	<- quit
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func receive(lAddr *net.UDPAddr) {
	conn, err := net.ListenUDP("udp", lAddr)
	checkError(err)

	defer func() {
		err := conn.Close()
		checkError(err)
	}()

	buf := make([]byte, 1024)
	for {
		n,addr,err := conn.ReadFromUDP(buf)
		fmt.Printf("Received %s from %s\n", string(buf[0:n]), addr)
		checkError(err)
	}
}

func transmit(sAddr *net.UDPAddr) {
	conn, err := net.ListenPacket("udp", ":0")
	checkError(err)
	defer func() {
		err := conn.Close()
		checkError(err)
	}()
	for {
		_, err = conn.WriteTo([]byte("data"), sAddr)
		checkError(err)
		time.Sleep(1*time.Second)
	}
}