package main

import (
	"errors"
	"fmt"
)

func getString(cmd Command, key string) (string, error) {
	if value, ok := cmd.Args[key]; ok {
		return fmt.Sprintf("%s", value), nil
	} else {
		return "", errors.New("Unknown key")
	}
}
