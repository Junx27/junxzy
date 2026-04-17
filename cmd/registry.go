package cmd

type Command interface {
	Name() string
	Execute(args []string)
}

var commandRegistry = map[string]Command{}

func RegisterCommand(cmd Command) {
	commandRegistry[cmd.Name()] = cmd
}
