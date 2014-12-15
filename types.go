package main

type PluginInterface interface {
	Create(cmd Command) error
	Update(cmd Command) error
	Delete(cmd Command) error
	IsValid(cmd Command) bool
	GetType() string
}

type Command struct {
	Plugin string
	Action string
	Args   map[string]interface{}
}
