package stream

// Borrowed in part from https://github.com/bluesky-social/indigo/blob/main/sonar/sonar.go

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/repomgr"
	log "github.com/sirupsen/logrus"

	// "github.com/bluesky-social/indigo/api/bsky"
	"time"

	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/events/schedulers/autoscaling"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/gorilla/websocket"
)

type Stream struct {
	// SocketURL is the full websocket path to the ATProto SubscribeRepos XRPC endpoint
	SocketURL *url.URL
	Events    chan string // TODO: structured data
}

func (s *Stream) BeginStreaming(ctx context.Context, workerCount int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	scalingSettings := autoscaling.DefaultAutoscaleSettings()
	scalingSettings.MaxConcurrency = workerCount
	scalingSettings.AutoscaleFrequency = time.Second

	pool := autoscaling.NewScheduler(scalingSettings, s.SocketURL.Host, s.HandleStreamEvent)

	log.Infof("connecting to WebSocket at: %s", s.SocketURL.String())
	c, _, err := websocket.DefaultDialer.Dial(s.SocketURL.String(), http.Header{
		"User-Agent": []string{"sonar/1.0"},
	})

	if err != nil {
		log.Infof("failed to connect to websocket: %v", err)
		return err
	}
	defer c.Close()

	go func() {
		err = events.HandleRepoStream(ctx, c, pool)
		log.Infof("HandleRepoStream returned unexpectedly: %+v...", err)
		cancel()
	}()

	select {
	case <-ctx.Done():
		log.Infof("shutting down...")
	}

	return nil
}

func (s *Stream) HandleStreamEvent(ctx context.Context, xe *events.XRPCStreamEvent) error {
	if xe.Error != nil {
		log.Errorf("error handling stream event: %+v", xe.Error)
	}

	if xe.RepoCommit != nil {
		return s.HandleRepoCommit(ctx, xe.RepoCommit)
	}
	return nil
}

func (s *Stream) HandleRepoCommit(ctx context.Context, evt *comatproto.SyncSubscribeRepos_Commit) error {
	// TODO: Actually do something with the event. Below is copied from the Sonar command.

	rr, err := repo.ReadRepoFromCar(ctx, bytes.NewReader(evt.Blocks))
	if err != nil {
		log.Warnf("failed to read repo from car: %+v", err)
		return nil
	}

	// Parse time from the event time string
	evtCreatedAt, err := time.Parse(time.RFC3339, evt.Time)
	if err != nil {
		log.Warnf("error parsing time: %+v", err)
		return nil
	} else {
		log.Infof("event created at: %+v", evtCreatedAt)
	}

	for _, op := range evt.Ops {
		collection := strings.Split(op.Path, "/")[0]

		ek := repomgr.EventKind(op.Action)
		log_wf := log.WithFields(log.Fields{"action": op.Action, "collection": collection})

		switch ek {
		case repomgr.EvtKindCreateRecord, repomgr.EvtKindUpdateRecord:
			// Grab the record from the merkel tree
			rc, rec, err := rr.GetRecord(ctx, op.Path)
			if err != nil {
				e := fmt.Errorf("getting record %s (%s) within seq %d for %s: %w", op.Path, *op.Cid, evt.Seq, evt.Repo, err)
				log_wf.Errorf("failed to get a record from the event: %+v", e)
				break
			}

			// Verify that the record cid matches the cid in the event
			if lexutil.LexLink(rc) != *op.Cid {
				e := fmt.Errorf("mismatch in record and op cid: %s != %s", rc, *op.Cid)
				log_wf.Errorf("failed to LexLink the record in the event: %+v", e)
				break
			}

			// Unpack the record and process it
			switch rec := rec.(type) {
			default:
				log_wf.Warnf("unknown record type: %+v", rec)
			}

			s.Events <- fmt.Sprintf("event: %+v", rec)

		case repomgr.EvtKindDeleteRecord:
		default:
			log.Warnf("unknown event kind from op action: %+v", op.Action)
		}
	}

	return nil
}
