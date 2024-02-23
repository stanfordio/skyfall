package stream

// Borrowed in part from https://github.com/bluesky-social/indigo/blob/main/sonar/sonar.go

import (
	"context"
	"sync"

	"github.com/bluesky-social/indigo/api/atproto"
	comatproto "github.com/bluesky-social/indigo/api/atproto"
	indigoutil "github.com/bluesky-social/indigo/util"
	"github.com/bluesky-social/indigo/xrpc"
	log "github.com/sirupsen/logrus"
	"github.com/stanfordio/skyfall/pkg/hydrator"
	// "github.com/bluesky-social/indigo/api/bsky"
)

type CarOutput struct {
	Did  string
	Data []byte
}

type RepoDump struct {
	PdsQueue     map[string]bool // Using a map as a set
	PdsCompleted map[string]bool // Using a map as a set
	Hydrator     *hydrator.Hydrator
	Output       chan CarOutput
}

type carPullRequest struct {
	pdsEndpoint string
	did         string
}

func (s *RepoDump) startRepoDownloader(ctx context.Context, carChan chan *carPullRequest, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			xrpcClient := xrpc.Client{
				Client: indigoutil.RobustHTTPClient(),
				Host:   "", // Will be set by the loop below
			}

			for downloadRequest := range carChan {
				// Download the car
				log.Debugf("Downloading car: %s from %s", downloadRequest.did, downloadRequest.pdsEndpoint)

				// Pull the bytes
				s.Hydrator.Ratelimit.Take()
				xrpcClient.Host = downloadRequest.pdsEndpoint
				repoBytes, err := atproto.SyncGetRepo(ctx, &xrpcClient, downloadRequest.did, "")
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

func (s *RepoDump) BeginDownloading(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// First, seed the PDS queue with our own PDS
	selfIdentity, err := s.Hydrator.LookupIdentity(s.Hydrator.AuthInfo.Handle)
	if err != nil {
		log.Errorf("Failed to get own PDS: %v", err)
		return err
	}

	// Start the downloader
	carDownloadChannel := make(chan *carPullRequest)
	var wg sync.WaitGroup
	go s.startRepoDownloader(ctx, carDownloadChannel, &wg)

	// Add our own PDS to the queue
	s.PdsQueue[selfIdentity.PDSEndpoint()] = true

	for len(s.PdsQueue) > 0 {
		// Pop the first PDS from the queue
		var pdsEndpoint string
		for pdsEndpoint = range s.PdsQueue {
			break
		}
		delete(s.PdsQueue, pdsEndpoint)

		xrpcClient := &xrpc.Client{
			Client: indigoutil.RobustHTTPClient(),
			Host:   pdsEndpoint,
		}

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
				log.Infof("Pulling DID from %s: %s", pdsEndpoint, r.Did)
				carDownloadChannel <- &carPullRequest{
					pdsEndpoint: pdsEndpoint,
					did:         r.Did,
				}
			}
		}

		// Wait for the downloaders to finish on this PDS before moving on
		// to the next one
		log.Infof("Waiting for downloaders to finish on %s", pdsEndpoint)
		wg.Wait()
		log.Infof("Downloaders finished on %s", pdsEndpoint)
	}

	<-ctx.Done()
	log.Infof("Shutting down...")

	return nil
}
