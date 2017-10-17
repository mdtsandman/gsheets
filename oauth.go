package gsheets

import (
	"net/http"
)

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewOauthClient(secret []byte, api string) (client *http.Client, err error) {
	config, err := google.JWTConfigFromJSON(secret, api)
	if err != nil {
		return client, err
	}
	client = config.Client(oauth2.NoContext)
	if err != nil {
		return client, err
	}
	return client, err
}
