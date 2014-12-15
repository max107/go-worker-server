package main

import (
	"../go-openvz"

	"fmt"
	"strconv"
)

type OpenvzPlugin struct {
	cmd_ctid  int
	cmd_wait  bool
	cmd_force bool
}

func (self *OpenvzPlugin) GetType() string {
	return "openvz"
}

func (self *OpenvzPlugin) IsValid(cmd Command) bool {
	ctid, _ := getString(cmd, "ctid")
	if len(ctid) != 0 {
		self.cmd_ctid, _ = strconv.Atoi(ctid)
	}

	return true
}

func (self *OpenvzPlugin) Create(cmd Command) error {
	var createArgs = make(map[string]bool)
	createArgs["layout"] = true
	createArgs["ostemplate"] = true
	createArgs["hostname"] = true
	createArgs["ipadd"] = true
	createArgs["diskspace"] = true

	var createParams = make(map[string]string)
	var params = make(map[string]string)

	for key, value := range cmd.Args {
		if key == "ctid" {
			continue
		}

		if _, ok := createArgs[key]; ok {
			createParams[key] = fmt.Sprintf("%s", value)
		} else {
			params[key] = fmt.Sprintf("%s", value)
		}
	}

	container := openvz.Container{Ctid: self.cmd_ctid}
	err := container.CreateFromMap(createParams, params)
	if err != nil {
		return err
	}

	return container.Start(true)
}

func (self *OpenvzPlugin) Update(cmd Command) error {
	var params = make(map[string]string)

	container := openvz.Container{Ctid: self.cmd_ctid}

	for key, value := range cmd.Args {
		if key == "ctid" {
			continue
		}

		params[key] = fmt.Sprintf("%s", value)
	}

	return container.SetFromMap(params)
}

func (self *OpenvzPlugin) Delete(cmd Command) error {
	container := openvz.Container{Ctid: self.cmd_ctid}
	return container.Delete(true)
}
