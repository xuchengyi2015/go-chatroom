package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var connMap = make(map[string]net.Conn)
var messageQueue = make(chan [2]string, 1000)
var quitChannel = make(chan bool)

var logFile *os.File
var logger *log.Logger

const (
	log_file_dic = "./log.txt"
	tcp          = "tcp"
	ip           = "127.0.0.1:8080"
)

func checkError(err error) {
	if err != nil {
		fmt.Printf("there is a error: %s\n", err.Error())
		os.Exit(1)
	}
}

func processInfo(conn net.Conn) {
	buf := make([]byte, 1024)

	//if conn has error, it should be removed by connMap.
	defer func(conn net.Conn) {
		addr := fmt.Sprint(conn.RemoteAddr())
		delete(connMap, addr)
		conn.Close()

		fmt.Println("this is time, online client:")
		for v := range connMap {
			fmt.Println(v)
		}
	}(conn)

	for {
		numOfBytes, err := conn.Read(buf)
		if err != nil {
			continue
		}

		if numOfBytes > 0 {
			remoteAddr := conn.RemoteAddr()

			senderAddr := fmt.Sprint(remoteAddr)
			message := [2]string{senderAddr, string(buf[:numOfBytes])}
			messageQueue <- message
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
		//if contents contains '#', it should be joined.
		msg := strings.Join(contents[1:], "#")

		if conn, ok := connMap[reciverAddr]; ok {
			fmt.Print(message)
			_, err := conn.Write([]byte(senderAddr + ": " + msg))
			if err != nil {
				fmt.Println("online conns send failure.")
			}
		}
	}
}

func main() {
	//log is vital importance!
	logFile, err := os.OpenFile(log_file_dic, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("create log file is failure.%s\n", err.Error())
		os.Exit(-1)
	}
	defer logFile.Close()

	logger = log.New(logFile, "\n", log.Ldate|log.Ltime|log.Llongfile)
	logger.Println("hello world")

	listen, err := net.Listen(tcp, ip)
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
