package auth

import (
	"crypto/sha256"
	"net/http"

	"github.com/tent/hawk-go"
)

var DefaultHawkHasher = sha256.New

type HawkRoundTripper struct {
	User  string
	Token string
}

func (t *HawkRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	creds := &hawk.Credentials{
		ID:   t.User,
		Key:  t.Token,
		Hash: DefaultHawkHasher,
	}

	auth := hawk.NewRequestAuth(req, creds, 0)

	req.Header.Set("Authorization", auth.RequestHeader())
	return http.DefaultTransport.RoundTrip(req)
}
