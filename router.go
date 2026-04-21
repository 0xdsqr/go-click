package click

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"
)

func runCommands[T any](ctx context.Context, env Env[T], commands []Command[T], args []string) error {
	normal, pass := splitPassthrough(args)
	if len(normal) == 0 {
		return printHelp(env, commands)
	}

	name := normal[0]
	for _, cmd := range commands {
		if cmd.Name == name {
			return runCommand(ctx, env, nil, cmd, normal[1:], pass)
		}
	}

	return fmt.Errorf("unknown command %q", name)
}

func runCommand[T any](ctx context.Context, env Env[T], path []string, cmd Command[T], args []string, pass []string) error {
	fs, showHelp := newCommandFlagSet(commandPath(path, cmd.Name), env.Stderr)
	if cmd.ConfigureFlags != nil {
		cmd.ConfigureFlags(fs)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	fullPath := append(path, cmd.Name)

	if showHelp.value {
		return printCommandHelp(env, fullPath, cmd, fs)
	}

	if len(cmd.Commands) > 0 && len(remaining) > 0 {
		name := remaining[0]
		for _, child := range cmd.Commands {
			if child.Name == name {
				return runCommand(ctx, env, fullPath, child, remaining[1:], pass)
			}
		}

		if cmd.Run == nil {
			return fmt.Errorf("unknown command %q", strings.Join(append(fullPath, name), " "))
		}
	}

	if len(pass) > 0 && !cmd.Passthrough {
		return fmt.Errorf("command %q does not accept passthrough args", cmd.Name)
	}
	if cmd.Run == nil {
		if len(cmd.Commands) > 0 {
			return printCommandHelp(env, fullPath, cmd, fs)
		}
		return fmt.Errorf("command %q requires a subcommand", cmd.Name)
	}

	return cmd.Run(ctx, env, remaining, pass)
}

func splitPassthrough(args []string) ([]string, []string) {
	for i, arg := range args {
		if arg == "--" {
			return args[:i], args[i+1:]
		}
	}
	return args, nil
}

type helpFlag struct {
	value bool
}

func newCommandFlagSet(name string, stderr io.Writer) (*flag.FlagSet, *helpFlag) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	if stderr != nil {
		fs.SetOutput(stderr)
	}

	showHelp := &helpFlag{}
	fs.BoolVar(&showHelp.value, "h", false, "show help")
	fs.BoolVar(&showHelp.value, "help", false, "show help")

	return fs, showHelp
}

func commandPath(path []string, name string) string {
	return strings.Join(append(path, name), " ")
}
