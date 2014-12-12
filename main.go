//
//  Hello World server.
//  Binds REP socket to tcp://*:5555
//  Expects "Hello" from client, replies with "World"
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"encoding/json"
	"fmt"
	"time"
)

type Command struct {
	Plugin string
	Args   map[string]interface{}
}

func receiveMessage(msg string) {
	var cmd Command

	if err := json.Unmarshal([]byte(msg), &cmd); err != nil {
		panic(err)
	}

	fmt.Println("Received %v", cmd)
}

func main() {
	//  Socket to talk to clients
	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	responder.Bind("tcp://*:5555")

	for {
		//  Wait for next request from client
		msg, _ := responder.Recv(0)
		receiveMessage(msg)

		//  Do some 'work'
		time.Sleep(time.Second)

		//  Send reply back to client
		reply := "World"
		responder.Send(reply, 0)
	}
}
