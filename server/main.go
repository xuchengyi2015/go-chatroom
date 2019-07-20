package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var connMap = make(map[string]net.Conn)
var messageQueue = make(chan [2]string, 1000)
var quitChannel = make(chan bool)

func checkError(err error) {
	if err != nil {
		fmt.Printf("there is a error: %s\n", err.Error())
		os.Exit(1)
	}
}

func processInfo(conn net.Conn) {
	buf := make([]byte, 1024)
	defer conn.Close()

	for {
		numOfBytes, err := conn.Read(buf)
		if err != nil {
			continue
		}

		if numOfBytes > 0 {
			remoteAddr := conn.RemoteAddr()
			//fmt.Print(remoteAddr)
			//fmt.Printf(": %s\n", string(buf[0:numOfBytes]))
			//
			//conn.Write([]byte(string(buf[:numOfBytes]) + ", too."))
			senderAddr := fmt.Sprint(remoteAddr)
			message := [2]string{senderAddr, string(buf[:numOfBytes])}
			messageQueue <- message
		}
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	checkError(err)
	defer listen.Close()

	fmt.Println("server is waiting ...")

	go handleMessage()
	for {
		conn, err := listen.Accept()
		checkError(err)

		go processInfo(conn)
		reciverAddr := fmt.Sprint(conn.RemoteAddr())
		connMap[reciverAddr] = conn
		for a := range connMap {
			fmt.Printf("%s is connected to server.\n", a)
		}
	}
}

func handleMessage() {
	for {
		select {
		case message := <-messageQueue:
			sendMessage(message)
		case <-quitChannel:
			break
		}
	}
}

func sendMessage(message [2]string) {
	contents := strings.Split(message[1], "#")
	senderAddr := message[0]
	if len(contents) > 1 {
		reciverAddr := contents[0]
		msg := contents[1]

		if conn, ok := connMap[reciverAddr]; ok {
			fmt.Print(message)
			_, err := conn.Write([]byte(senderAddr + ": " + msg))
			if err != nil {
				fmt.Println("online conns send failure.")
			}
		}
	}
}
