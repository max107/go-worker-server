package main

import (
	"labix.org/v2/mgo"

	"errors"
	"fmt"
	"log"
	"time"
)

type OpenvzPlugin struct {
	container *Container
}

func (self *OpenvzPlugin) GetType() string {
	return "openvz"
}

func (self *OpenvzPlugin) getString(cmd Command, key string) (string, error) {
	if value, ok := cmd.Args[key]; ok {
		return fmt.Sprintf("%s", value), nil
	} else {
		return "", errors.New("Unknown key")
	}
}

func (self *OpenvzPlugin) IsValid(cmd Command) bool {
	username, _ := self.getString(cmd, "username")
	if len(username) == 0 {
		return false
	}
	self.cmd_username = username

	password, _ := self.getString(cmd, "password")
	if len(password) == 0 {
		return false
	}
	self.cmd_password = password

	database, _ := self.getString(cmd, "database")
	if len(database) == 0 {
		return false
	}
	self.cmd_database = database

	return true
}

func (self *OpenvzPlugin) Process(cmd Command) error {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", self.Username, self.Password, self.Hostname, self.Port, self.cmd_database)
	session, err := mgo.DialWithTimeout(url, time.Second)
	if err != nil {
		return err
	}
	database := session.DB(self.cmd_database)

	user := &mgo.User{
		Username: self.cmd_username,
		Password: self.cmd_password,
		Roles:    []mgo.Role{mgo.RoleReadWrite},
	}

	err := database.UpsertUser(user)
	if err != nil {
		return err
	}

	session.Close()

	return nil
}
