package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var clients = make(map[net.Conn]bool) // 保存所有客户端连接
var mutex = &sync.Mutex{}             // 用于保护 clients 共享资源

// 广播消息给所有客户端，除了发送者
func broadcastMessage(sender net.Conn, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		if client != sender { // 不发送给发送者
			_, err := client.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Println("Error broadcasting to client:", err)
				client.Close()
				delete(clients, client) // 移除连接断开的客户端
			}
		}
	}
}

// 处理单个客户端连接
func handleClient(conn net.Conn) {
	defer conn.Close()
	mutex.Lock()
	clients[conn] = true // 将新连接的客户端加入列表
	mutex.Unlock()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading message:", err)
			mutex.Lock()
			delete(clients, conn) // 移除断开的客户端
			mutex.Unlock()
			return
		}
		message := string(buf[:n])
		fmt.Println("Received message:", message)

		// 广播消息给其他客户端
		broadcastMessage(conn, strings.TrimSpace(message))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8030")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8030...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}
