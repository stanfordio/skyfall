package main

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/bluesky-social/indigo/xrpc"
	"github.com/stanfordio/skyfall/pkg/auth"
	"github.com/stanfordio/skyfall/pkg/hydrator"
	stream "github.com/stanfordio/skyfall/pkg/stream"
	"github.com/stanfordio/skyfall/pkg/utils"
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
					&cli.StringFlag{
						Name:  "output-file",
						Usage: "file to write output to (if specified, will attempt to backfill from the most recent event in the file)",
						Value: "output.jsonl",
					},
					&cli.Int64Flag{
						Name:  "backfill-seq",
						Usage: "seq to backfill from (if specified, will override the seqno extracted from the output file)",
						Value: 0,
					},
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		&cli.Int64Flag{
			Name:  "cache-size",
			Usage: "maximum size of the cache, in bytes",
			Value: 1 << 32,
		},
		&cli.StringFlag{
			Name:  "handle",
			Usage: "handle to authenticate with, e.g., miles.land or det.bsky.social",
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
		log.Fatalf("Failed to create authenticator: %+v", err)
		return nil, err
	}

	authInfo, err := authenticator.Authenticate(cctx.String("handle"), cctx.String("password"))
	if err != nil {
		log.Fatalf("Failed to authenticate: %+v", err)
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
		log.Fatalf("Failed to authenticate: %+v", err)
		return err
	}

	u, err := url.Parse("wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos")
	if err != nil {
		log.Fatalf("Failed to parse ws-url: %+v", err)
		return err
	}

	hydrator, err := hydrator.MakeHydrator(cctx.Context, cctx.Int64("cache-size"), authInfo)
	if err != nil {
		log.Fatalf("Failed to create hydrator: %+v", err)
		return err
	}

	var lastSeq int64 = cctx.Int64("backfill-seq")

	if lastSeq == 0 {
		log.Infof("No backfill seq specified, so attempting to backfill from the last line of the output file")

		lastLine, err := utils.GetLastLine(cctx.String("output-file"))
		if err != nil {
			log.Warnf("Unable to read last line of output file for backfill: %+v", err)
			log.Warnf("Continuing without backfill...")
		}

		// Try to parse the last line as JSON, then pull out the last "seq" field
		lastData := make(map[string]interface{})
		err = json.Unmarshal([]byte(lastLine), &lastData)
		if err != nil {
			log.Warnf("Unable to parse last line as JSON: %+v", err)
			log.Warnf("Continuing without backfill...")
		}
		lastSeqFloat := lastData["_Seq"]
		if lastSeqFloat == nil {
			log.Warnf("Unable to find seq in last line of output file")
			log.Warnf("Continuing without backfill...")
		} else {
			lastSeq = int64(lastSeqFloat.(float64))
			log.Infof("Backfilling from inferred seq (from output file): %d", lastSeq)
		}
	} else {
		log.Infof("Backfilling from provided seq: %d", lastSeq)
	}

	// Open the file
	f, err := os.OpenFile(cctx.String("output-file"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %+v", err)
		return err
	}
	defer f.Close()

	s := stream.Stream{
		SocketURL:   u,
		Output:      make(chan []byte),
		Hydrator:    hydrator,
		BackfillSeq: lastSeq,
	}

	go func() {
		err = s.BeginStreaming(ctx, cctx.Int("worker-count"))
		log.Fatalf("Streaming ended unexpectedly: %+v", err)
		cancel()
	}()

	go func() {
		for {
			e := <-s.Output
			if _, err := f.Write(append(e, byte('\n'))); err != nil {
				log.Errorf("Failed to write output: %+v", err)
				cancel()
			}
		}
	}()

	select {
	case <-signals:
		cancel()
		log.Infof("Shutting down on signal")
	case <-ctx.Done():
		log.Infof("Shutting down on context done")
	}

	return nil
}
