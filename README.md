# go-click

<p align="center">
  <a href="https://github.com/0xdsqr/go-click/actions/workflows/test.yml">
    <img alt="ci" src="https://github.com/0xdsqr/go-click/actions/workflows/test.yml/badge.svg?branch=master">
  </a>
</p>

tiny cli plumbing for lightweight or personal go projects. `go-click` is a tiny
helper for small projects. you should probably just copy it into your own repo,
but i'm lazy. if you want something more standard or more fully featured, use
[cobra](https://github.com/spf13/cobra).

## ⇁ TOC

* [What](#-what)
* [Installation](#-installation)
* [Getting Started](#-getting-started)
* [API](#-api)
* [Contributing](#-contributing)
* [MIT](#-mit)

## ⇁ What

you can really only do these things with it:

* shared root state
* root flags
* command flags
* nested commands
* passthrough args after `--`
* injected `stdout` / `stderr`

that's the whole idea.

## ⇁ Installation

```bash
go get github.com/0xdsqr/go-click
```

## ⇁ Getting Started

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	click "github.com/0xdsqr/go-click"
)

type root struct {
	verbose bool
}

func main() {
	app := click.App[root]{
		Name: "demo",
		ConfigureRoot: func(cfg *root) {
			cfg.verbose = false
		},
		ConfigureRootFlags: func(fs *flag.FlagSet, cfg *root) {
			fs.BoolVar(&cfg.verbose, "v", false, "enable verbose output")
		},
		Commands: []click.Command[root]{
			{
				Name: "hello",
				Run: func(_ context.Context, env click.Env[root], _ []string, _ []string) error {
					if env.Root.verbose {
						fmt.Fprintln(env.Stderr, "verbose mode enabled")
					}
					fmt.Fprintln(env.Stdout, "hello from go-click")
					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

## ⇁ API

docs:

* `App`
* `Command`
* `Env`
* `https://pkg.go.dev/github.com/0xdsqr/go-click`

<details>
<summary>basic command</summary>

```go
app := click.App[struct{}]{
	Name: "demo",
	Commands: []click.Command[struct{}]{
		{
			Name: "hello",
			Run: func(_ context.Context, env click.Env[struct{}], _ []string, _ []string) error {
				fmt.Fprintln(env.Stdout, "hello")
				return nil
			},
		},
	},
}
```

</details>

<details>
<summary>shared root state</summary>

```go
type root struct {
	host string
}

app := click.App[root]{
	Name: "demo",
	ConfigureRoot: func(cfg *root) {
		cfg.host = "local"
	},
	Commands: []click.Command[root]{
		{
			Name: "host",
			Run: func(_ context.Context, env click.Env[root], _ []string, _ []string) error {
				fmt.Fprintln(env.Stdout, env.Root.host)
				return nil
			},
		},
	},
}
```

</details>

<details>
<summary>root flags</summary>

```go
type root struct {
	verbose bool
}

app := click.App[root]{
	Name: "demo",
	ConfigureRootFlags: func(fs *flag.FlagSet, cfg *root) {
		fs.BoolVar(&cfg.verbose, "v", false, "enable verbose output")
	},
	Commands: []click.Command[root]{
		{
			Name: "hello",
			Run: func(_ context.Context, env click.Env[root], _ []string, _ []string) error {
				if env.Root.verbose {
					fmt.Fprintln(env.Stderr, "verbose mode enabled")
				}
				fmt.Fprintln(env.Stdout, "hello")
				return nil
			},
		},
	},
}
```

</details>

<details>
<summary>nested commands</summary>

```go
app := click.App[struct{}]{
	Name: "demo",
	Commands: []click.Command[struct{}]{
		{
			Name: "project",
			Commands: []click.Command[struct{}]{
				{
					Name: "list",
					Run: func(_ context.Context, env click.Env[struct{}], _ []string, _ []string) error {
						fmt.Fprintln(env.Stdout, "listing projects")
						return nil
					},
				},
			},
		},
	},
}
```

usage:

```bash
demo project list
```

</details>

<details>
<summary>command flags</summary>

```go
var formatJSON bool

app := click.App[struct{}]{
	Name: "demo",
	Commands: []click.Command[struct{}]{
		{
			Name: "project",
			Commands: []click.Command[struct{}]{
				{
					Name: "list",
					ConfigureFlags: func(fs *flag.FlagSet) {
						fs.BoolVar(&formatJSON, "json", false, "output JSON")
					},
					Run: func(_ context.Context, env click.Env[struct{}], _ []string, _ []string) error {
						if formatJSON {
							fmt.Fprintln(env.Stdout, `["a","b"]`)
							return nil
						}
						fmt.Fprintln(env.Stdout, "a\nb")
						return nil
					},
				},
			},
		},
	},
}
```

usage:

```bash
demo project list --json
```

</details>

<details>
<summary>passthrough args with --</summary>

```go
app := click.App[struct{}]{
	Name: "demo",
	Commands: []click.Command[struct{}]{
		{
			Name:        "exec",
			Passthrough: true,
			Run: func(_ context.Context, env click.Env[struct{}], args []string, pass []string) error {
				fmt.Fprintln(env.Stdout, "normal args:", args)
				fmt.Fprintln(env.Stdout, "passthrough args:", pass)
				return nil
			},
		},
	},
}
```

usage:

```bash
demo exec build target -- --watch --verbose
```

`args` will be `["build", "target"]`

`pass` will be `["--watch", "--verbose"]`

</details>

<details>
<summary>custom stdout / stderr</summary>

```go
var stdout bytes.Buffer
var stderr bytes.Buffer

app := click.App[struct{}]{
	Name:   "demo",
	Stdout: &stdout,
	Stderr: &stderr,
}
```

</details>

## ⇁ Contributing

some people have started using this for whatever reason. i am probably not
going to accept mrs.

clone it. use the flake. do whatever you want with it.

```bash
nix develop
nix flake check
```

## ⇁ MIT

do whatever you want with it.
