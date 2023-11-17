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
