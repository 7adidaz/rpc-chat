### Event-Driven Chat room using RPC.

How to run it? `$ ./test_mine.sh` or you can run server and multiple clients using `$ go run server.go` and in another terminal `$ go run client.go`.

Server responsplities: 
    - recognise clients (clients ping the server Server.Ping)
    - assign them ports to run on  (Server.Recognise)
    - recieve messages and forward them (Server.Send)

Client: 
    - ping the server for port
    - runs on that port 
    - recieve messages (Client.Recieve)

### Connection Diagram: 

![Connection Diagram](./RPC_connection_diagram.png)
