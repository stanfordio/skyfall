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
	hydrator "github.com/stanfordio/skyfall/pkg/hydrator"
)

type Stream struct {
	// SocketURL is the full websocket path to the ATProto SubscribeRepos XRPC endpoint
	SocketURL   *url.URL
	Output      chan map[string]interface{}
	Hydrator    *hydrator.Hydrator
	BackfillSeq int64
}

func (s *Stream) BeginStreaming(ctx context.Context, workerCount int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	scalingSettings := autoscaling.DefaultAutoscaleSettings()
	scalingSettings.MaxConcurrency = workerCount
	scalingSettings.AutoscaleFrequency = time.Second

	pool := autoscaling.NewScheduler(scalingSettings, s.SocketURL.Host, s.HandleStreamEvent)

	var socketUrl string = s.SocketURL.String()
	// If we're backfilling, add the backfill seq to the socket url
	if s.BackfillSeq > 0 {
		socketUrl = fmt.Sprintf("%s?cursor=%d", socketUrl, s.BackfillSeq)
	}

	log.Infof("Connecting to WebSocket at: %s", socketUrl)
	c, _, err := websocket.DefaultDialer.Dial(socketUrl, http.Header{
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

			// Log the action performed
			hydrated["Action"] = op.Action

			// Include the event sequence number
			hydrated["Seq"] = evt.Seq

			s.Output <- hydrated

		case repomgr.EvtKindDeleteRecord:
			// Not much we can do here, since we don't have the record anymore; just log the action
			hydrated, err := s.Hydrator.Hydrate(map[string]interface{}{"CreatedAt": time.Now().Format(time.RFC3339), "Item": op.Path, "LexiconTypeID": strings.Split(op.Path, "/")[0]}, actorDid)
			if err != nil {
				log_wf.Errorf("Failed to hydrate record: %+v", err)
				error = err
				break
			}

			hydrated["Action"] = op.Action
			hydrated["Seq"] = evt.Seq

			s.Output <- hydrated
		default:
			log.Warnf("Unknown event kind from op action: %+v", op.Action)
		}
	}

	return
}
