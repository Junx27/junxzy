package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "junxzy",
	Short: "Happy coding with Junxzy CLI! 🚀",
	Run: func(cmd *cobra.Command, args []string) {
		printBanner()
		startREPL()
	},
}

func printBanner() {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	white := color.New(color.FgWhite, color.Bold).SprintFunc()

	fmt.Println()
	fmt.Println(cyan("   ██╗") + magenta("  ██╗   ██╗  ") + yellow("███╗   ██╗") + green("  ██╗  ██╗  ") + red("███████╗") + blue("  ██╗   ██╗"))
	fmt.Println(cyan("   ██║") + magenta("  ██║   ██║  ") + yellow("████╗  ██║") + green("  ╚██╗██╔╝  ") + red("╚════██║") + blue("  ╚██╗ ██╔╝"))
	fmt.Println(cyan("   ██║") + magenta("  ██║   ██║  ") + yellow("██╔██╗ ██║") + green("   ╚███╔╝   ") + red("    ██╔╝") + blue("   ╚████╔╝ "))
	fmt.Println(cyan("██ ██║") + magenta("  ██║   ██║  ") + yellow("██║╚██╗██║") + green("   ██╔██╗   ") + red("   ██╔╝ ") + blue("    ╚██╔╝  "))
	fmt.Println(cyan("╚████╔╝") + magenta(" ╚██████╔╝  ") + yellow("██║ ╚████║") + green("  ██╔╝ ██╗  ") + red("  ██████╗") + blue("    ██║   "))
	fmt.Println(cyan(" ╚═══╝ ") + magenta("  ╚═════╝   ") + yellow("╚═╝  ╚═══╝") + green("  ╚═╝  ╚═╝  ") + red("  ╚═════╝") + blue("    ╚═╝   "))
	fmt.Println()

	fmt.Println(cyan("  ╔══════════════════════════════════════╗"))
	fmt.Println(cyan("  ║") + white("   Welcome to Junxzy CLI! 🚀          ") + cyan("║"))
	fmt.Println(cyan("  ║") + white("   Happy coding with Golang 🔥        ") + cyan("║"))
	fmt.Println(cyan("  ╚══════════════════════════════════════╝"))
	fmt.Println()
	fmt.Println(white("  Junxzy CLI is your personal developer assistant"))
	fmt.Println(white("  for generating modules, managing projects, and speeding up your workflow."))
	fmt.Println(yellow("  Type ") + green("'help'") + yellow(" to see all commands."))
	fmt.Println(yellow("  Type ") + red("'exit'") + yellow(" to exit."))
	fmt.Println()
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
