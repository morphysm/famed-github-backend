package providers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/awnumar/memguard"
	"github.com/golang-jwt/jwt/v4"
)

const (
	acceptHeader = "application/vnd.github.v3+json"
)

// Copied from "github.com/bradleyfalzon/ghinstallation/v2"
// Adapted to use "github.com/awnumar/memguard"

// AppsTransport provides a http.RoundTripper by wrapping an existing
// http.RoundTripper and provides GitHub Apps authentication as a
// GitHub App.
//
// Client can also be overwritten, and is useful to change to one which
// provides retry logic if you do experience retryable errors.
//
// See https://developer.github.com/apps/building-integrations/setting-up-and-registering-github-apps/about-authentication-options-for-github-apps/
type AppsTransport struct {
	BaseURL    string            // BaseURL is the scheme and host for GitHub API, defaults to https://api.github.com
	Client     *http.Client      // Client to use to refresh tokens, defaults to http.Client with provided transport
	tr         http.RoundTripper // tr is the underlying roundtripper being wrapped
	keyEnclave *memguard.Enclave // keyEnclave memguard.Enclave containing is the GitHub App's private key
	appID      int64             // appID is the GitHub App's ID
}

// NewAppsTransport returns an AppsTransport using a memguard.Enclave containing a crypto/rsa.(*PrivateKey).
func NewAppsTransport(baseUrl string, tr http.RoundTripper, appID int64, key *memguard.Enclave) *AppsTransport {
	return &AppsTransport{
		BaseURL:    baseUrl,
		Client:     &http.Client{Transport: tr},
		tr:         tr,
		keyEnclave: key,
		appID:      appID,
	}
}

// RoundTrip implements http.RoundTripper interface.
func (t *AppsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// GitHub rejects expiry and issue timestamps that are not an integer,
	// while the jwt-go library serializes to fractional timestamps.
	// Truncate them before passing to jwt-go.
	iss := time.Now().Add(-30 * time.Second).Truncate(time.Second)
	exp := iss.Add(2 * time.Minute)
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(iss),
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    strconv.FormatInt(t.appID, 10),
	}
	bearer := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Decrypt the key into a local copy
	keyBuf, err := t.keyEnclave.Open()
	if err != nil {
		memguard.SafePanic(err)
	}
	defer keyBuf.Destroy() // Destroy the copy when we return

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBuf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("could not parse private keyEnclave: %s", err)
	}

	ss, err := bearer.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("could not sign jwt: %s", err)
	}

	req.Header.Set("Authorization", "Bearer "+ss)
	req.Header.Add("Accept", acceptHeader)

	resp, err := t.tr.RoundTrip(req)
	return resp, err
}
