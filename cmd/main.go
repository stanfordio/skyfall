package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	stream "github.com/stanfordio/skyfall/pkg"
	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
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
				Action:  streamCmd,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "worker-count",
						Usage: "number of workers to scale to",
						Value: 4,
					},
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "ws-url",
			Usage: "full websocket path to the ATProto SubscribeRepos XRPC endpoint",
			Value: "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func streamCmd(cctx *cli.Context) error {
	ctx := cctx.Context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	u, err := url.Parse(cctx.String("ws-url"))
	if err != nil {
		log.Fatalf("failed to parse ws-url: %+v", err)
	}

	s := stream.Stream{
		SocketURL: u,
		Events:    make(chan string),
	}

	go func() {
		err = s.BeginStreaming(ctx, cctx.Int("worker-count"))
		log.Fatalf("streaming ended unexpectedly: %+v", err)
		cancel()
	}()

	go func() {
		for {
			e := <-s.Events
			fmt.Printf(e)
		}
	}()

	select {
	case <-signals:
		cancel()
		log.Infof("shutting down on signal")
	case <-ctx.Done():
		log.Infof("shutting down on context done")
	}

	return nil
}
