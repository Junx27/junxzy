package cli

import "testing"

type fakeCommand struct {
	name string
}

func (f *fakeCommand) Name() string {
	return f.name
}

func (f *fakeCommand) Execute(args []string) {}

func TestRegisterCommand(t *testing.T) {
	commandRegistry = map[string]Command{}

	registered := &fakeCommand{name: "demo"}
	RegisterCommand(registered)

	got, ok := commandRegistry["demo"]
	if !ok {
		t.Fatalf("expected command to be registered")
	}

	if got != registered {
		t.Fatalf("registered command mismatch")
	}
}

func TestRegisterCommandsPopulatesDefaultCommands(t *testing.T) {
	commandRegistry = map[string]Command{}

	registerCommands()

	want := []string{"help", "clear", "exit", "make:module", "init", "simulate"}
	for _, name := range want {
		if _, ok := commandRegistry[name]; !ok {
			t.Fatalf("expected command %q to be registered", name)
		}
	}
}