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

	ui.Start("Membuat folder project...")
	time.Sleep(1 * time.Second)
	ui.StopSuccess("Folder project dibuat")

	ui.Start("Generate services...")
	time.Sleep(2 * time.Second)
	ui.StopSuccess("Services berhasil dibuat")

	ui.Start("Setup gateway...")
	time.Sleep(1 * time.Second)
	ui.StopSuccess("Gateway siap")

	ui.Start("Generate docker-compose...")
	time.Sleep(2 * time.Second)
	ui.StopSuccess("docker-compose dibuat")

	ui.Start("Finalizing project...")
	time.Sleep(1 * time.Second)
	ui.StopSuccess("Project siap digunakan")

	ui.Success("Simulasi selesai 🚀")
}
