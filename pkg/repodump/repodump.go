package stream

// Borrowed in part from https://github.com/bluesky-social/indigo/blob/main/sonar/sonar.go

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	log "github.com/sirupsen/logrus"
	"github.com/stanfordio/skyfall/pkg/hydrator"
	"github.com/stanfordio/skyfall/pkg/utils"
	// "github.com/bluesky-social/indigo/api/bsky"
)

type CarOutput struct {
	Did  string
	Data []byte
}

type intermediateState struct {
	PdsCursor string
}

type RepoDump struct {
	PdsCursor             string
	SkipDids              func(string) bool
	IntermediateStatePath string
	Hydrator              *hydrator.Hydrator
	Output                chan CarOutput
}

type carPullRequest struct {
	pdsEndpoint string
	did         string
}

func (s *RepoDump) startRepoDownloader(_ctx context.Context, numWorkers int, carChan chan *carPullRequest, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			for downloadRequest := range carChan {
				// Download the car
				if s.SkipDids(downloadRequest.did) {
					log.Infof("Skipping car: %s from %s (likely already downloaded)", downloadRequest.did, downloadRequest.pdsEndpoint)
					continue
				}
				log.Infof("Downloading car: %s from %s", downloadRequest.did, downloadRequest.pdsEndpoint)

				// Pull the bytes
				repoBytes, err := s.Hydrator.GetRepoBytes(downloadRequest.did, downloadRequest.pdsEndpoint)
				if err != nil {
					log.Errorf("Failed to download car %s from %s: %v", downloadRequest.did, downloadRequest.pdsEndpoint, err)
					continue
				}

				// Write the bytes to the output channel
				s.Output <- CarOutput{
					Did:  downloadRequest.did,
					Data: repoBytes,
				}
			}
			wg.Done()
		}()
	}
}

func (s *RepoDump) saveIntermediateStateToDisk() error {
	// Saves the pull queue and the completed queue to disk so that we can
	// resume the download later if needed.

	state := intermediateState{
		PdsCursor: s.PdsCursor,
	}

	// Marshall into json
	out, err := json.Marshal(state)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.IntermediateStatePath, out, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *RepoDump) loadIntermediateStateFromDisk() error {
	// Loads the pull queue and the completed queue from disk.

	// If the file exists, load it
	if _, err := os.Stat(s.IntermediateStatePath); err == nil {
		// Load the file
		in, err := os.ReadFile(s.IntermediateStatePath)
		if err != nil {
			return err
		}

		// Unmarshall the json
		state := intermediateState{
			PdsCursor: s.PdsCursor,
		}
		err = json.Unmarshal(in, &state)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *RepoDump) BeginDownloading(ctx context.Context, numWorkers int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start the downloader
	carDownloadChannel := make(chan *carPullRequest) // A channel of channels; each subchannel is a queue of carPullRequests from the same PDS
	var wg sync.WaitGroup
	go s.startRepoDownloader(ctx, numWorkers, carDownloadChannel, &wg)

	err := s.loadIntermediateStateFromDisk()
	if err != nil {
		log.Errorf("Failed to load intermediate state from disk, so will not be resuming from previous pull: %v", err)
	}

	pdsEndpoint := "https://bsky.network"

	xrpcClient := &xrpc.Client{
		Client: utils.RetryingHTTPClient(),
		Host:   pdsEndpoint,
	}

	// Create the channel and add it to the downloaders
	for {
		s.Hydrator.Ratelimit.Take()
		out, err := comatproto.SyncListRepos(ctx, xrpcClient, s.PdsCursor, 1000)
		if err != nil {
			log.Errorf("Failed to get list of repos: %v", err)
			return err
		}
		if len(out.Repos) == 0 {
			log.Infof("Finished pulling DIDs from: %s", pdsEndpoint)
			break
		}
		s.PdsCursor = *out.Cursor
		err = s.saveIntermediateStateToDisk()
		if err != nil {
			log.Errorf("Failed to save intermediate state to disk: %v", err)
			return err
		}

		// Go through and pull each repo
		for _, r := range out.Repos {
			log.Infof("Pulling CAR from %s: %s", pdsEndpoint, r.Did)
			carDownloadChannel <- &carPullRequest{
				pdsEndpoint: pdsEndpoint,
				did:         r.Did,
			}
		}
	}

	// Close the channel so that the downloaders know that they are done
	close(carDownloadChannel)

	// Wait for the downloaders to finish on this PDS before moving on
	// to the next one
	log.Infof("Waiting for downloaders to finish on %s", pdsEndpoint)
	wg.Wait()
	log.Infof("Downloaders finished on %s", pdsEndpoint)

	<-ctx.Done()
	log.Infof("Shutting down...")

	return nil
}
