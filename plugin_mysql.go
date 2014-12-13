package main

import (
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"

	"errors"
	"fmt"
	"log"
)

type MysqlPlugin struct {
	Hostname string
	Port     int
	Username string
	Password string

	cmd_username string
	cmd_password string
	cmd_database string
}

func (self *MysqlPlugin) GetType() string {
	return "mysql"
}

func (self *MysqlPlugin) getString(cmd Command, key string) (string, error) {
	if value, ok := cmd.Args[key]; ok {
		return fmt.Sprintf("%s", value), nil
	} else {
		return "", errors.New("Unknown key")
	}
}

func (self *MysqlPlugin) IsValid(cmd Command) bool {
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

func (self *MysqlPlugin) Process(cmd Command) error {
	log.Printf("%v", self)
	log.Printf("%v", cmd)

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True", self.Username, self.Password, "mysql"))
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE DATABASE ?", self.cmd_database)
	if err != nil {
		return err
	}

	_, err = db.Exec("GRANT ALL PRIVILEGES ON ?.* TO ?@localhost IDENTIFIED BY '?'", self.cmd_database, self.cmd_username, self.cmd_password)
	if err != nil {
		return err
	}

	return nil
}
