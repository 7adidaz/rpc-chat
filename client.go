package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"

	"os"
	"strconv"
)

/*
TODO
there are two implementations
1) pooling
	- the client will dial the rpc of the coordinating server
	- the client will call the remote procedure on the server to send a message
	- the client can fetch all of the messages history from the server using remote procedure call


2) event-driven [BONUS]
	- a client starts by looking for a port to establish it's server on (like giving my phone number to my friends to call me)
	- a client can send a message through an infinite loop waiting for input text, this message will be broadcasted to other clients through an rpc call on the server
	- a client can also receive messages simultaneously using the GO keyword
	- so a client here is a server + a client at the same time
*/

type Client int

func (c *Client) Recieve(
	args struct {
		Name    string
		Message string
	}, ack *bool) error {

	*ack = true
	fmt.Println(args.Name, ": ", args.Message)
	fmt.Print("> ")
	return nil
}

func main() {
	server, err := rpc.Dial("tcp", "0.0.0.0:3000")
	if err != nil {
		log.Fatal("Error Connecting to the server.")
	}

	// assign me a port.
	var port int
	err = server.Call("Server.Ping", struct{}{}, &port)
	if err != nil {
		log.Fatal("Ping err: ", err)
	}

	// run in this port.
	address := "0.0.0.0:" + strconv.Itoa(port)
	add, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", add)
	if err != nil {
		log.Fatal(err)
	}

	client := new(Client)
	rpc.Register(client)
	go rpc.Accept(inbound)

	//------------------

	in := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	name, _, err := in.ReadLine()

	id := &struct {
		Port int
		Name string
	}{port, string(name)}
	var reply bool = false

	// call the server to recognise me
	// as running on port x with name y
	err = server.Call("Server.Recognise", id, &reply)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Print("> ")
		line, _, err := in.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		var ack bool
		message := &struct {
			Port    int
			Message string
		}{port, string(line)}

		err = server.Call("Server.Send", message, &ack)
		if err != nil {
			log.Fatal("send: ", err)
		}
	}
}
