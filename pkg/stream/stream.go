package stream

// Borrowed in part from https://github.com/bluesky-social/indigo/blob/main/sonar/sonar.go

import (
	"bytes"
	"context"
	"encoding/json"
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
	hydrator "github.com/stanfordio/skyfall/pkg/hydrator"
)

type Stream struct {
	// SocketURL is the full websocket path to the ATProto SubscribeRepos XRPC endpoint
	SocketURL *url.URL
	Output    chan []byte
	Hydrator  *hydrator.Hydrator
}

func (s *Stream) BeginStreaming(ctx context.Context, workerCount int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	scalingSettings := autoscaling.DefaultAutoscaleSettings()
	scalingSettings.MaxConcurrency = workerCount
	scalingSettings.AutoscaleFrequency = time.Second

	pool := autoscaling.NewScheduler(scalingSettings, s.SocketURL.Host, s.HandleStreamEvent)

	log.Infof("Connecting to WebSocket at: %s", s.SocketURL.String())
	c, _, err := websocket.DefaultDialer.Dial(s.SocketURL.String(), http.Header{
		"User-Agent": []string{"sonar/1.0"},
	})

	if err != nil {
		log.Infof("Failed to connect to websocket: %v", err)
		return err
	}
	defer c.Close()

	go func() {
		err = events.HandleRepoStream(ctx, c, pool)
		log.Infof("HandleRepoStream returned unexpectedly: %+v...", err)
		cancel()
	}()

	<-ctx.Done()
	log.Infof("Shutting down...")

	return nil
}

func (s *Stream) HandleStreamEvent(ctx context.Context, xe *events.XRPCStreamEvent) error {
	if xe.Error != nil {
		log.Errorf("Error handling stream event: %+v", xe.Error)
	}

	if xe.RepoCommit != nil {
		return s.HandleRepoCommit(ctx, xe.RepoCommit)
	} else {
		log.Warnf("Unknown stream event: %+v", xe)
	}
	return nil
}

func (s *Stream) HandleRepoCommit(ctx context.Context, evt *comatproto.SyncSubscribeRepos_Commit) (error error) {
	rr, err := repo.ReadRepoFromCar(ctx, bytes.NewReader(evt.Blocks))
	if err != nil {
		log.Warnf("Failed to read repo from car: %+v", err)
		return nil
	}

	// Extract the actor (i.e., whose repo is this?)
	actorDid := rr.RepoDid()

	error = nil

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
				log_wf.Errorf("Failed to get a record from the event: %+v", e)
				error = e
				break
			}

			// Verify that the record cid matches the cid in the event
			if lexutil.LexLink(rc) != *op.Cid {
				e := fmt.Errorf("mismatch in record and op cid: %s != %s", rc, *op.Cid)
				log_wf.Errorf("Failed to LexLink the record in the event: %+v", e)
				error = e
				break
			}

			// Hydrate the record
			hydrated, err := s.Hydrator.Hydrate(rec, actorDid)
			if err != nil {
				log_wf.Errorf("Failed to hydrate record: %+v", err)
				error = err
			}

			val, err := json.Marshal(hydrated)

			if err != nil {
				log.Errorf("Failed to marshal record: %+v", err)
				break
			}

			s.Output <- val

		case repomgr.EvtKindDeleteRecord:
			log.Warnf("Delete record not implemented yet: %+c", op)
		default:
			log.Warnf("Unknown event kind from op action: %+v", op.Action)
		}
	}

	return
}
