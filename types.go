package main

type PluginInterface interface {
	Process(cmd Command) error
	IsValid(cmd Command) bool
	GetType() string
}

type Command struct {
	Plugin string
	Args   map[string]interface{}
}
