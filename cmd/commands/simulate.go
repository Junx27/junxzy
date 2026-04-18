package commands

import (
	"time"

	"github.com/Junx27/junxzy/internal/ui"
)

type SimulateCommand struct{}

var simulateSleep = time.Sleep

func (s SimulateCommand) Name() string {
	return "simulate"
}

func (s SimulateCommand) Execute(args []string) {

	ui.RunStep("Creating project folder", func() {
		simulateSleep(1 * time.Second)
	})

	ui.RunStep("Generate services", func() {
		simulateSleep(2 * time.Second)
	})

	ui.RunStep("Setup gateway", func() {
		simulateSleep(1 * time.Second)
	})

	ui.RunStep("Generate docker-compose", func() {
		simulateSleep(2 * time.Second)
	})

	ui.RunStep("Finalizing project", func() {
		simulateSleep(1 * time.Second)
	})

	ui.Success("Simulation completed 🚀")
}
