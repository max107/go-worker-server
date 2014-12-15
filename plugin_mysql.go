package main

import (
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
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

func (self *MysqlPlugin) IsValid(cmd Command) bool {
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

func (self *MysqlPlugin) Create(cmd Command) error {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True", self.Username, self.Password, "mysql"))
	if err != nil {
		return err
	}

	db.Exec(fmt.Sprintf("CREATE DATABASE %s", self.cmd_database))
	db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO %s@localhost IDENTIFIED BY '%s'", self.cmd_database, self.cmd_username, self.cmd_password))

	return nil
}

func (self *MysqlPlugin) Delete(cmd Command) error {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True", self.Username, self.Password, "mysql"))
	if err != nil {
		return err
	}

	db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", self.cmd_database))

	return nil
}

func (self *MysqlPlugin) Update(cmd Command) error {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True", self.Username, self.Password, "mysql"))
	if err != nil {
		return err
	}

	db.Exec(fmt.Sprintf("SET PASSWORD FOR '%s'@localhost = PASSWORD('%s');", self.cmd_username, self.cmd_password))

	return nil
}
