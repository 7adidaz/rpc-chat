package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

var participant = map[int]string{}
var participantRPC = map[int]*rpc.Client{}
var serverPort = 3000
var maxPort = 3000

type Server int

func (s *Server) Send(args struct {
	Port    int
	Message string
}, ack *bool) error {

	fmt.Println(participant[args.Port], ": ", args.Message)
	message := &struct {
		Name    string
		Message string
	}{participant[args.Port], args.Message}

	for clientPort, clientRPC := range participantRPC {
		if clientPort != args.Port {
			var clientAck bool = false
			clientRPC.Call("Client.Recieve", message, &clientAck)

			// clients disconnected, so remove them
			if !clientAck {
				delete(participant, clientPort)
				delete(participantRPC, clientPort)
			}
		}
	}
	*ack = true
	return nil
}

// assign a port to the client
func (s *Server) Ping(_ struct{}, reply *int) error {
	maxPort++
	*reply = maxPort
	return nil
}

// dial the client and save the connection
func (*Server) Recognise(args struct {
	Port int
	Name string
}, reply *bool) error {
	fmt.Println(args.Name, " Joined the chat.")

	participant[args.Port] = string(args.Name)
	address := "0.0.0.0:" + strconv.Itoa(args.Port)

	client, err := rpc.Dial("tcp", address)

	if err != nil {
		return err
	}
	participantRPC[args.Port] = client
	*reply = true

	return nil
}

func main() {
	address := "0.0.0.0:" + strconv.Itoa(serverPort)
	add, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", add)
	if err != nil {
		log.Fatal(err)
	}

	server := new(Server)
	rpc.Register(server)
	rpc.Accept(inbound)
}
