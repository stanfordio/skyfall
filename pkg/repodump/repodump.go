package stream

// Borrowed in part from https://github.com/bluesky-social/indigo/blob/main/sonar/sonar.go

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"sync"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/repo"
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
	PdsQueue     map[string]bool // Using a map as a set
	PdsCompleted map[string]bool // Using a map as a set
}

type RepoDump struct {
	PdsQueue              map[string]bool // Using a map as a set
	PdsCompleted          map[string]bool // Using a map as a set
	SkipDids              func(string) bool
	IntermediateStatePath string
	Hydrator              *hydrator.Hydrator
	Output                chan CarOutput
}

type carPullRequest struct {
	pdsEndpoint string
	did         string
}

func (s *RepoDump) startRepoDownloader(ctx context.Context, numWorkers int, carChan chan chan *carPullRequest, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		go func() {
			for carDownloadChannel := range carChan {
				wg.Add(1)
				for downloadRequest := range carDownloadChannel {
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

					// Parse the repo so that we can pull all the identities in the repo
					repo, err := repo.ReadRepoFromCar(ctx, bytes.NewReader(repoBytes))
					if err != nil {
						log.Errorf("Unable to read repo %s: %s", downloadRequest.did, err)
					}

					// Pull all the identities
					identities, err := s.Hydrator.GetIdentitiesInRepo(repo)
					if err != nil {
						log.Errorf("Unable to read identities in repo %s: %s", downloadRequest.did, err)
					}

					log.Infof("Found %d identities in repo %s", len(identities), downloadRequest.did)

					// Find the unique PDSes
					didFindNewPds := false
					for _, identity := range identities {
						if _, ok := s.PdsCompleted[identity.PDSEndpoint()]; !ok {
							if _, ok := s.PdsQueue[identity.PDSEndpoint()]; !ok {
								log.Infof("Adding PDS to queue: %s", identity.PDSEndpoint())
								s.PdsQueue[identity.PDSEndpoint()] = true
								didFindNewPds = true
							}
						}
					}

					if didFindNewPds {
						// Save the state to disk
						err = s.saveIntermediateStateToDisk()
						if err != nil {
							log.Errorf("Failed to save intermediate state to disk: %v", err)
						}
					}

					// Write the bytes to the output channel
					s.Output <- CarOutput{
						Did:  downloadRequest.did,
						Data: repoBytes,
					}
				}
				wg.Done()
			}
		}()
	}
}

func (s *RepoDump) saveIntermediateStateToDisk() error {
	// Saves the pull queue and the completed queue to disk so that we can
	// resume the download later if needed.

	state := intermediateState{
		PdsQueue:     s.PdsQueue,
		PdsCompleted: s.PdsCompleted,
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
			PdsQueue:     make(map[string]bool),
			PdsCompleted: make(map[string]bool),
		}
		err = json.Unmarshal(in, &state)

		if err != nil {
			return err
		}

		importedCount := 0

		// Set the state additively to the current state
		for k, v := range state.PdsCompleted {
			s.PdsCompleted[k] = v
			importedCount++
		}
		for k, v := range state.PdsQueue {
			if _, ok := s.PdsCompleted[k]; !ok {
				s.PdsQueue[k] = v
				importedCount++
			}
		}

		log.Infof("Imported %d PDSes from intermediate state", importedCount)
	}

	return nil
}

func (s *RepoDump) BeginDownloading(ctx context.Context, numWorkers int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// First, seed the PDS queue with our own PDS
	selfIdentity, err := s.Hydrator.LookupIdentity(s.Hydrator.AuthInfo.Handle)
	if err != nil {
		log.Errorf("Failed to get own PDS: %v", err)
		return err
	}

	// Start the downloader
	carDownloadChannels := make(chan chan *carPullRequest) // A channel of channels; each subchannel is a queue of carPullRequests from the same PDS
	var wg sync.WaitGroup
	go s.startRepoDownloader(ctx, numWorkers, carDownloadChannels, &wg)

	err = s.loadIntermediateStateFromDisk()
	if err != nil {
		log.Errorf("Failed to load intermediate state from disk, so will not be resuming from previous pull: %v", err)
	}

	// Add our own PDS to the queue
	s.PdsQueue[selfIdentity.PDSEndpoint()] = true

	for len(s.PdsQueue) > 0 {
		// Save the state to disk
		err = s.saveIntermediateStateToDisk()
		if err != nil {
			log.Errorf("Failed to save intermediate state to disk: %v", err)
		}

		// Pop the first PDS from the queue; we will move it to the completed queue
		// once we are done with it (at the end of this loop iteration)
		var pdsEndpoint string
		for pdsEndpoint = range s.PdsQueue {
			break
		}

		xrpcClient := &xrpc.Client{
			Client: utils.RetryingHTTPClient(),
			Host:   pdsEndpoint,
		}

		// Create the channel and add it to the downloaders
		channel := make(chan *carPullRequest)
		carDownloadChannels <- channel

		cursor := ""
		for {
			s.Hydrator.Ratelimit.Take()
			out, err := comatproto.SyncListRepos(ctx, xrpcClient, cursor, 1000)
			if err != nil {
				log.Errorf("Failed to get list of repos: %v", err)
				return err
			}
			if len(out.Repos) == 0 {
				log.Infof("Finished pulling DIDs from: %s", pdsEndpoint)
				break
			}
			cursor = *out.Cursor

			// Go through and pull each repo
			for _, r := range out.Repos {
				log.Infof("Pulling CAR from %s: %s", pdsEndpoint, r.Did)
				channel <- &carPullRequest{
					pdsEndpoint: pdsEndpoint,
					did:         r.Did,
				}
			}
		}

		// Close the channel so that the downloaders know that they are done
		close(channel)

		// Wait for the downloaders to finish on this PDS before moving on
		// to the next one
		log.Infof("Waiting for downloaders to finish on %s", pdsEndpoint)
		wg.Wait()
		log.Infof("Downloaders finished on %s", pdsEndpoint)

		delete(s.PdsQueue, pdsEndpoint)
		s.PdsCompleted[pdsEndpoint] = true
	}

	// Close the channel for the downloaders
	close(carDownloadChannels)

	<-ctx.Done()
	log.Infof("Shutting down...")

	return nil
}
