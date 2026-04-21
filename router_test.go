package click

import (
	"bytes"
	"context"
	"flag"
	"strings"
	"testing"
)

func TestRunCommandsReturnsUnknownCommand(t *testing.T) {
	err := runCommands(context.Background(), Env[struct{}]{}, nil, []string{"wat"})
	if err == nil {
		t.Fatal("expected an error for an unknown command")
	}
	if !strings.Contains(err.Error(), `unknown command "wat"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunCommandRejectsUnexpectedPassthrough(t *testing.T) {
	err := runCommand(
		context.Background(),
		Env[struct{}]{},
		nil,
		Command[struct{}]{
			Name: "hello",
			Run: func(context.Context, Env[struct{}], []string, []string) error {
				return nil
			},
		},
		nil,
		[]string{"tail"},
	)
	if err == nil {
		t.Fatal("expected passthrough validation error")
	}
	if !strings.Contains(err.Error(), "does not accept passthrough args") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrintHelpWithNoCommands(t *testing.T) {
	var stdout bytes.Buffer

	err := printHelp(Env[struct{}]{Stdout: &stdout}, nil)
	if err != nil {
		t.Fatalf("PrintHelp returned error: %v", err)
	}
	if !strings.Contains(stdout.String(), "No commands available.") {
		t.Fatalf("unexpected help output: %q", stdout.String())
	}
}

func TestRunCommandReturnsSubcommandErrorWhenRunIsNil(t *testing.T) {
	err := runCommand(context.Background(), Env[struct{}]{}, nil, Command[struct{}]{Name: "root"}, nil, nil)
	if err == nil {
		t.Fatal("expected subcommand error")
	}
	if !strings.Contains(err.Error(), `requires a subcommand`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunCommandParsesCommandFlagsBeforeRun(t *testing.T) {
	var (
		jsonOutput bool
		gotArgs    []string
	)

	err := runCommand(
		context.Background(),
		Env[struct{}]{},
		nil,
		Command[struct{}]{
			Name: "list",
			ConfigureFlags: func(fs *flag.FlagSet) {
				fs.BoolVar(&jsonOutput, "json", false, "output JSON")
			},
			Run: func(_ context.Context, _ Env[struct{}], args []string, _ []string) error {
				gotArgs = append([]string(nil), args...)
				return nil
			},
		},
		[]string{"--json", "projects"},
		nil,
	)
	if err != nil {
		t.Fatalf("runCommand returned error: %v", err)
	}
	if !jsonOutput {
		t.Fatal("expected command flag to be parsed")
	}
	if strings.Join(gotArgs, ",") != "projects" {
		t.Fatalf("unexpected args: %v", gotArgs)
	}
}

func TestRunCommandParsesNestedCommandFlags(t *testing.T) {
	var jsonOutput bool

	err := runCommand(
		context.Background(),
		Env[struct{}]{},
		nil,
		Command[struct{}]{
			Name: "project",
			Commands: []Command[struct{}]{
				{
					Name: "list",
					ConfigureFlags: func(fs *flag.FlagSet) {
						fs.BoolVar(&jsonOutput, "json", false, "output JSON")
					},
					Run: func(_ context.Context, _ Env[struct{}], _ []string, _ []string) error {
						return nil
					},
				},
			},
		},
		[]string{"list", "--json"},
		nil,
	)
	if err != nil {
		t.Fatalf("runCommand returned error: %v", err)
	}
	if !jsonOutput {
		t.Fatal("expected nested command flag to be parsed")
	}
}

func TestRunCommandPrintsCommandHelp(t *testing.T) {
	var stdout bytes.Buffer

	err := runCommand(
		context.Background(),
		Env[struct{}]{Stdout: &stdout},
		nil,
		Command[struct{}]{
			Name:        "project",
			Description: "work with projects",
			Usage:       "demo project",
			Commands: []Command[struct{}]{
				{Name: "list", Description: "list projects"},
			},
		},
		[]string{"--help"},
		nil,
	)
	if err != nil {
		t.Fatalf("runCommand returned error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "Usage: demo project") {
		t.Fatalf("help output missing usage: %q", got)
	}
	if !strings.Contains(got, "Subcommands:") {
		t.Fatalf("help output missing subcommands: %q", got)
	}
	if !strings.Contains(got, "-help") {
		t.Fatalf("help output missing flags: %q", got)
	}
}
