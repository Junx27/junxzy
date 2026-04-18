package commands

import "fmt"

type HelloCommand struct{}

func (h HelloCommand) Name() string {
	return "hello"
}

func (h HelloCommand) Execute(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: hello <name>")
		return
	}

	fmt.Printf("Hello, %s!\n", args[0])
}