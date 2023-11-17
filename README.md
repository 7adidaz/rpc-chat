### Event-driven chat room using RPC.

How to run it? 
- `$ ./test_mine.sh` or
-  you can run a server and multiple clients using `$ go run server.go` and in another terminal `$ go run client.go`.

Server responsibilities: 
- recognize clients (clients ping the server Server.Ping)
- assign them ports to run on  (Server.Recognise)
- receive messages and forward them (Server.Send)

Client: 
- Ping the server for port
- runs on that port 
- receive messages (Client.Recieve)

### Connection Diagram: 

![Connection Diagram](./RPC_connection_diagram.png)
