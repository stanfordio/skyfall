package pull

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"sort"
	"sync"

	"github.com/bluesky-social/indigo/repo"
	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"
	"github.com/stanfordio/skyfall/pkg/census"
	"github.com/stanfordio/skyfall/pkg/hydrator"
	// "github.com/bluesky-social/indigo/api/bsky"
)

type intermediateState struct {
	FirstUnpulledDidIndex uint64 // All DIDs up to this one have been pulled, 1-indexed
}

type Pull struct {
	CensusPath                  string
	IntermediateStatePath       string
	Output                      chan map[string]interface{}
	Hydrator                    *hydrator.Hydrator
	PdsEndpoint                 string
	FirstUnpulledDidIndex       uint64   // 1-indexed, initialize to 0 by default
	RecentlyPulledCensusIndices []uint64 // initialize to empty slice by default
	CompletedIndicesChannel     chan uint64
}

type carPullRequest struct {
	pdsEndpoint     string
	did             string
	censusFileIndex uint64 // Given Bluesky's current size, this would overflow if we used uint32
}

func (s *Pull) handleDownloadRequest(ctx context.Context, downloadRequest *carPullRequest) error {
	// Download the car
	log.Infof("Downloading car: %s from %s", downloadRequest.did, downloadRequest.pdsEndpoint)

	// Eventually the state management goroutine that we've pulled this DID
	defer func() { s.CompletedIndicesChannel <- downloadRequest.censusFileIndex }()

	// Pull the bytes
	repoBytes, err := s.Hydrator.GetRepoBytes(downloadRequest.did, downloadRequest.pdsEndpoint)
	if err != nil {
		log.Errorf("Failed to download car %s from %s: %v", downloadRequest.did, downloadRequest.pdsEndpoint, err)
		return err
	}
	repo, err := repo.ReadRepoFromCar(ctx, bytes.NewReader(repoBytes))
	if err != nil {
		log.Errorf("Failed to read car %s from %s: %v", downloadRequest.did, downloadRequest.pdsEndpoint, err)
		return err
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
		return err
	}

	return nil
}

func (s *Pull) startDownloader(ctx context.Context, numWorkers int, carChan chan *carPullRequest, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			for downloadRequest := range carChan {
				s.handleDownloadRequest(ctx, downloadRequest)
			}
			wg.Done()
		}()
	}
}

func (s *Pull) saveIntermediateStateToDisk() error {
	// Saves the pull queue and the completed queue to disk so that we can
	// resume the download later if needed.

	state := intermediateState{
		FirstUnpulledDidIndex: s.FirstUnpulledDidIndex,
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
			FirstUnpulledDidIndex: s.FirstUnpulledDidIndex,
		}
		err = json.Unmarshal(in, &state)

		if err != nil {
			return err
		}

		s.FirstUnpulledDidIndex = state.FirstUnpulledDidIndex
	}

	return nil
}

func (s *Pull) keepIntermediateStateUpdated() {
	for completedIndex := range s.CompletedIndicesChannel {
		// Add the completed index to the list of completed indices
		s.RecentlyPulledCensusIndices = append(s.RecentlyPulledCensusIndices, completedIndex)

		// Sort the list of completed indices
		sort.Slice(s.RecentlyPulledCensusIndices, func(i, j int) bool {
			return s.RecentlyPulledCensusIndices[i] < s.RecentlyPulledCensusIndices[j]
		})

		// While the first element in the list is the next one in the sequence,
		// increment the first element and remove it from the list
		log.Debugf("First unpulled DID index: %d, earliest unprocessed: %d", s.FirstUnpulledDidIndex, s.RecentlyPulledCensusIndices[0])
		for len(s.RecentlyPulledCensusIndices) > 0 && s.RecentlyPulledCensusIndices[0] <= s.FirstUnpulledDidIndex {
			s.FirstUnpulledDidIndex = s.RecentlyPulledCensusIndices[0] + 1
			s.RecentlyPulledCensusIndices = s.RecentlyPulledCensusIndices[1:]
		}

		// Save the state to disk
		err := s.saveIntermediateStateToDisk()
		if err != nil {
			log.Fatalf("Failed to save intermediate state to disk: %v", err)
		}
	}
}

func (s *Pull) BeginDownloading(ctx context.Context, numWorkers int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := s.loadIntermediateStateFromDisk()
	if err != nil {
		log.Errorf("Failed to load intermediate state from disk: %v", err)
		return err
	}

	// Start the downloader
	carDownloadChannel := make(chan *carPullRequest, 10000)
	var wg sync.WaitGroup
	go s.startDownloader(ctx, numWorkers, carDownloadChannel, &wg)

	// Start the intermediate state manager
	go s.keepIntermediateStateUpdated()

	// Open the census file
	censusFile, err := os.Open(s.CensusPath)
	if err != nil {
		log.Errorf("Failed to open census file: %v", err)
		return err
	}
	censusFileScanner := bufio.NewScanner(censusFile)

	// Create the channel and add it to the downloaders
	index := uint64(1)
	for censusFileScanner.Scan() {
		line := censusFileScanner.Text()

		lineIndex := index
		index++

		// First, check if we've already pulled this DID
		if index < s.FirstUnpulledDidIndex {
			// Skip this DID
			continue
		}

		// Unmarshal the line
		var repoInfo census.CensusFileEntry
		err := json.Unmarshal([]byte(line), &repoInfo)
		if err != nil {
			log.Infof("Failed to decode census file line : %v", err)
			return err
		}

		// Go through and pull each repo
		carDownloadChannel <- &carPullRequest{
			pdsEndpoint:     s.PdsEndpoint,
			did:             repoInfo.Did,
			censusFileIndex: lineIndex, // 1-indexed
		}

		err = s.saveIntermediateStateToDisk()
	}

	// Close the channel so that the downloaders know that they are done
	close(carDownloadChannel)

	// Wait for the downloaders to finish on this PDS before moving on
	// to the next one
	log.Infof("Waiting for downloaders to finish on %s", s.PdsEndpoint)
	wg.Wait()
	log.Infof("Downloaders finished on %s", s.PdsEndpoint)

	// Close the output channel
	close(s.Output)

	<-ctx.Done()
	log.Infof("Shutting down...")

	return nil
}
