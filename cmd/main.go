package main

import (
	"os"
	"fmt"
	"log"
	"github.com/urfave/cli/v2"
)

var log = logging.Logger("bigsky")

func main() {
	run(os.Args)
}

func run(args []string) {
	app := &cli.App{
		Name: "skyfall",
		Usage: "A simple CLI for Bluesky data ingest",
		Version: "prerelease",
		Commands: []*cli.Command{
            {
                Name:    "test",
                Aliases: []string{"t"},
                Usage:   "test command",
                Action: func(cCtx *cli.Context) error {
                    fmt.Println("test")
                    return nil
                },
            }
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
