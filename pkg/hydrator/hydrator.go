package hydrator

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/DmitriyVTitov/size"
	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"
	"github.com/stanfordio/skyfall/pkg/utils"
	"go.uber.org/ratelimit"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/atproto/crypto"
	"github.com/bluesky-social/indigo/atproto/identity"
	atpidentity "github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/repo"
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
	Ratelimit         ratelimit.Limiter // Rate limiting for authenticated endpoints. May be called by other packages whenever they make a rate-limited request.
}

var didRegex = regexp.MustCompile(`did:plc:[a-zA-Z0-9]+`)

func MakeHydrator(ctx context.Context, cacheSize int64, authInfo *xrpc.AuthInfo) (*Hydrator, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e8, // number of keys to track frequency of
		MaxCost:     cacheSize,
		BufferItems: 64, // number of keys per Get buffer
		Cost: func(value interface{}) int64 {
			val := int64(size.Of(value))
			return val
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %+v", err)
	}

	h := Hydrator{
		Cache:   cache,
		Context: ctx,
		Client: &xrpc.Client{
			Client: utils.RetryingHTTPClient(),
			Host:   "https://public.api.bsky.app", // We generally want to use the public.api.bsky.app host for all requests, since they are doing the indexing (and it's public with big rate limits)
			Auth:   authInfo,
		},
		IdentityDirectory: atpidentity.DefaultDirectory(),
		AuthInfo:          authInfo,
		Ratelimit:         ratelimit.New(1000), // 1000 requests per second; empirically we find that this is fine
	}

	return &h, nil
}

func namespaceKey(namespace string, key string) string {
	return fmt.Sprintf("%s:%s", namespace, key)
}

func (h *Hydrator) LookupIdentity(identifier string) (identity *atpidentity.Identity, err error) {
	key := namespaceKey("identity", identifier)

	// Check the cache first
	cachedValue, found := h.Cache.Get(key)

	if found && cachedValue != nil {
		if cachedError, isErr := cachedValue.(error); isErr {
			log.Debugf("Cached error for %s: %v", identifier, cachedError)
			return nil, cachedError // Return the cached error
		}
		identity = cachedValue.(*atpidentity.Identity)
		return
	}

	log.Debugf("Cache miss for %s", identifier)

	h.Ratelimit.Take()
	resolvedIdentifier, error := syntax.ParseAtIdentifier(identifier)
	if error != nil {
		err = error
		return
	}

	identity, err = h.IdentityDirectory.Lookup(h.Context, *resolvedIdentifier)
	if err != nil {
		h.Cache.SetWithTTL(key, err, 0, time.Duration(1)*time.Hour*24)
		return
	}

	h.Cache.SetWithTTL(key, identity, 0, time.Duration(1)*time.Hour*24)

	return
}

func (h *Hydrator) lookupProfileFromIdentity(identity *atpidentity.Identity) (profile *bsky.ActorDefs_ProfileViewDetailed, err error) {
	if identity == nil {
		return nil, fmt.Errorf("identity is nil")
	}

	key := namespaceKey("profile", identity.Handle.String())

	// Check the cache first
	cachedValue, found := h.Cache.Get(key)

	if found && cachedValue != nil {
		if cachedError, isErr := cachedValue.(error); isErr {
			log.Warnf("Found cached error for %s: %v", identity.Handle.String(), cachedError)
			return nil, cachedError
		}
		profile = cachedValue.(*bsky.ActorDefs_ProfileViewDetailed)
		return
	}

	h.Ratelimit.Take()
	profile, err = bsky.ActorGetProfile(h.Context, h.Client, identity.Handle.String())

	if err != nil { // Cache if error looking up profile like suspended
		log.Warnf("Profile lookup failed for identity %s: %s", identity.Handle.String(), err)
		h.Cache.SetWithTTL(key, err, 0, time.Duration(1)*time.Hour*24)
		return nil, err
	}

	h.Cache.SetWithTTL(key, profile, 0, time.Duration(1)*time.Hour*24)

	return profile, nil
}

