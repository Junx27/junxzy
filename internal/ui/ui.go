package ui

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	cyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	green  = color.New(color.FgGreen, color.Bold).SprintFunc()
	red    = color.New(color.FgRed, color.Bold).SprintFunc()
	yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
)

// Step: proses berjalan
func Step(message string) {
	fmt.Println(cyan("➜ ") + message)
}

// Success: berhasil
func Success(message string) {
	fmt.Println(green("✔ ") + message)
}

// Error: gagal
func Error(message string) {
	fmt.Println(red("✖ ") + message)
}

// Info: tambahan info
func Info(message string) {
	fmt.Println(yellow("➤ ") + message)
}
