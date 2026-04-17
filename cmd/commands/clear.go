package commands

import "fmt"

type ClearCommand struct{}

func (c ClearCommand) Name() string {
	return "clear"
}

func (c ClearCommand) Execute(args []string) {
	fmt.Print("\033[H\033[2J")
}