func (h *Hydrator) lookupProfile(did string) (profile *bsky.ActorDefs_ProfileViewDetailed, err error) {
	identity, err := h.LookupIdentity(did)
	if err != nil {
		return
	}
	profile, err = h.lookupProfileFromIdentity(identity)
	return
}

func (h *Hydrator) lookupPost(atUrl string) (post *bsky.FeedDefs_PostView, err error) {
	key := namespaceKey("post", atUrl)

	// Check the cache first
	cachedValue, found := h.Cache.Get(key)

	if found && cachedValue != nil {
		if cachedError, isErr := cachedValue.(error); isErr {
			log.Warnf("Cached error for %s: %v", atUrl, cachedError)
			return nil, cachedError
		}
		post = cachedValue.(*bsky.FeedDefs_PostView)
		return
	}

	log.Debugf("Cache miss for %s", atUrl)

	h.Ratelimit.Take()
	output, err := bsky.FeedGetPosts(h.Context, h.Client, []string{atUrl})

	if err != nil {
		log.Errorf("Unable to fetch post at %s: %s", atUrl, err)
		return
	}

	if len(output.Posts) == 0 {
		err = fmt.Errorf("no posts found for %s", atUrl)
	}

	if err != nil { // caching miss so we don't keep checking
		h.Cache.SetWithTTL(key, err, 0, time.Duration(1)*time.Hour*24)
		return
	}

	post = output.Posts[0]

	h.Cache.SetWithTTL(key, post, 0, time.Duration(1)*time.Hour*24)

	return
}

func (h *Hydrator) flattenIdentity(identity *atpidentity.Identity) (result map[string]interface{}, err error) {
	if identity == nil {
		return nil, fmt.Errorf("identity is nil")
	}

	result = make(map[string]interface{})

	result["DID"] = identity.Handle.String()
	result["Handle"] = identity.Handle

	var pk crypto.PublicKey
	pk, _ = identity.PublicKey()
	if err != nil {
		log.Warnf("Failed to get public key for actor: %s, %s", identity.Handle, err)
	} else {

		result["DIDKey"] = pk.DIDKey()
	}

	result["PDS"] = identity.PDSEndpoint()

	return
}

func (h *Hydrator) flattenProfile(profile *bsky.ActorDefs_ProfileViewBasic) (result map[string]interface{}) {
	if profile == nil {
		return nil
	}

	result = make(map[string]interface{})

	result["Avatar"] = profile.Avatar
	result["DisplayName"] = profile.DisplayName
	result["Handle"] = profile.Handle
	result["DID"] = profile.Did

	return
}

func (h *Hydrator) flattenActorProfile(profile *bsky.ActorProfile) (result map[string]interface{}) {
	if profile == nil {
		return nil
	}

	result = make(map[string]interface{})

	result["DisplayName"] = profile.DisplayName
	result["Description"] = profile.Description

	return
}

func (h *Hydrator) flattenFullProfile(profile *bsky.ActorDefs_ProfileViewDetailed) (result map[string]interface{}) {
	if profile == nil {
		return nil
	}

	result = make(map[string]interface{})

	result["Avatar"] = profile.Avatar
	result["DisplayName"] = profile.DisplayName
	result["Handle"] = profile.Handle
	result["DID"] = profile.Did
	result["Description"] = profile.Description
	result["FollowersCount"] = profile.FollowersCount
	result["FollowsCount"] = profile.FollowsCount
	result["PostsCount"] = profile.PostsCount
	result["IndexedAt"] = profile.IndexedAt

	return
}

