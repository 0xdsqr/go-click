package click

import (
	"context"
	"flag"
	"io"
)

// Env contains the values passed to a command handler at runtime.
type Env[T any] struct {
	Stdout io.Writer
	Stderr io.Writer
	Root   T
}

// Command describes one CLI command in the tree.
type Command[T any] struct {
	// Name is the token users type to select this command.
	Name string
	// Description is a short help string shown in command listings.
	Description string
	// Usage overrides the default command path shown in help output.
	Usage string
	// Commands are this command's nested subcommands.
	Commands []Command[T]
	// ConfigureFlags registers flags for this specific command.
	ConfigureFlags func(*flag.FlagSet)
	// Passthrough allows args after `--` to be passed separately to Run.
	Passthrough bool
	// Run executes the command with normal and passthrough arguments.
	Run func(context.Context, Env[T], []string, []string) error
}
