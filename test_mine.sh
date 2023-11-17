#!/bin/bash

go run server.go & 
sleep 1

gnome-terminal -- /bin/bash -c "go run client.go" & 
gnome-terminal -- /bin/bash -c "go run client.go" & 
gnome-terminal -- /bin/bash -c "go run client.go" & 
gnome-terminal -- /bin/bash -c "go run client.go" & 


sleep 10

echo "thanks"


kill $(jobs -p)
kill $(ps aux | grep 'server' | awk '{print $2}')
