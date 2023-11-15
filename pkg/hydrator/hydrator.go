package hydrator

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bluesky-social/indigo/atproto/identity"
	atpidentity "github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	indigoutil "github.com/bluesky-social/indigo/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/dgraph-io/ristretto"
	"github.com/mitchellh/mapstructure"
)

type Hydrator struct {
	Cache             *ristretto.Cache
	Context           context.Context
	Client            *xrpc.Client
	IdentityDirectory identity.Directory
}

func MakeHydrator(ctx context.Context, cacheSize int64, hostDomain string) (*Hydrator, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e8, // number of keys to track frequency of
		MaxCost:     cacheSize,
		BufferItems: 64, // number of keys per Get buffer
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %+v", err)
	}

	h := Hydrator{
		Cache:   cache,
		Context: ctx,
		Client: &xrpc.Client{
			Client: indigoutil.RobustHTTPClient(),
			Host:   fmt.Sprintf("https://%s", hostDomain),
		},
		IdentityDirectory: identity.DefaultDirectory(),
	}

	return &h, nil
}

func (h *Hydrator) lookupIdentity(identifier string) (identity *identity.Identity, err error) {
	// Check the cache first
	cachedValue, found := h.Cache.Get(identifier)

	if found && cachedValue != nil {
		identity = cachedValue.(*atpidentity.Identity)
		return
	}

	log.Debugf("Cache miss for %s", identifier)

	resolvedIdentifier, error := syntax.ParseAtIdentifier(identifier)
	if error != nil {
		err = error
		return
	}

	identity, err = h.IdentityDirectory.Lookup(h.Context, *resolvedIdentifier)
	if err != nil {
		return
	}

	h.Cache.SetWithTTL(identifier, identity, 1, time.Duration(1)*time.Hour*24)

	return
}

func (h *Hydrator) Hydrate(val interface{}, actorDid string) (result map[string]interface{}, err error) {
	err = nil

	err = mapstructure.Decode(val, &result)

	if err != nil {
		return
	}

	// Resolve full identity and profile information for the actor
	identity, err := h.lookupIdentity(actorDid)
	if err != nil {
		log.Warnf("Failed to get profile for actor: %s", actorDid)
		identity = nil
	}

	// Add the actorDid and profile to the map
	result["_ActorDid"] = actorDid
	result["_ActorIdentity"] = identity

	return
}
