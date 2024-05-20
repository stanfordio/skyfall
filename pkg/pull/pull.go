package pull

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"sync"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"
	"github.com/stanfordio/skyfall/pkg/hydrator"
	"github.com/stanfordio/skyfall/pkg/utils"
	// "github.com/bluesky-social/indigo/api/bsky"
)

type intermediateState struct {
	PdsCursor string
}

type Pull struct {
	PdsCursor             string
	IntermediateStatePath string
	Output                chan map[string]interface{}
	Hydrator              *hydrator.Hydrator
}

type carPullRequest struct {
	pdsEndpoint string
	did         string
}

func (s *Pull) startDownloader(ctx context.Context, numWorkers int, carChan chan *carPullRequest, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			for downloadRequest := range carChan {
				// Download the car
				log.Infof("Downloading car: %s from %s", downloadRequest.did, downloadRequest.pdsEndpoint)

				// Pull the bytes
				repoBytes, err := s.Hydrator.GetRepoBytes(downloadRequest.did, downloadRequest.pdsEndpoint)
				if err != nil {
					log.Errorf("Failed to download car %s from %s: %v", downloadRequest.did, downloadRequest.pdsEndpoint, err)
					continue
				}
				repo, err := repo.ReadRepoFromCar(ctx, bytes.NewReader(repoBytes))
				if err != nil {
					log.Errorf("Failed to read car %s from %s: %v", downloadRequest.did, downloadRequest.pdsEndpoint, err)
					continue
				}

				actorDid := repo.RepoDid()

				// Hydrate the car and send it to the output
				err = repo.ForEach(ctx, "", func(k string, v cid.Cid) error {
					// Grab the record from the merkel tree
					_, rec, err := repo.GetRecord(ctx, k)
					if err != nil {
						log.Errorf("Failed to get record %s from car %s from %s: %v", v.String(), downloadRequest.did, downloadRequest.pdsEndpoint, err)
						return err
					}

					// Hydrate the record
					hydrated, err := s.Hydrator.Hydrate(rec, actorDid)
					if err != nil {
						log.Errorf("Failed to hydrate record: %+v", err)
						return err
					}

					// Output the record (it'll be thrown into BigQuery or the
					// output file)
					s.Output <- hydrated

					return nil
				})
				if err != nil {
					log.Errorf("Failed to hydrate car %s from %s: %v", downloadRequest.did, downloadRequest.pdsEndpoint, err)
					continue
				}
			}
			wg.Done()
		}()
	}
}

func (s *Pull) saveIntermediateStateToDisk() error {
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

func (s *Pull) loadIntermediateStateFromDisk() error {
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

		s.PdsCursor = state.PdsCursor
	}

	return nil
}

func (s *Pull) BeginDownloading(ctx context.Context, numWorkers int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start the downloader
	carDownloadChannel := make(chan *carPullRequest, 10000)
	var wg sync.WaitGroup
	go s.startDownloader(ctx, numWorkers, carDownloadChannel, &wg)

	err := s.loadIntermediateStateFromDisk()
	if err != nil {
		log.Errorf("Failed to load intermediate state from disk, so will not be resuming from previous pull: %v", err)
	} else {
		log.Infof("Loaded intermediate state from disk; cursor = %s", s.PdsCursor)
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
		} else {
			log.Infof("Got %d repos from %s (cursor = %s)", len(out.Repos), pdsEndpoint, s.PdsCursor)
		}
		if len(out.Repos) == 0 {
			log.Infof("Finished pulling DIDs from: %s", pdsEndpoint)
			break
		}
		s.PdsCursor = *out.Cursor
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

		err = s.saveIntermediateStateToDisk()
	}

	// Close the channel so that the downloaders know that they are done
	close(carDownloadChannel)

	// Wait for the downloaders to finish on this PDS before moving on
	// to the next one
	log.Infof("Waiting for downloaders to finish on %s", pdsEndpoint)
	wg.Wait()
	log.Infof("Downloaders finished on %s", pdsEndpoint)

	// Close the output channel
	close(s.Output)

	<-ctx.Done()
	log.Infof("Shutting down...")

	return nil
}
