package gowebdav

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func TestParseLineMissingValue(t *testing.T) {
	login, pass := parseLine("machine example.com login")
	if login != "" || pass != "" {
		t.Fatalf("got login=%q pass=%q, want both empty", login, pass)
	}
}

func TestReadConfigMatchesHostnameWithoutPort(t *testing.T) {
	dir := t.TempDir()
	netrc := filepath.Join(dir, ".netrc")
	content := "machine other.example login other password nope\n" +
		"machine example.com\n" +
		"  login demo\n" +
		"  password secret\n"
	if err := os.WriteFile(netrc, []byte(content), 0600); err != nil {
		t.Fatalf("write netrc: %v", err)
	}

	login, pass := ReadConfig("https://example.com:8443/webdav", netrc)
	if login != "demo" || pass != "secret" {
		t.Fatalf("got login=%q pass=%q, want demo/secret", login, pass)
	}
}

func TestStatInvalidXMLReturnsError(t *testing.T) {
	srv := httptest.NewServer(basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PROPFIND":
			w.WriteHeader(207)
			_, _ = io.WriteString(w, `<d:multistatus xmlns:d="DAV:"><d:response>`)
		default:
			w.WriteHeader(200)
		}
	})))
	defer srv.Close()

	cli := NewClient(srv.URL, "user", "password")
	info, err := cli.Stat("/broken.txt")
	if err == nil {
		t.Fatalf("got info=%v, want xml parse error", info)
	}

	pe, ok := err.(*os.PathError)
	if !ok {
		t.Fatalf("got %T, want *os.PathError", err)
	}
	if pe.Op != "Stat" {
		t.Fatalf("got op=%q, want Stat", pe.Op)
	}
}

func TestReadStreamRangeFallbackWithoutServerRangeSupport(t *testing.T) {
	content := []byte("hello gowebdav\n")
	srv := httptest.NewServer(basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write(content)
	})))
	defer srv.Close()

	cli := NewClient(srv.URL, "user", "password")

	stream, err := cli.ReadStreamRange("/hello.txt", 4, 4)
	if err != nil {
		t.Fatalf("range read with length: %v", err)
	}
	defer stream.Close()

	got, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("read limited range: %v", err)
	}
	if !bytes.Equal(got, []byte("o go")) {
		t.Fatalf("got %q, want %q", got, "o go")
	}

	stream, err = cli.ReadStreamRange("/hello.txt", 6, 0)
	if err != nil {
		t.Fatalf("range read to end: %v", err)
	}
	defer stream.Close()

	got, err = io.ReadAll(stream)
	if err != nil {
		t.Fatalf("read tail range: %v", err)
	}
	if !bytes.Equal(got, []byte("gowebdav\n")) {
		t.Fatalf("got %q, want %q", got, "gowebdav\n")
	}
}

func TestRemoveAllPreservesUnderlyingError(t *testing.T) {
	cli := NewClient("https://example.com", "user", "password")
	sentinel := errors.New("network down")
	cli.SetTransport(roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, sentinel
	}))

	err := cli.RemoveAll("/hello.txt")
	if !errors.Is(err, sentinel) {
		t.Fatalf("got %v, want wrapped %v", err, sentinel)
	}
}
