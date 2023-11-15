package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/stanfordio/skyfall/pkg/hydrator"
	stream "github.com/stanfordio/skyfall/pkg/stream"
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
						Value: 32,
					},
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "host-domain",
			Usage: "host domain for the xrpc client",
			Value: "bsky.network",
		},
		&cli.Int64Flag{
			Name:  "cache-size",
			Usage: "maximum size of the cache, in bytes",
			Value: 1 << 32,
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

	u, err := url.Parse(fmt.Sprintf("wss://%s/xrpc/com.atproto.sync.subscribeRepos", cctx.String("host-domain")))
	if err != nil {
		log.Fatalf("failed to parse ws-url: %+v", err)
	}

	hydrator, err := hydrator.MakeHydrator(cctx.Context, cctx.Int64("cache-size"), cctx.String("host-domain"))
	if err != nil {
		log.Fatalf("failed to create hydrator: %+v", err)
	}

	s := stream.Stream{
		SocketURL: u,
		Output:    make(chan []byte),
		Hydrator:  hydrator,
	}

	go func() {
		err = s.BeginStreaming(ctx, cctx.Int("worker-count"))
		log.Fatalf("streaming ended unexpectedly: %+v", err)
		cancel()
	}()

	go func() {
		for {
			e := <-s.Output
			fmt.Println(string(e))
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