func (h *Hydrator) flattenFacets(facets []*bsky.RichtextFacet) (hashtags []string, urls []string) {
	hashtags = []string{}
	urls = []string{}
	if facets != nil {
		for _, facet := range facets {
			if facet != nil {
				features := facet.Features
				for _, feature := range features {
					if feature.RichtextFacet_Tag != nil {
						tag := feature.RichtextFacet_Tag.Tag
						hashtags = append(hashtags, tag)
					}
					if feature.RichtextFacet_Link != nil {
						url := feature.RichtextFacet_Link.Uri
						urls = append(urls, url)
					}
				}
			}
		}
	}
	return
}

func (h *Hydrator) flattenPostView(post *bsky.FeedDefs_PostView) (result map[string]interface{}) {
	if post == nil {
		return nil
	}

	result = make(map[string]interface{})

	if post.Author == nil {
		result["Author"] = nil
	} else {
		result["Author"] = h.flattenProfile(post.Author)
	}
	result["CID"] = post.Cid
	result["LikeCount"] = post.LikeCount
	result["RepostCount"] = post.RepostCount
	result["ReplyCount"] = post.ReplyCount
	result["LikeCount"] = post.LikeCount
	result["URI"] = post.Uri

	rec := post.Record.Val.(*bsky.FeedPost)
	result["Text"] = rec.Text
	result["CreatedAt"] = rec.CreatedAt

	result["Langs"] = rec.Langs
	if rec.Langs == nil {
		result["Langs"] = []string{}
	}

	if rec.Embed != nil {
		result["Embed"] = h.flattenEmbed(rec.Embed)
	}

	hashtags, urls := h.flattenFacets(rec.Facets)
	result["Hashtags"] = hashtags
	result["URLs"] = urls

	return
}

func (h *Hydrator) flattenPost(post *bsky.FeedPost) (result map[string]interface{}) {
	if post == nil {
		return nil
	}

	result = make(map[string]interface{})

	result["Text"] = post.Text
	result["CreatedAt"] = post.CreatedAt

	result["Langs"] = post.Langs
	if post.Langs == nil {
		result["Langs"] = []string{}
	}

	if post.Reply != nil {
		// For some reason replies can lack a parent
		if post.Reply.Parent != nil {
			result["ReplyParentCID"] = post.Reply.Parent.Cid
		}
	}

	if post.Embed != nil {
		result["Embed"] = h.flattenEmbed(post.Embed)
	}

	hashtags, urls := h.flattenFacets(post.Facets)
	result["Hashtags"] = hashtags
	result["URLs"] = urls

	return
}

func (h *Hydrator) extractAllDids(str string) []string {
	// Nasty hack that runs a Regex over the string to extract all the DIDs. We
	// do this so we can reliably extract all DIDs from records without having
	// to worry about Bluesky changing their schema.

	// Find all matches of the regular expression in the string
	matches := didRegex.FindAllString(str, -1)

	// Return the slice of extracted DIDs
	return matches
}

func (h *Hydrator) GetIdentitiesInRepo(repo *repo.Repo) ([]atpidentity.Identity, error) {
	identities := make([]atpidentity.Identity, 0)
	identitiesFound := make(map[string]bool)

	err := repo.ForEach(h.Context, "", func(k string, v cid.Cid) error {
		// Get the record
		_, rec, err := repo.GetRecord(h.Context, k)
		if err != nil {
			log.Errorf("Unable to parse CID %s: %s", v.String(), err)
			return err
		}

		recJson, err := json.MarshalIndent(rec, "", "  ")
		if err != nil {
			return err
		}
		dids := h.extractAllDids(string(recJson))

		// Lookup all the identities
		for i := range dids {
			did := dids[i]
			if identitiesFound[did] {
				continue
			}
			identitiesFound[did] = true
			identity, err := h.LookupIdentity(did)
			if err != nil {
				log.Errorf("Failed to lookup identity for %s: %s", did, err)
				continue
			}
			identities = append(identities, *identity)
		}

		return nil
	})

	return identities, err
}

