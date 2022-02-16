package main

import (
	"fmt"
	"net"
	"time"
)

func readTCP(conn net.Conn, r chan string) {
	fmt.Println("readTCP")
	buf := make([]byte, 2048)
	for {
		n, _ := conn.Read(buf)
		r <- string(buf[:n])
	}
}

func writeTCP(conn net.Conn, w chan string) {
	for {
		msg := <-w
		conn.Write(append([]byte(msg), 0))
		fmt.Print("Write TCP")
		//time.Sleep(time.Second)
	}
}

func writeSomething(write_ch chan string) {
	for {
		write_ch <- "1110"
		time.Sleep(time.Second)
		fmt.Print("Write stuff <3")
	}
}

func main() {
	server_addr_str := "127.0.0.1:15657"
	server_addr, err := net.ResolveTCPAddr("tcp", server_addr_str)
	if err != nil {
		fmt.Print("err 3")
		panic(err)
	}
	// local_addr_str := "10.100.23.245:30000"
	// local_addr, _ := net.ResolveTCPAddr("tcp", local_addr_str)

	conn, err := net.DialTCP("tcp", nil, server_addr)
	if err != nil {
		fmt.Print("err 1")
		panic(err)
	}

	defer conn.Close()
	ln, err := net.Listen("tcp", ":9004")
	if err != nil {
		fmt.Print("err 2")
		panic(err)
	}

	defer ln.Close()

	read_ch := make(chan string)
	write_ch := make(chan string)
	//write_ch2 := make(chan string)
	//go readTCP(conn, read_ch)
	go writeTCP(conn, write_ch)
	go readTCP(conn, read_ch)
	//go writeTCP(conn, write_ch2)
	go writeSomething(write_ch)

	write_ch <- "Connect to: 10.100.23.245:9004\000"

	//write_ch <- "Heihei!"

	for {
		fmt.Println("start for-loop")
		c, _ := ln.Accept()

		go readTCP(c, read_ch)
		go writeTCP(c, write_ch)
		go writeSomething(write_ch)

		for {
			select {
			case msg := <-read_ch:
				fmt.Println(msg)
			}
		}
	}
}
