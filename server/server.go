package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
	if err != nil {
		fmt.Println("There is a error!")
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.

	for {
		conn, _ := ln.Accept()
		conns <- conn
	}

}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.

	for {
		//fmt.Println("Awaiting string")
		reader := bufio.NewReader(client)
		msg, _ := reader.ReadString('\n')
		//fmt.Println("Received a message")
		msg1 := Message{
			sender:  clientid,
			message: msg,
		}
		msgs <- msg1
	}

}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, _ := net.Listen("tcp", *portPtr)

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			fmt.Println("Received a new connection")
			clientID := len(clients) + 1
			// - add the client to the clients channel
			clients[clientID] = conn
			// - start to asynchronously handle messages from this client
			go handleClient(conn, clientID, msgs)
		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			//fmt.Println("In the message case")
			fmt.Printf("from client: %d\n", msg.sender)
			fmt.Println(msg.message)
			for clientId := range clients {

				if clientId != msg.sender {
					//fmt.Println("in the if")
					//conn, _ := ln.Accept()
					//fmt.Println(msg.message)
					//fmt.Fprintln(conn, "OK")

					fmt.Fprintln(clients[clientId], msg.message)
					fmt.Fprintln(clients[clientId], msg.sender)
					fmt.Fprintln(clients[clientId], "OK")
				}
			}
		}
	}
}