func (h *Hydrator) GetRepoBytes(actorDid string, pdsEndpoint string) ([]byte, error) {
	key := namespaceKey("repo", actorDid)

	// Check the cache first
	cachedValue, found := h.Cache.Get(key)

	if found && cachedValue != nil {
		repo := cachedValue.([]byte)
		return repo, nil
	}

	xrpcc := xrpc.Client{
		Host: pdsEndpoint,
	}

	h.Ratelimit.Take()
	repoBytes, err := atproto.SyncGetRepo(h.Context, &xrpcc, actorDid, "")
	if err != nil {
		return nil, err
	}

	// Set the cache
	h.Cache.SetWithTTL(key, repoBytes, 0, time.Duration(1)*time.Hour*24)

	return repoBytes, nil

}

func (h *Hydrator) invalidateRepoCache(actorDid string) {
	h.Cache.Del(actorDid)
}

func (h *Hydrator) flattenEmbed(embed *bsky.FeedPost_Embed) (result map[string]interface{}) {
	if embed == nil {
		return nil
	}

	result = make(map[string]interface{})

	// Three types of embeds: external links, images, records, and records with media
	if embed.EmbedExternal != nil && embed.EmbedExternal.External != nil {
		externalEmbedResult := make(map[string]interface{})
		externalEmbedResult["URI"] = embed.EmbedExternal.External.Uri
		externalEmbedResult["Title"] = embed.EmbedExternal.External.Title
		externalEmbedResult["Description"] = embed.EmbedExternal.External.Description
		result["External"] = externalEmbedResult
	} else {
		result["External"] = nil
	}

	if embed.EmbedImages != nil && len(embed.EmbedImages.Images) > 0 {
		images := make([]map[string]interface{}, 0)
		for _, image := range embed.EmbedImages.Images {
			imageResult := make(map[string]interface{})
			imageResult["Alt"] = image.Alt
			if image.Image != nil {
				imageResult["BlobLink"] = image.Image.Ref.String()
				imageResult["MimeType"] = image.Image.MimeType
				if image.AspectRatio != nil {
					imageResult["Width"] = image.AspectRatio.Width
					imageResult["Height"] = image.AspectRatio.Height
				}
				imageResult["MimeType"] = image.Image.MimeType
			}
			images = append(images, imageResult)
		}
		result["Images"] = images
	} else {
		result["Images"] = []map[string]interface{}{}
	}

	if embed.EmbedRecord != nil && embed.EmbedRecord.Record != nil {
		recordEmbedResult := make(map[string]interface{})
		recordEmbedResult["CID"] = embed.EmbedRecord.Record.Cid
		recordEmbedResult["URI"] = embed.EmbedRecord.Record.Uri
		recordEmbedResult["Type"] = embed.EmbedRecord.LexiconTypeID
		result["Record"] = recordEmbedResult
	}

	if embed.EmbedRecordWithMedia != nil && embed.EmbedRecordWithMedia.Record != nil {
		recordEmbedResult := make(map[string]interface{})
		if embed.EmbedRecordWithMedia.Record.Record != nil {
			recordEmbedResult["CID"] = embed.EmbedRecordWithMedia.Record.Record.Cid
			recordEmbedResult["URI"] = embed.EmbedRecordWithMedia.Record.Record.Uri
		}
		recordEmbedResult["Type"] = embed.EmbedRecordWithMedia.LexiconTypeID

		media := make([]map[string]interface{}, 0)
		if embed.EmbedRecordWithMedia.Media != nil {
			if embed.EmbedRecordWithMedia.Media.EmbedImages != nil {
				for _, image := range embed.EmbedRecordWithMedia.Media.EmbedImages.Images {
					mediaResult := make(map[string]interface{})
					mediaResult["Alt"] = image.Alt
					if image.Image != nil {
						mediaResult["BlobLink"] = image.Image.Ref.String()
						mediaResult["MimeType"] = image.Image.MimeType
					}

					if image.AspectRatio != nil {
						mediaResult["Width"] = image.AspectRatio.Width
						mediaResult["Height"] = image.AspectRatio.Height
					}

					media = append(media, mediaResult)
				}
			}
		}
		result["EmbedRecordMedia"] = media

		result["Record"] = recordEmbedResult
	}

	return
}

