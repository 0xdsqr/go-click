package click

import (
	"context"
	"fmt"
)

func runCommands[T any](ctx context.Context, env Env[T], commands []Command[T], args []string) error {
	if len(args) == 0 {
		return printHelp(env, commands)
	}

	name := args[0]
	for _, cmd := range commands {
		if cmd.Name == name {
			return runCommand(ctx, env, cmd, args[1:])
		}
	}

	return fmt.Errorf("unknown command %q", name)
}

func runCommand[T any](ctx context.Context, env Env[T], cmd Command[T], args []string) error {
	if len(cmd.Commands) > 0 && len(args) > 0 && args[0] != "--" {
		for _, child := range cmd.Commands {
			if child.Name == args[0] {
				return runCommand(ctx, env, child, args[1:])
			}
		}
	}

	normal, pass := splitPassthrough(args)

	if len(pass) > 0 && !cmd.Passthrough {
		return fmt.Errorf("command %q does not accept passthrough args", cmd.Name)
	}
	if cmd.Run == nil {
		return fmt.Errorf("command %q requires a subcommand", cmd.Name)
	}

	return cmd.Run(ctx, env, normal, pass)
}

func splitPassthrough(args []string) ([]string, []string) {
	for i, arg := range args {
		if arg == "--" {
			return args[:i], args[i+1:]
		}
	}
	return args, nil
}
