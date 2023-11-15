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

func MakeHydrator(ctx context.Context, cacheSize int64, authInfo *xrpc.AuthInfo) (*Hydrator, error) {
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

func (h *Hydrator) lookupProfileFromIdentity(identity *atpidentity.Identity) (profile *bsky.ActorDefs_ProfileViewDetailed, err error) {
	// Check the cache first
	cachedValue, found := h.Cache.Get(identity.Handle.String())

	if found && cachedValue != nil {
		profile = cachedValue.(*bsky.ActorDefs_ProfileViewDetailed)
		return
	}

	profile, err = bsky.ActorGetProfile(h.Context, h.Client, identity.Handle.String())

	// Set the cache
	if err != nil {
		h.Cache.SetWithTTL(identity.Handle.String(), profile, 1, time.Duration(1)*time.Hour*24)
	}

	return
}

func (h *Hydrator) lookupProfile(did string) (profile *bsky.ActorDefs_ProfileViewDetailed, err error) {
	identity, err := h.lookupIdentity(did)
	if err != nil {
		return
	}
	profile, err = h.lookupProfileFromIdentity(identity)
	return
}

func (h *Hydrator) lookupPost(atUrl string) (post *bsky.FeedDefs_PostView, err error) {
	// Check the cache first
	cachedValue, found := h.Cache.Get(atUrl)

	if found && cachedValue != nil {
		post = cachedValue.(*bsky.FeedDefs_PostView)
		return
	}

	log.Debugf("Cache miss for %s", atUrl)

	output, err := bsky.FeedGetPosts(h.Context, h.Client, []string{atUrl})
	if err != nil {
		return
	}

	if len(output.Posts) == 0 {
		err = fmt.Errorf("no posts found for %s", atUrl)
		return
	}

	post = output.Posts[0]

	h.Cache.SetWithTTL(atUrl, post, 1, time.Duration(1)*time.Hour*24)

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

	profile, err := h.lookupProfileFromIdentity(identity)
	if err != nil {
		log.Warnf("Failed to get profile for actor: %s", actorDid)
		profile = nil
	}

	// Add the actorDid and profile to the map
	result["_ActorDid"] = actorDid
	result["_ActorIdentity"] = identity
	result["_ActorProfile"] = profile

	// Add the pulled time to the map, iso8601
	result["_PulledTimestamp"] = time.Now().Format(time.RFC3339)

	// Depending on the type, add additional information
	switch val := val.(type) {
	// For types, it can be helpful to look at https://github.com/bluesky-social/indigo/blob/49a1572716a6cccde22022c4264b62acbab43bc2/sonar/sonar.go#L227
	case *bsky.FeedLike:
		// Lookup the actual post (basic author info will be included)
		post, err := h.lookupPost(val.Subject.Uri)
		if err != nil {
			log.Warnf("Failed to get post for like: %s", val.Subject.Uri)
			post = nil
		}
		result["_LikedPost"] = post
	case *bsky.FeedRepost:
		// Lookup the actual post (basic author info will be included)
		post, err := h.lookupPost(val.Subject.Uri)
		if err != nil {
			log.Warnf("Failed to get post for repost: %s", val.Subject.Uri)
			post = nil
		}
		result["_RepostedPost"] = post
	case *bsky.GraphBlock:
		// Lookup the blocked user
		profile, err := h.lookupProfile(val.Subject)
		if err != nil {
			log.Warnf("Failed to get profile for blocked user: %s", val.Subject)
			profile = nil
		}
		result["_BlockedProfile"] = profile
	case *bsky.GraphFollow:
		// Lookup the followed user
		profile, err := h.lookupProfile(val.Subject)
		if err != nil {
			log.Warnf("Failed to get profile for followed user: %s", val.Subject)
			profile = nil
		}
		result["_FollowedProfile"] = profile
	case *bsky.ActorProfile:
		// Nothing to do
	case *bsky.FeedPost:
		// Nothing to do
	}

	return
}
