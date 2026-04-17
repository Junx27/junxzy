package ui

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var (
	green = color.New(color.FgGreen, color.Bold).SprintFunc()
	red   = color.New(color.FgRed, color.Bold).SprintFunc()
)

var s *spinner.Spinner

func Success(message string) {
	fmt.Println(green("✔ ") + message)
}

// Start loading animation
func Start(message string) {
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Start()
}

// Stop loading animation (success)
func StopSuccess(message string) {
	if s != nil {
		s.Stop()
	}
	fmt.Println(green("✔ ") + message)
}

// Stop loading animation (error)
func StopError(message string) {
	if s != nil {
		s.Stop()
	}
	fmt.Println(red("✖ ") + message)
}
