package cmd

import (
	"bufio"
	"fmt"

	"os"
	"strings"

	"github.com/Junx27/junxzy/cmd/commands"
	"github.com/fatih/color"
)

func startREPL() {
	// register semua command
	registerCommands()

	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(cyan("junxzy") + " » ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		cmdName := parts[0]
		args := parts[1:]

		cmd, exists := commandRegistry[cmdName]
		if !exists {
			fmt.Println(red("Command tidak dikenal:"), cmdName)
			continue
		}

		cmd.Execute(args)
	}
}

func registerCommands() {
	RegisterCommand(commands.HelpCommand{})
	RegisterCommand(commands.ClearCommand{})
	RegisterCommand(commands.ExitCommand{})
	RegisterCommand(commands.MakeModuleCommand{})
}
