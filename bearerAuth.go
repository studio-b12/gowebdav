package gowebdav

import (
	"fmt"
	"net/http"
)

// BearerAuth structure holds our bearer token
type BearerAuth struct {
	token string
}

// NewBearerAuth creates a new bearer token authenticator.
func NewBearerAuth(token string) Authenticator {
	return &BearerAuth{token: token}
}

// NewBearerAuthClient creates a new client instance with preemptive bearer authentication.
func NewBearerAuthClient(uri, token string) *Client {
	return NewAuthClient(uri, NewPreemptiveAuth(NewBearerAuth(token)))
}

// Authorize the current request
func (b *BearerAuth) Authorize(c *http.Client, rq *http.Request, path string) error {
	rq.Header.Set("Authorization", "Bearer "+b.token)
	return nil
}

// Verify verifies if the authentication
func (b *BearerAuth) Verify(c *http.Client, rs *http.Response, path string) (redo bool, err error) {
	if rs.StatusCode == http.StatusUnauthorized {
		err = NewPathError("Authorize", path, rs.StatusCode)
	}
	return
}

// Close cleans up all resources
func (b *BearerAuth) Close() error {
	return nil
}

// Clone creates a Copy of itself
func (b *BearerAuth) Clone() Authenticator {
	return b
}

// String toString
func (b *BearerAuth) String() string {
	return fmt.Sprintf("BearerAuth token: %t", b.token != "")
}
