package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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

func (self *PgsqlPlugin) IsValid(cmd Command) bool {
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

func (self *PgsqlPlugin) Create(cmd Command) error {
	db, err := gorm.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", self.Username, "postgres", self.Password))
	if err != nil {
		return err
	}

	db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", self.cmd_username, self.cmd_password))
	db.Exec(fmt.Sprintf("CREATE DATABASE %s", self.cmd_database))
	db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s to %s", self.cmd_database, self.cmd_username))

	return nil
}

func (self *PgsqlPlugin) Delete(cmd Command) error {
	db, err := gorm.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", self.Username, "postgres", self.Password))
	if err != nil {
		return err
	}

	db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", self.cmd_database))

	return nil
}

func (self *PgsqlPlugin) Update(cmd Command) error {
	db, err := gorm.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", self.Username, "postgres", self.Password))
	if err != nil {
		return err
	}

	db.Exec(fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s'", self.cmd_username, self.cmd_password))

	return nil
}
