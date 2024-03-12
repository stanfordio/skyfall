package auth

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stanfordio/skyfall/pkg/utils"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/xrpc"
)

type Authenticator struct {
	Context context.Context
	Client  *xrpc.Client
}

func MakeAuthenticator(ctx context.Context) (*Authenticator, error) {
	a := Authenticator{
		Context: ctx,
		Client: &xrpc.Client{
			Client: utils.RetryingHTTPClient(),
		},
	}

	return &a, nil
}

func findPersonalDataServerEndpoint(identifier string) (string, error) {
	resolvedIdentifier, error := syntax.ParseAtIdentifier(identifier)
	if error != nil {
		return "", error
	}

	identity, err := identity.DefaultDirectory().Lookup(context.TODO(), *resolvedIdentifier)
	if err != nil {
		return "", err
	}

	// Print the identity
	log.Infof("Found identity: %+v", identity)

	endpoint := identity.PDSEndpoint()

	if endpoint == "" {
		return "", fmt.Errorf("no PDS endpoint found for %s", identifier)
	}

	return endpoint, nil
}

func (a *Authenticator) Authenticate(identifier string, password string) (*xrpc.AuthInfo, error) {
	// First we need to lookup where we authenticate; then we authenticate there
	pdsEndpoint, err := findPersonalDataServerEndpoint(identifier)
	if err != nil {
		return nil, err
	}

	// Print/format the pds endpoint
	log.Infof("Authenticating with PDS endpoint: %s", pdsEndpoint)

	// Set the host to the PDS endpoint
	a.Client.Host = pdsEndpoint

	// Hit the PDS endpoint to authenticate
	output, err := atproto.ServerCreateSession(a.Context, a.Client, &atproto.ServerCreateSession_Input{
		Identifier: identifier,
		Password:   password,
	})

	if err != nil {
		return nil, err
	}

	info := xrpc.AuthInfo{
		AccessJwt:  output.AccessJwt,
		RefreshJwt: output.RefreshJwt,
		Did:        output.Did,
		Handle:     output.Handle,
	}

	a.Client.Auth = &info
	go a.refreshTokenContinuously(output)
	return &info, nil
}

func (a *Authenticator) refreshTokenContinuously(authOutput *atproto.ServerCreateSession_Output) {
	// Put the refresh token into the access token slot. Janky, but this is what Bluesky expects.
	// We intentionally create a new AuthInfo here.
	a.Client.Auth = &xrpc.AuthInfo{
		AccessJwt: authOutput.RefreshJwt,
		Did:       authOutput.Did,
		Handle:    authOutput.Handle,
	}

	// Send a refresh request every minute
	for {
		_, error := atproto.ServerRefreshSession(a.Context, a.Client)

		if error != nil {
			log.Errorf("Error refreshing token: %+v", error)
		} else {
			log.Debugf("Successfully refreshed access token")
		}

		time.Sleep(15 * time.Second)
	}
}
