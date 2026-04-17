package click

import "fmt"

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
