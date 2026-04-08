package gowebdav

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func bearerAuthHandler(token string, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer "+token {
			h.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Bearer realm="dav"`)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func TestNewBearerAuth(t *testing.T) {
	a := NewBearerAuth("token").(*BearerAuth)

	ex := "BearerAuth token: true"
	if a.String() != ex {
		t.Error("expected: " + ex + " got: " + a.String())
	}

	if a.Clone() != a {
		t.Error("expected the same instance")
	}

	if a.Close() != nil {
		t.Error("expected close without errors")
	}
}

func TestBearerAuthAuthorize(t *testing.T) {
	a := NewBearerAuth("token")
	rq, _ := http.NewRequest("GET", "http://localhost/", nil)
	if err := a.Authorize(nil, rq, "/"); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if rq.Header.Get("Authorization") != "Bearer token" {
		t.Error("got wrong Authorization header: " + rq.Header.Get("Authorization"))
	}
}

func TestPreemptiveBearerAuth(t *testing.T) {
	auth := NewPreemptiveAuth(NewBearerAuth("token"))
	n, b := auth.NewAuthenticator(nil)
	if b != nil {
		t.Error("expected body to be nil")
	}
	if n == nil {
		t.Fatal("expected authenticator")
	}

	srv, _, _ := newAuthSrv(t, func(h http.Handler) http.HandlerFunc {
		return bearerAuthHandler("token", h)
	})
	defer srv.Close()
	cli := NewAuthClient(srv.URL, auth)
	if err := cli.Connect(); err != nil {
		t.Fatalf("got error: %v, want nil", err)
	}
}

func TestBearerAutoAuthNegotiation(t *testing.T) {
	srv := httptest.NewServer(bearerAuthHandler("token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	defer srv.Close()

	cli := NewClient(srv.URL, "", "token")
	if err := cli.Connect(); err != nil {
		t.Fatalf("got error: %v, want nil", err)
	}
}

func TestNewBearerAuthClient(t *testing.T) {
	srv := httptest.NewServer(bearerAuthHandler("token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	defer srv.Close()

	cli := NewBearerAuthClient(srv.URL, "token")
	if err := cli.Connect(); err != nil {
		t.Fatalf("got error: %v, want nil", err)
	}
}
