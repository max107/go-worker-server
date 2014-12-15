package main

import (
	"labix.org/v2/mgo"

	"fmt"
	"time"
)

type MongoPlugin struct {
	Hostname string
	Port     int
	Username string
	Password string

	cmd_username string
	cmd_password string
	cmd_database string
}

func (self *MongoPlugin) GetType() string {
	return "mongo"
}

func (self *MongoPlugin) IsValid(cmd Command) bool {
	username, _ := getString(cmd, "username")
	if len(username) == 0 {
		return false
	}
	self.cmd_username = username

	password, _ := getString(cmd, "password")
	if len(password) == 0 {
		return false
	}
	self.cmd_password = password

	database, _ := getString(cmd, "database")
	if len(database) == 0 {
		return false
	}
	self.cmd_database = database

	return true
}

func (self *MongoPlugin) Create(cmd Command) error {
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

	database.UpsertUser(user)
	session.Close()

	return nil
}

func (self *MongoPlugin) Delete(cmd Command) error {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", self.Username, self.Password, self.Hostname, self.Port, self.cmd_database)
	session, err := mgo.DialWithTimeout(url, time.Second)
	if err != nil {
		return err
	}
	database := session.DB(self.cmd_database)

	err = database.RemoveUser(self.cmd_username)
	if err != nil {
		return err
	}

	err = database.DropDatabase()
	if err != nil {
		return err
	}

	session.Close()

	return nil
}

func (self *MongoPlugin) Update(cmd Command) error {
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

	database.UpsertUser(user)

	session.Close()

	return nil
}
