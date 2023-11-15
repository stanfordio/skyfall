package hydrator

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bluesky-social/indigo/api/bsky"
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
	AuthInfo          *xrpc.AuthInfo
	Context           context.Context
	Client            *xrpc.Client
	IdentityDirectory identity.Directory
}

func MakeHydrator(ctx context.Context, cacheSize int64, hostDomain string, authInfo *xrpc.AuthInfo) (*Hydrator, error) {
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
			Host:   "https://bsky.social", // We generally want to use the bsky.social host for all requests, since they are doing the indexing
			Auth:   authInfo,
		},
		IdentityDirectory: atpidentity.DefaultDirectory(),
		AuthInfo:          authInfo,
	}

	return &h, nil
}

func (h *Hydrator) makeXRPCClientForHost(baseUrl string) *xrpc.Client {
	// TODO: Cache these

	return &xrpc.Client{
		Client: h.Client.Client,
		Host:   baseUrl, // It's called Host, but it's really the base URL
		Auth:   h.AuthInfo,
	}
}

func (h *Hydrator) getXRPCClientForPDS(identity *atpidentity.Identity) *xrpc.Client {
	return h.makeXRPCClientForHost(identity.PDSEndpoint())
}

func (h *Hydrator) lookupIdentity(identifier string) (identity *atpidentity.Identity, err error) {
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

func (h *Hydrator) lookupProfile(identity *atpidentity.Identity) (profile *bsky.ActorDefs_ProfileViewDetailed, err error) {
	profile, err = bsky.ActorGetProfile(h.Context, h.Client, identity.Handle.String())
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

	profile, err := h.lookupProfile(identity)
	if err != nil {
		log.Warnf("Failed to get profile for actor: %s", actorDid)
		profile = nil
	}

	// Add the actorDid and profile to the map
	result["_ActorDid"] = actorDid
	result["_ActorIdentity"] = identity
	result["_ActorProfile"] = profile

	return
}
