package gowebdav

import (
	"fmt"
	"net/http"
)

type TokenPlacement string

const (
	PlacementHeader TokenPlacement = "header"
	PlacementQuery  TokenPlacement = "query"
)

// TokenAuth structure holds token, its name and placement (header or path)
type TokenAuth struct {
	token     string         // same token
	name      string         // name of token (for example: Bearer, OAuth etc.)
	placement TokenPlacement // "header" or "query"
}

func NewTokenAuth(name, token string, placement TokenPlacement) *TokenAuth {
	return &TokenAuth{
		name:      name,
		token:     token,
		placement: placement,
	}
}

func NewHeaderTokenAuthorizer(name, token string) Authorizer {
	return NewPreemptiveAuth(&TokenAuth{
		name:      name,
		token:     token,
		placement: PlacementHeader,
	})
}

func NewQueryTokenAuthorizer(name, token string) Authorizer {
	return NewPreemptiveAuth(&TokenAuth{
		name:      name,
		token:     token,
		placement: PlacementQuery,
	})
}

// Authorize the current request
func (t *TokenAuth) Authorize(c *http.Client, rq *http.Request, path string) error {
	switch t.placement {
	case "header":
		if t.name != "" {
			rq.Header.Set("Authorization", t.name+" "+t.token)
		} else {
			rq.Header.Set("Authorization", t.token)
		}
	case "query":
		q := rq.URL.Query()

		if t.name == "" {
			return fmt.Errorf("query token name is required")
		}

		q.Set(t.name, t.token)
		rq.URL.RawQuery = q.Encode()
	default:
		return fmt.Errorf("invalid placement: %s.", t.placement)
	}
	return nil
}

// Verify verifies if the authentication
func (t *TokenAuth) Verify(c *http.Client, rs *http.Response, path string) (redo bool, err error) {
	if rs.StatusCode == http.StatusUnauthorized {
		err = NewPathError("Authorize", path, rs.StatusCode)
	}
	return
}

// Close cleans up all resources
func (t *TokenAuth) Close() error {
	return nil
}

// Clone creates a Copy of itself
func (t *TokenAuth) Clone() Authenticator {
	// no copy due to read only access
	return t
}

// String toString
func (t *TokenAuth) String() string {
	return fmt.Sprintf("TokenAuth token name: %s, placement: %s", t.name, t.placement)
}
