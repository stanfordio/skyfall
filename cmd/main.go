package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
    "log"
	"os"
)

func main() {
	run(os.Args)
}

func run(args []string) {
	app := &cli.App{
		Name:    "skyfall",
		Usage:   "A simple CLI for Bluesky data ingest",
		Version: "prerelease",
		Commands: []*cli.Command{
			{
				Name:    "stream",
				Aliases: []string{"t"},
				Usage:   "Sip from the firehose",
				Action:  stream,
			},
		},
	}

    app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "ws-url",
			Usage: "full websocket path to the ATProto SubscribeRepos XRPC endpoint",
			Value: "wss://bsky.social/xrpc/com.atproto.sync.subscribeRepos",
		},
    }

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func stream(cctx *cli.Context) error {
    u, err := url.Parse(cctx.String("ws-url"))
	if err != nil {
		log.Fatalf("failed to parse ws-url: %+v", err)
	}


}
