package commands

import (
	"fmt"
	"os"
)

type ExitCommand struct{}

func (e ExitCommand) Name() string {
	return "exit"
}

func (e ExitCommand) Execute(args []string) {
	fmt.Println("Bye! 👋")
	os.Exit(0)
}
