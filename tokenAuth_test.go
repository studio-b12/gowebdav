package gowebdav

import (
	"net/http"
	"testing"
)

func TestNewTokenAuth(t *testing.T) {
	a := &TokenAuth{
		token:     "secret-token",
		name:      "Bearer",
		placement: PlacementHeader,
	}

	ex := "TokenAuth token name: Bearer, placement: header"
	if a.String() != ex {
		t.Errorf("expected %q, got %q", ex, a.String())
	}

	if a.Clone() != a {
		t.Error("expected the same instance")
	}

	if err := a.Close(); err != nil {
		t.Error("expected close without errors")
	}
}

func TestTokenAuthAuthorizeHeaderWithName(t *testing.T) {
	a := &TokenAuth{
		token:     "abc123",
		name:      "Bearer",
		placement: PlacementHeader,
	}

	rq, _ := http.NewRequest("GET", "http://localhost/", nil)

	if err := a.Authorize(nil, rq, "/"); err != nil {
		t.Fatalf("Authorize returned error: %v", err)
	}

	want := "Bearer abc123"
	got := rq.Header.Get("Authorization")

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestTokenAuthAuthorizeHeaderWithoutName(t *testing.T) {
	a := &TokenAuth{
		token:     "abc123",
		placement: PlacementHeader,
	}

	rq, _ := http.NewRequest("GET", "http://localhost/", nil)

	if err := a.Authorize(nil, rq, "/"); err != nil {
		t.Fatalf("Authorize returned error: %v", err)
	}

	want := "abc123"
	got := rq.Header.Get("Authorization")

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestTokenAuthAuthorizeQuery(t *testing.T) {
	a := &TokenAuth{
		token:     "abc123",
		name:      "access_token",
		placement: PlacementQuery,
	}

	rq, _ := http.NewRequest("GET", "http://localhost/?page=1", nil)

	if err := a.Authorize(nil, rq, "/"); err != nil {
		t.Fatalf("Authorize returned error: %v", err)
	}

	got := rq.URL.Query().Get("access_token")
	if got != "abc123" {
		t.Errorf("expected query token %q, got %q", "abc123", got)
	}

	if rq.URL.Query().Get("page") != "1" {
		t.Error("existing query parameters were lost")
	}
}

func TestTokenAuthAuthorizeQueryWithoutName(t *testing.T) {
	a := &TokenAuth{
		token:     "abc123",
		placement: PlacementQuery,
	}

	rq, _ := http.NewRequest("GET", "http://localhost/", nil)

	if err := a.Authorize(nil, rq, "/"); err == nil {
		t.Fatal("expected error")
	}
}

func TestTokenAuthAuthorizeInvalidPlacement(t *testing.T) {
	a := &TokenAuth{
		token:     "abc123",
		name:      "Bearer",
		placement: "invalid",
	}

	rq, _ := http.NewRequest("GET", "http://localhost/", nil)

	if err := a.Authorize(nil, rq, "/"); err == nil {
		t.Fatal("expected error")
	}
}

func TestTokenAuthVerifyUnauthorized(t *testing.T) {
	a := &TokenAuth{}

	rs := &http.Response{
		StatusCode: http.StatusUnauthorized,
	}

	redo, err := a.Verify(nil, rs, "/")

	if redo {
		t.Error("expected redo=false")
	}

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTokenAuthVerifyOK(t *testing.T) {
	a := &TokenAuth{}

	rs := &http.Response{
		StatusCode: http.StatusOK,
	}

	redo, err := a.Verify(nil, rs, "/")

	if redo {
		t.Error("expected redo=false")
	}

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
