package gsheets

import (
	"net/http"
)

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewOauthClient(secret []byte, api string, ctx context.Context) (client *http.Client, err error) {
	config, err := google.JWTConfigFromJSON(secret, api)
	if err != nil {
		return client, err
	}
	if ctx == nil {
		client = config.Client(oauth2.NoContext)
	} else {
		client = config.Client(ctx)
	}
	return client, nil
}
