package click

import (
	"context"
	"flag"
	"io"
	"os"
)

type App[T any] struct {
	// Name is used when constructing the root flag set.
	Name string
	// ConfigureRootFlags registers flags shared across the whole command tree.
	ConfigureRootFlags func(*flag.FlagSet, *T)
	// Commands are the top-level commands available in the app.
	Commands []Command[T]
	// Stdout overrides the default command output writer.
	Stdout io.Writer
	// Stderr overrides the default error output writer.
	Stderr io.Writer
}

// Run parses root flags and dispatches the matching command.
func (a App[T]) Run(ctx context.Context, args []string) error {
	var root T

	fs := flag.NewFlagSet(a.Name, flag.ContinueOnError)
	fs.SetOutput(a.stderr())

	if a.ConfigureRootFlags != nil {
		a.ConfigureRootFlags(fs, &root)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	env := Env[T]{
		Stdout: a.stdout(),
		Stderr: a.stderr(),
		Root:   root,
	}

	return runCommands(ctx, env, a.Commands, fs.Args())
}

func (a App[T]) stdout() io.Writer {
	if a.Stdout != nil {
		return a.Stdout
	}
	return os.Stdout
}

func (a App[T]) stderr() io.Writer {
	if a.Stderr != nil {
		return a.Stderr
	}
	return os.Stderr
}
