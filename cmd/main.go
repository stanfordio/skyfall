package main

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/ipfs/go-cid"
	"github.com/stanfordio/skyfall/pkg/auth"
	"github.com/stanfordio/skyfall/pkg/hydrator"
	"github.com/stanfordio/skyfall/pkg/output"
	pull "github.com/stanfordio/skyfall/pkg/pull"
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
				Name:   "stream",
				Usage:  "Sip from the firehose",
				Action: streamCmd,
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
				Name:   "pull",
				Usage:  "Pull all content and write it to a file or BigQuery",
				Action: pullCmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "intermediate-state",
						Usage: "file to intermediate state to (important for resumption)",
						Value: "pull-intermediate-state.json",
					},
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
				},
			},
			{
				Name:   "hydrate",
				Usage:  "Hydrate a folder of .car files into the same format as the stream",
				Action: hydrateCmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Usage:    "folder or file to read data from",
						Required: true,
					},
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
					&cli.StringFlag{
						Name:  "output-bq-table",
						Usage: "name of a BigQuery table to output to in ID form (e.g., dgap_bsky.example_table)",
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

	handle := os.Getenv("BLUESKY_HANDLE")
	if handle == "" {
		handle = cctx.String("handle")
		if handle == "" {
			log.Fatal("No handle provided")
			return nil, errors.New("No handle provided")
		}
	}

	password := os.Getenv("BLUESKY_PASSWORD")
	if password == "" {
		password = cctx.String("password")
		if password == "" {
			log.Fatal("No password provided")
			return nil, errors.New("No password provided")
		}
	}

	authInfo, err := authenticator.Authenticate(handle, password)
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

	outputChannel := make(chan map[string]interface{}, 10000)

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

func pullCmd(cctx *cli.Context) error {
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

	// Create the output channel
	outputChannel := make(chan map[string]interface{}, 10000)

	// Create a client
	client := &pull.Pull{
		PdsCursor:             "",
		IntermediateStatePath: cctx.String("intermediate-state"),
		Output:                outputChannel,
		Hydrator:              hydrator,
	}

	// Setup the output
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

	// Start downloading repos
	go func() {
		err := client.BeginDownloading(ctx, cctx.Int("worker-count"))
		log.Errorf("Downloading ended unexpectedly: %+v", err)
		cancel()
	}()

	// Start the output stream
	go output.StreamOutput(ctx)

	select {
	case <-signals:
		cancel()
		log.Infof("Shutting down on signal")
	case <-ctx.Done():
		log.Infof("Shutting down on context done")
	}

	return nil
}

func hydrateCmd(cctx *cli.Context) error {
	ctx := cctx.Context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Create a client
	log.Infof("Authenticating...")
	authInfo, err := authenticate(cctx)
	if err != nil {
		log.Fatalf("Failed to authenticate: %+v", err)
		return err
	}

	log.Infof("Creating hydrator...")
	hydrator, err := hydrator.MakeHydrator(cctx.Context, cctx.Int64("cache-size"), authInfo)
	if err != nil {
		log.Fatalf("Failed to create hydrator: %+v", err)
		return err
	}

	outputChannel := make(chan map[string]interface{}, 10000)

	log.Infof("Creating output...")
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

	// Find all the CARs; i.e., every `.car` file in the input folder (could be nested)
	// and then hydrate them.
	carFiles := make(chan string, 10000)
	go func() {
		defer close(carFiles)
		carFilesCount := 0
		input := cctx.String("input")
		err = filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".car") {
				carFiles <- path
				carFilesCount++
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Failed to walk input folder: %+v", err)
		}

		log.Infof("Found %d car files to hydrate in %s.", carFilesCount, input)
	}()

	// Spawn workers to hydrate the CARs
	for i := 0; i < cctx.Int("worker-count"); i++ {
		go func() {
			for carFile := range carFiles {
				log.Infof("Hydrating %s", carFile)
				// Read the car file
				data, err := os.ReadFile(carFile)
				if err != nil {
					log.Errorf("Failed to read car file: %+v", err)
					continue
				}
				repo, err := repo.ReadRepoFromCar(ctx, bytes.NewReader(data))
				if err != nil {
					log.Errorf("Failed to read repo from car: %+v", err)
					continue
				}
				// Extract the actor (i.e., whose repo is this?)
				actorDid := repo.RepoDid()

				// Hydrate the repo
				err = repo.ForEach(ctx, "", func(k string, v cid.Cid) error {
					// Get the record
					_, rec, err := repo.GetRecord(ctx, k)
					if err != nil {
						log.Errorf("Unable to parse CID %s from %s: %s", v.String(), actorDid, err)
						return err
					}

					// Hydrate the record
					hydrated, err := hydrator.Hydrate(rec, actorDid)
					if err != nil {
						log.Errorf("Failed to hydrate record: %+v", err)
						return err
					}

					// Write the hydrated record to the output
					outputChannel <- hydrated

					return nil
				})

				if err != nil {
					log.Errorf("Failed to hydrate repo: %+v", err)
				}
			}
		}()
	}

	go output.StreamOutput(ctx)

	select {
	case <-signals:
		cancel()
		log.Infof("Shutting down on signal")
	case <-ctx.Done():
		log.Infof("Shutting down on context done")
	}

	return nil
}
