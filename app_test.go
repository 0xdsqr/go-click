package click

import (
	"bytes"
	"context"
	"flag"
	"strings"
	"testing"
)

func TestAppRunPrintsHelpWhenNoCommandIsGiven(t *testing.T) {
	var stdout bytes.Buffer

	app := App[struct{}]{
		Name:   "demo",
		Stdout: &stdout,
		Commands: []Command[struct{}]{
			{Name: "hello", Description: "print a greeting"},
		},
	}

	if err := app.Run(context.Background(), nil); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "Available commands:") {
		t.Fatalf("help output missing heading: %q", got)
	}
	if !strings.Contains(got, "hello") {
		t.Fatalf("help output missing command: %q", got)
	}
}

func TestAppRunPassesRootFlagsAndPassthrough(t *testing.T) {
	type rootFlags struct {
		Verbose bool
		Prefix  string
	}

	var (
		gotRoot rootFlags
		gotArgs []string
		gotPass []string
	)

	app := App[rootFlags]{
		Name: "demo",
		ConfigureRoot: func(root *rootFlags) {
			root.Prefix = "demo"
		},
		ConfigureRootFlags: func(fs *flag.FlagSet, root *rootFlags) {
			fs.BoolVar(&root.Verbose, "v", false, "enable verbose output")
		},
		Commands: []Command[rootFlags]{
			{
				Name:        "hello",
				Passthrough: true,
				Run: func(_ context.Context, env Env[rootFlags], args []string, pass []string) error {
					gotRoot = env.Root
					gotArgs = append([]string(nil), args...)
					gotPass = append([]string(nil), pass...)
					return nil
				},
			},
		},
	}

	err := app.Run(context.Background(), []string{"-v", "hello", "one", "--", "two", "three"})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !gotRoot.Verbose {
		t.Fatalf("expected verbose root flag to be true")
	}
	if gotRoot.Prefix != "demo" {
		t.Fatalf("expected prefix to be initialized, got %q", gotRoot.Prefix)
	}
	if strings.Join(gotArgs, ",") != "one" {
		t.Fatalf("unexpected args: %v", gotArgs)
	}
	if strings.Join(gotPass, ",") != "two,three" {
		t.Fatalf("unexpected passthrough args: %v", gotPass)
	}
}

func TestAppRunPrintsRootHelp(t *testing.T) {
	var stdout bytes.Buffer

	app := App[struct{}]{
		Name:   "demo",
		Stdout: &stdout,
		ConfigureRootFlags: func(fs *flag.FlagSet, _ *struct{}) {
			fs.Bool("v", false, "enable verbose output")
		},
		Commands: []Command[struct{}]{
			{Name: "hello", Description: "print a greeting"},
		},
	}

	if err := app.Run(context.Background(), []string{"--help"}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "Usage: demo [flags] <command>") {
		t.Fatalf("help output missing usage: %q", got)
	}
	if !strings.Contains(got, "Commands:") {
		t.Fatalf("help output missing commands heading: %q", got)
	}
	if !strings.Contains(got, "-v") {
		t.Fatalf("help output missing root flag: %q", got)
	}
}
