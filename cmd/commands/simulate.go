package commands

import (
	"time"

	"github.com/Junx27/junxzy/internal/ui"
)

type SimulateCommand struct{}

func (s SimulateCommand) Name() string {
	return "simulate"
}

func (s SimulateCommand) Execute(args []string) {

	ui.RunStep("Membuat folder project", func() {
		time.Sleep(1 * time.Second)
	})

	ui.RunStep("Generate services", func() {
		time.Sleep(2 * time.Second)
	})

	ui.RunStep("Setup gateway", func() {
		time.Sleep(1 * time.Second)
	})

	ui.RunStep("Generate docker-compose", func() {
		time.Sleep(2 * time.Second)
	})

	ui.RunStep("Finalizing project", func() {
		time.Sleep(1 * time.Second)
	})

	ui.Success("Simulasi selesai 🚀")
}
