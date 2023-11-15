package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/bluesky-social/indigo/xrpc"
	"github.com/stanfordio/skyfall/pkg/auth"
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
		&cli.StringFlag{
			Name:  "handle",
			Usage: "handle to authenticate with, e.g., @miles.land or @det.bsky.social",
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "password to authenticate with",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func authenticate(cctx *cli.Context) (*xrpc.AuthInfo, error) {
	authenticator, err := auth.MakeAuthenticator(cctx.Context)

	if err != nil {
		log.Fatalf("failed to create authenticator: %+v", err)
		return nil, err
	}

	authInfo, err := authenticator.Authenticate(cctx.String("handle"), cctx.String("password"))
	if err != nil {
		log.Fatalf("failed to authenticate: %+v", err)
		return nil, err
	}

	return authInfo, nil
}

func streamCmd(cctx *cli.Context) error {
	ctx := cctx.Context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Create a client
	authInfo, err := authenticate(cctx)
	if err != nil {
		log.Fatalf("failed to authenticate: %+v", err)
		return err
	}

	u, err := url.Parse(fmt.Sprintf("wss://%s/xrpc/com.atproto.sync.subscribeRepos", cctx.String("host-domain")))
	if err != nil {
		log.Fatalf("failed to parse ws-url: %+v", err)
		return err
	}

	hydrator, err := hydrator.MakeHydrator(cctx.Context, cctx.Int64("cache-size"), cctx.String("host-domain"), authInfo)
	if err != nil {
		log.Fatalf("failed to create hydrator: %+v", err)
		return err
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
