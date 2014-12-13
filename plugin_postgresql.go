package main

import (
	"github.com/jinzhu/gorm"

	_ "github.com/lib/pq"

	"errors"
	"fmt"
	"log"
)

type PgsqlPlugin struct {
	Hostname string
	Port     int
	Username string
	Password string

	cmd_username string
	cmd_password string
	cmd_database string
}

func (self *PgsqlPlugin) GetType() string {
	return "pgsql"
}

func (self *PgsqlPlugin) getString(cmd Command, key string) (string, error) {
	if value, ok := cmd.Args[key]; ok {
		return fmt.Sprintf("%s", value), nil
	} else {
		return "", errors.New("Unknown key")
	}
}

func (self *PgsqlPlugin) IsValid(cmd Command) bool {
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

func (self *PgsqlPlugin) Process(cmd Command) error {
	log.Printf("%v", self)
	log.Printf("%v", cmd)

	db, err := gorm.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", self.Username, "postgres", self.Password))
	if err != nil {
		return err
	}

	db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", self.cmd_username, self.cmd_password))
	db.Exec(fmt.Sprintf("CREATE DATABASE %s", self.cmd_database))
	db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s to %s", self.cmd_database, self.cmd_username))

	return nil
}
