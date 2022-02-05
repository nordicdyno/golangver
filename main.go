package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/nordicdyno/golangver/golang"
)

func main() {
	newApp().runE(os.Args)
}

const goBinPathDefault = "/usr/local/bin/go"

type app struct {
	app       *cli.App
	goBinPath string
	verbose   bool
}

func newApp() *app {
	a := &app{}
	a.app = &cli.App{
		Name:  "golangver",
		Usage: "golang binaries version manager",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Usage:       "verbose (debug) output",
				Aliases:     []string{"v"},
				Destination: &a.verbose,
			},
			&cli.StringFlag{
				Name:        "go-bin",
				Usage:       "symlink to Go binary",
				Value:       goBinPathDefault,
				Destination: &a.goBinPath,
			},
		},
	}
	a.addFlags()
	return a
}

func (a *app) run(args []string) error {
	return a.app.Run(args)
}

func (a *app) runE(args []string) {
	if err := a.run(args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func (a *app) addFlags() {
	var forceInstall bool
	cInstall := &cli.Command{
		Name:  "get",
		Usage: "fetch version from https://go.dev/dl/",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "force",
				Aliases:     []string{"f"},
				Usage:       "remove if distro exists locally before fecth",
				Destination: &forceInstall,
			},
		},
		Action: func(cliCtx *cli.Context) error {
			args := cliCtx.Args()
			version := args.Get(0)
			if version == "" {
				return fmt.Errorf("version is not provided")
			}
			if version[0] == 'v' {
				version = version[1:]
			}

			if err := golang.Install(version, forceInstall); err != nil {
				return err
			}
			fmt.Println()
			// fmt.Println("Current version is")
			fmt.Printf("Now you can use Go %s:\n", version)
			fmt.Println("* with command:", "v-service go use", version)
			fmt.Println(" OR")
			fmt.Printf("* set Go path: export PATH=%s:$PATH\n", mustGoBinByVersion(version))
			return nil
		},
	}

	var listOpts golang.ListOpts
	cList := &cli.Command{
		Name:  "list",
		Usage: "show available versions",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "remotes",
				Usage:       "show available Go versions >= 1.13 (only latest minors if -a flag is not set)",
				Aliases:     []string{"r"},
				Destination: &listOpts.ShowRemotes,
			},
			&cli.BoolFlag{
				Name:        "all",
				Usage:       "show all versions (>= 1.13)",
				Aliases:     []string{"a"},
				Destination: &listOpts.ShowAllRemotes,
			},
			&cli.BoolFlag{
				Name:        "outdated",
				Usage:       "show all versions, even outdated (<1.13)",
				Aliases:     []string{"o"},
				Destination: &listOpts.ShowOutdated,
			},
		},
		Action: func(cliCtx *cli.Context) error {
			return golang.List(a.goBinPath, listOpts)
		},
	}

	cUse := &cli.Command{
		Name:  "use",
		Usage: "use provided version",
		Action: func(cliCtx *cli.Context) error {
			args := cliCtx.Args()
			version := args.Get(0)
			if version == "" {
				return fmt.Errorf("go version should be provided")
			}

			if version[0] == filepath.Separator {
				return golang.UseBinary(a.goBinPath, version)
			}

			if version[0] == 'v' {
				version = version[1:]
			}
			err := golang.UseVersion(a.goBinPath, version)
			if err != nil {
				return fmt.Errorf("switch to %s is failed: %w", version, err)
			}
			return nil
		},
	}

	a.app.Commands = append(a.app.Commands, cInstall, cList, cUse)
}

func mustGoBinByVersion(version string) string {
	goBin, err := golang.BinaryPath(version)
	if err != nil {
		panic(err)
	}
	return goBin
}
