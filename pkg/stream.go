package stream

// Borrowed in part from https://github.com/bluesky-social/indigo/blob/main/sonar/sonar.go

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"log"

	"github.com/araddon/dateparse"
	comatproto "github.com/bluesky-social/indigo/api/atproto"

	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/repomgr"
	"github.com/bluesky-social/indigo/api/bsky"
)

type Stream struct {
	// SocketURL is the full websocket path to the ATProto SubscribeRepos XRPC endpoint
	SocketURL string
}

func (s *Stream) HandleStreamEvent(ctx context.Context, xe *events.XRPCStreamEvent) error {
	if xe.Error != nil {
		// TODO: Log the error
	}

	if xe.RepoCommit != nil {
		return s.HandleRepoCommit(ctx, xe.RepoCommit)
	}
	return nil
}

func (s *Stream) HandleRepoCommit(ctx context.Context, evt *comatproto.SyncSubscribeRepos_Commit) error {
	// TODO: Actually do something with the event. Below is copied from the Sonar command.
	
	processedAt := time.Now()

	rr, err := repo.ReadRepoFromCar(ctx, bytes.NewReader(evt.Blocks))
	if err != nil {
		log.Errorf("failed to read repo from car: %+v", err)
		return nil
	}

	if evt.Rebase {
		log.Debug("rebase")
	}

	// Parse time from the event time string
	evtCreatedAt, err := time.Parse(time.RFC3339, evt.Time)
	if err != nil {
		log.Errorf("error parsing time: %+v", err)
		return nil
	}
	
	for _, op := range evt.Ops {
		collection := strings.Split(op.Path, "/")[0]

		ek := repomgr.EventKind(op.Action)
		log = log.With("action", op.Action, "collection", collection)

		opsProcessedCounter.WithLabelValues(op.Action, collection, s.SocketURL).Inc()

		switch ek {
		case repomgr.EvtKindCreateRecord, repomgr.EvtKindUpdateRecord:
			// Grab the record from the merkel tree
			rc, rec, err := rr.GetRecord(ctx, op.Path)
			if err != nil {
				e := fmt.Errorf("getting record %s (%s) within seq %d for %s: %w", op.Path, *op.Cid, evt.Seq, evt.Repo, err)
				log.Errorf("failed to get a record from the event: %+v", e)
				break
			}

			// Verify that the record cid matches the cid in the event
			if lexutil.LexLink(rc) != *op.Cid {
				e := fmt.Errorf("mismatch in record and op cid: %s != %s", rc, *op.Cid)
				log.Errorf("failed to LexLink the record in the event: %+v", e)
				break
			}

			var recCreatedAt time.Time
			var parseError error

			// Unpack the record and process it
			switch rec := rec.(type) {
			default:
				log.Warnf("unknown record type: %+v", rec)
			}
			if parseError != nil {
				log.Errorf("error parsing time: %+v", parseError)
				continue
			}
			if !recCreatedAt.IsZero() {
				lastEvtCreatedAtGauge.WithLabelValues(s.SocketURL).Set(float64(recCreatedAt.UnixNano()))
				lastEvtCreatedRecordCreatedGapGauge.WithLabelValues(s.SocketURL).Set(float64(evtCreatedAt.Sub(recCreatedAt).Seconds()))
				lastRecordCreatedEvtProcessedGapGauge.WithLabelValues(s.SocketURL).Set(float64(processedAt.Sub(recCreatedAt).Seconds()))
			}

		case repomgr.EvtKindDeleteRecord:
		default:
			log.Warnf("unknown event kind from op action: %+v", op.Action)
		}
	}

	eventProcessingDurationHistogram.WithLabelValues(s.SocketURL).Observe(time.Since(processedAt).Seconds())
	return nil
}