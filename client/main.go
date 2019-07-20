package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const address string = "127.0.0.1:8080"

func checkError(err error) {
	if err != nil {
		fmt.Printf("there is a error: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	conn, err := net.Dial("tcp", address)
	checkError(err)
	defer conn.Close()

	go sendMessage(conn)

	buf := make([]byte, 1024)
	for {
		numOfBytes, err := conn.Read(buf)
		if err != nil {
			continue
		}

		if numOfBytes > 0 {
			fmt.Printf("%s\n", string(buf[:numOfBytes]))
		}
	}
	fmt.Println("client send a msg.")
}

func sendMessage(conn net.Conn) {
	defer conn.Close()

	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)

		if strings.ToUpper(input) == "EXIT" {
			fmt.Printf("client:%s is exited.\n", address)
			conn.Close()
			break
		}

		_, err := conn.Write(data)
		if err != nil {
			fmt.Printf("client:%s connection is failure.\n", err.Error())
			conn.Close()
			os.Exit(0)
			break
		}
	}
}
