package ui

import "fmt"

func RunStep(message string, fn func()) {
	Start(message)
	defer func() {
		if r := recover(); r != nil {
			StopError(fmt.Sprintf("Error: %v", r))
			return
		}
	}()
	fn()
	StopSuccess(message + " selesai")
}
