package click

import (
	"flag"
	"fmt"
	"strings"
)

func printHelp[T any](env Env[T], commands []Command[T]) error {
	if len(commands) == 0 {
		_, err := fmt.Fprintln(env.Stdout, "No commands available.")
		return err
	}

	if _, err := fmt.Fprintln(env.Stdout, "Available commands:"); err != nil {
		return err
	}

	for _, cmd := range commands {
		if cmd.Description == "" {
			if _, err := fmt.Fprintf(env.Stdout, "  %s\n", cmd.Name); err != nil {
				return err
			}
			continue
		}

		if _, err := fmt.Fprintf(env.Stdout, "  %s\t%s\n", cmd.Name, cmd.Description); err != nil {
			return err
		}
	}

	return nil
}

func printCommandHelp[T any](env Env[T], path []string, cmd Command[T], fs *flag.FlagSet) error {
	usage := strings.Join(path, " ")
	if cmd.Usage != "" {
		usage = cmd.Usage
	}

	if _, err := fmt.Fprintf(env.Stdout, "Usage: %s\n", usage); err != nil {
		return err
	}
	if cmd.Description != "" {
		if _, err := fmt.Fprintf(env.Stdout, "\n%s\n", cmd.Description); err != nil {
			return err
		}
	}

	if fs != nil && hasVisibleFlags(fs) {
		if _, err := fmt.Fprintln(env.Stdout, "\nFlags:"); err != nil {
			return err
		}
		fs.SetOutput(env.Stdout)
		fs.PrintDefaults()
	}

	if len(cmd.Commands) == 0 {
		return nil
	}

	if _, err := fmt.Fprintln(env.Stdout, "\nSubcommands:"); err != nil {
		return err
	}

	for _, child := range cmd.Commands {
		if child.Description == "" {
			if _, err := fmt.Fprintf(env.Stdout, "  %s\n", child.Name); err != nil {
				return err
			}
			continue
		}

		if _, err := fmt.Fprintf(env.Stdout, "  %s\t%s\n", child.Name, child.Description); err != nil {
			return err
		}
	}

	return nil
}

func hasVisibleFlags(fs *flag.FlagSet) bool {
	hasFlags := false
	fs.VisitAll(func(*flag.Flag) {
		hasFlags = true
	})
	return hasFlags
}
