package click

import (
	"context"
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
	// Usage is reserved for future, more detailed help output.
	Usage string
	// Commands are this command's nested subcommands.
	Commands []Command[T]
	// Passthrough allows args after `--` to be passed separately to Run.
	Passthrough bool
	// Run executes the command with normal and passthrough arguments.
	Run func(context.Context, Env[T], []string, []string) error
}
