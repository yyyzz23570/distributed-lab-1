package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func readFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		fmt.Print("Server message: " + message)
	}
}

func main() {
	// 使用 flag 包从命令行参数中读取服务器 IP 和端口
	serverIP := flag.String("server-ip", "127.0.0.1", "Server IP address")
	port := flag.String("port", "8030", "Server port")
	flag.Parse()

	// 拼接服务器地址
	address := *serverIP + ":" + *port

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to server at", address)

	go readFromServer(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}
