package commands

import "fmt"

type HelpCommand struct{}

func (h HelpCommand) Name() string {
	return "help"
}

func (h HelpCommand) Execute(args []string) {
	fmt.Println("============ List command ============")
	fmt.Println("-- help")
	fmt.Println("-- hello [name]")
	fmt.Println("-- clear")
	fmt.Println("-- exit")
	fmt.Println("-- make:module [name]")
}
