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
	}

	var (
		gotRoot rootFlags
		gotArgs []string
		gotPass []string
	)

	app := App[rootFlags]{
		Name: "demo",
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
	if strings.Join(gotArgs, ",") != "one" {
		t.Fatalf("unexpected args: %v", gotArgs)
	}
	if strings.Join(gotPass, ",") != "two,three" {
		t.Fatalf("unexpected passthrough args: %v", gotPass)
	}
}