func (h *Hydrator) Hydrate(val interface{}, actorDid string) (result map[string]interface{}, err error) {
	err = nil

	result = make(map[string]interface{})
	full := make(map[string]interface{})
	projection := make(map[string]interface{})
	err = mapstructure.Decode(val, &full)

	if err != nil {
		return
	}

	// Resolve full identity and profile information for the actor
	identity, err := h.LookupIdentity(actorDid)
	if err != nil {
		log.Warnf("Failed to lookup identity for actor %s: %s", actorDid, err)
		identity = nil
	}

	profile, err := h.lookupProfileFromIdentity(identity)
	if err != nil {
		log.Warnf("Failed to lookup profile for actor %s: %s", actorDid, err)
		profile = nil
	}

	// Add key metadata to the outer map
	result["Type"] = full["LexiconTypeID"]
	result["CreatedAt"] = full["CreatedAt"]
	result["PulledTimestamp"] = time.Now().Format(time.RFC3339)

	// Add the actorDid and profile to the map
	full["_ActorDid"] = actorDid
	full["_ActorIdentity"] = identity
	full["_ActorProfile"] = profile
	flat, err := h.flattenIdentity(identity)
	if err != nil {
		log.Warnf("Failed to flatten identity %s: %s", actorDid, err)
		flat = nil
	}
	projection["Actor"] = flat

	// Depending on the type, add additional information
	switch val := val.(type) {
	// For types, it can be helpful to look at https://github.com/bluesky-social/indigo/blob/49a1572716a6cccde22022c4264b62acbab43bc2/sonar/sonar.go#L227
	case *bsky.FeedLike:
		// Lookup the actual post (basic author info will be included)
		if val.Subject != nil {
			post, err := h.lookupPost(val.Subject.Uri)
			if err != nil {
				log.Warnf("Failed to get post for like: %s", val.Subject.Uri)
				post = nil
			}
			full["_LikedPost"] = post
			projection["LikedPost"] = h.flattenPostView(post)
		} else {
			log.Warn("No Subject in Like")
		}
	case *bsky.FeedRepost:
		// Lookup the actual post (basic author info will be included)
		if val.Subject != nil {
			post, err := h.lookupPost(val.Subject.Uri)
			if err != nil {
				log.Warnf("Failed to get post for repost: %s", val.Subject.Uri)
				post = nil
			}
			full["_RepostedPost"] = post
			projection["RepostedPost"] = h.flattenPostView(post)
		} else {
			log.Warn("No Subject in Repost")
		}
	case *bsky.GraphBlock:
		// Lookup the blocked user
		profile, err := h.lookupProfile(val.Subject)
		if err != nil {
			log.Warnf("Failed to get profile for blocked user: %s", val.Subject)
			profile = nil
		}
		full["_BlockedProfile"] = profile
		projection["BlockedProfile"] = h.flattenFullProfile(profile)
	case *bsky.GraphFollow:
		// Lookup the followed user
		profile, err := h.lookupProfile(val.Subject)
		if err != nil {
			log.Warnf("Failed to get profile for followed user: %s", val.Subject)
			profile = nil
		}
		full["_FollowedProfile"] = profile
		projection["FollowedProfile"] = h.flattenFullProfile(profile)
	case *bsky.ActorProfile:
		projection["Profile"] = h.flattenActorProfile(val)
	case *bsky.FeedPost:
		projection["Post"] = h.flattenPost(val)
	}

	// Add the full object to the result
	result["Full"] = full
	result["Projection"] = projection

	return
}
