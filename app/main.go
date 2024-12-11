package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	PORT = 4421
	ADDRESS = "127.0.0.1"
)

func main() {
	listener, err := net.Listen("tcp", ADDRESS + ":" + strconv.Itoa(PORT))
	if err != nil {
		log.Fatalln("Error starting tcp server: ", err)
	}	
	defer listener.Close()

	fmt.Printf("Running tcp server on address: %s:%d\n", ADDRESS, PORT)
	
	for {
		conn, err := listener.Accept()
		
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		
		fmt.Println(conn.RemoteAddr())
		conn.Close()
	}
}