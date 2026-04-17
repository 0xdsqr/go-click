package click_test

import (
	"context"
	"flag"
	"fmt"

	click "github.com/0xdsqr/go-click"
)

func ExampleApp() {
	type rootFlags struct {
		Verbose bool
	}

	app := click.App[rootFlags]{
		Name: "demo",
		ConfigureRootFlags: func(fs *flag.FlagSet, root *rootFlags) {
			fs.BoolVar(&root.Verbose, "v", false, "enable verbose output")
		},
		Commands: []click.Command[rootFlags]{
			{
				Name: "hello",
				Run: func(_ context.Context, env click.Env[rootFlags], _ []string, _ []string) error {
					fmt.Fprintln(env.Stdout, "hello from go-click")
					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), []string{"hello"}); err != nil {
		panic(err)
	}

	// Output:
	// hello from go-click
}
