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
	"github.com/stanfordio/skyfall/pkg/output"
	repodump "github.com/stanfordio/skyfall/pkg/repodump"
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
					&cli.StringFlag{
						Name:  "output-file",
						Usage: "file to write output to (if specified, will attempt to backfill from the most recent event in the file)",
						Value: "output.jsonl",
					},
					&cli.BoolFlag{
						Name:  "stringify-full",
						Usage: "whether to stringify the full event in file output (if true, the JSON will be stringified; this is helpful when you want output to match what would be sent to BigQuery)",
						Value: false,
					},
					&cli.StringFlag{
						Name:  "output-bq-table",
						Usage: "name of a BigQuery table to output to in ID form (e.g., dgap_bsky.example_table)",
					},
					&cli.Int64Flag{
						Name:  "backfill-seq",
						Usage: "seq to backfill from (if specified, will override the seqno extracted from the output file/bigquery table)",
						Value: 0,
					},
					&cli.BoolFlag{
						Name:  "autorestart",
						Usage: "automatically restart the stream if it dies",
						Value: true,
					},
				},
			},
			{
				Name:    "repodump",
				Aliases: []string{"d"},
				Usage:   "Dump everyone's repos (as CAR) into a folder",
				Action:  repodumpCmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "output-folder",
						Usage: "folder to write repos to",
						Value: "output",
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

	outputChannel := make(chan map[string]interface{})

	output, err := output.NewOutput(cctx, outputChannel)
	if err != nil {
		log.Fatalf("Failed to create output: %+v", err)
		return err
	}

	// Setup the output
	err = output.Setup()
	if err != nil {
		log.Fatalf("Failed to setup output: %+v", err)
		return err
	}

	var lastSeq int64 = cctx.Int64("backfill-seq")

	if lastSeq == 0 {
		log.Infof("No backfill seq specified, so attempting to backfill from the last line of the output file...")
		seqno, err := output.GetBackfillSeqno()
		if err != nil {
			log.Warnf("Failed to get backfill seqno: %+v", err)
			log.Warnf("Continuing without backfill...")
		} else {
			log.Infof("Backfilling from seq: %d", seqno)
			lastSeq = seqno
		}
	} else {
		log.Infof("Backfilling from provided seq: %d", lastSeq)
	}

	s := stream.Stream{
		SocketURL:   u,
		Output:      outputChannel,
		Hydrator:    hydrator,
		BackfillSeq: lastSeq,
	}

	go func() {
		for {
			err = s.BeginStreaming(ctx, cctx.Int("worker-count"))
			log.Errorf("Streaming ended unexpectedly: %+v", err)

			if !cctx.Bool("autorestart") {
				log.Infof("Exiting...")
				break
			} else {
				log.Infof("Restarting stream...")
			}
		}
		cancel()
	}()

	go output.StreamOutput(ctx)

	if cctx.Bool("autorestart") {
		log.Infof("Autorestart is enabled! Stream will restart if it dies...")
	}

	select {
	case <-signals:
		cancel()
		log.Infof("Shutting down on signal")
	case <-ctx.Done():
		log.Infof("Shutting down on context done")
	}

	return nil
}

func repodumpCmd(cctx *cli.Context) error {
	ctx := cctx.Context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Authenticate
	authInfo, err := authenticate(cctx)
	if err != nil {
		log.Fatalf("Failed to authenticate: %+v", err)
		return err
	}

	hydrator, err := hydrator.MakeHydrator(cctx.Context, cctx.Int64("cache-size"), authInfo)
	if err != nil {
		log.Fatalf("Failed to create hydrator: %+v", err)
		return err
	}

	// Create a client
	client := &repodump.RepoDump{
		PdsQueue:     make(map[string]bool),
		PdsCompleted: make(map[string]bool),
		Hydrator:     hydrator,
		Output:       make(chan repodump.CarOutput),
	}

	// Start downloading repos
	go func() {
		err := client.BeginDownloading(ctx)
		log.Errorf("Downloading ended unexpectedly: %+v", err)
		cancel()
	}()

	// Start writing repos to the output folder
	go func() {
		for carOutput := range client.Output {
			// The folder is the output folder, followed by the first two
			// characters of the DID, followed by the second two characters of
			// the DID e.g., output/ab/12/did:plc:ab1234.car. This is to prevent too
			// many files in a single directory. Note that we have to skip the `did:plc:`
			// prefix.
			folder := fmt.Sprintf("%s/%s/%s", cctx.String("output-folder"), carOutput.Did[8:10], carOutput.Did[10:12])
			locationOnDisk := fmt.Sprintf("%s/%s.car", folder, carOutput.Did)

			// Ensure necessary directories exist
			err := os.MkdirAll(folder, 0755)
			if err != nil {
				log.Errorf("Failed to create output folder: %+v", err)
				continue
			}

			// Write the car to disk
			err = os.WriteFile(locationOnDisk, carOutput.Data, 0644)
			if err != nil {
				log.Errorf("Failed to write car to output folder: %+v", err)
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
