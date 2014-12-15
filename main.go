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
	"github.com/mimicloud/easyconfig"
	"log"
	"strings"
)

var config = struct {
	AuthVerbose bool   `json:"auth_verbose"`
	Bind        string `json:"bind"`
	Allow       string `json:"allow"`

	MysqlUsername string `json:"mysql_username"`
	MysqlPassword string `json:"mysql_password"`

	PgsqlUsername string `json:"pgsql_username"`
	PgsqlPassword string `json:"pgsql_password"`

	MongoUsername string `json:"mongo_password"`
	MongoPassword string `json:"mongo_password"`
}{}

var (
	plugins = make(map[string]PluginInterface)
)

func init() {
	easyconfig.Parse("./config.json", &config)

	plugins["mysql"] = &MysqlPlugin{
		Username: config.MysqlUsername,
		Password: config.MysqlPassword,
	}
	plugins["pgsql"] = &PgsqlPlugin{
		Username: config.PgsqlUsername,
		Password: config.PgsqlPassword,
	}
	plugins["mongo"] = &MongoPlugin{
		Username: config.MongoPassword,
		Password: config.MongoPassword,
	}
	plugins["openvz"] = &OpenvzPlugin{}

	// Get some indication of what the authenticator is deciding
	zmq.AuthSetVerbose(config.AuthVerbose)

	// Start the authentication engine. This engine
	// allows or denies incoming connections (talking to the libzmq
	// core over a protocol called ZAP).
	zmq.AuthStart()

	// Whitelist our address; any other address will be rejected
	zmq.AuthAllow(config.Allow)
}

func receiveMessage(msg string) string {
	var cmd Command
	var err error

	log.Printf("%s", msg)
	if err := json.Unmarshal([]byte(msg), &cmd); err != nil {
		return fmt.Sprintf("%v", err)
	}

	if plugin, ok := plugins[cmd.Plugin]; ok {
		if plugin.IsValid(cmd) {
			action := strings.ToLower(cmd.Action)

			if len(action) == 0 {
				return "action not set"
			}

			if action == "create" {
				err = plugin.Create(cmd)
				log.Printf("%v", err)
			} else if action == "delete" {
				err = plugin.Delete(cmd)
			} else if action == "update" {
				err = plugin.Update(cmd)
			} else {
				return fmt.Sprintf("unknown action: %s", cmd.Action)
			}

			if err != nil {
				return fmt.Sprintf("%v", err)
			} else {
				return "ok"
			}
		} else {
			return "Invalid plugin configuration"
		}
	}

	return fmt.Sprintf("Unknown plugin %v", cmd)
}

func main() {
	//  Socket to talk to clients
	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	responder.Bind(config.Bind)

	for {
		//  Wait for next request from client
		msg, _ := responder.Recv(0)
		reply := receiveMessage(msg)
		responder.Send(reply, 0)
	}
}
