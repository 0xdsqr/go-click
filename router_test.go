package click

import (
	"bytes"
	"context"
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
		Command[struct{}]{
			Name: "hello",
			Run: func(context.Context, Env[struct{}], []string, []string) error {
				return nil
			},
		},
		[]string{"--", "tail"},
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
	err := runCommand(context.Background(), Env[struct{}]{}, Command[struct{}]{Name: "root"}, nil)
	if err == nil {
		t.Fatal("expected subcommand error")
	}
	if !strings.Contains(err.Error(), `requires a subcommand`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
