package gowebdav

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDigestPartsProtocolCompatibility(t *testing.T) {
	rs := &http.Response{
		Header: http.Header{
			"Www-Authenticate": []string{`Digest realm="testrealm@host.com", qop="auth-int, auth", nonce="abc123", algorithm=md5-sess, opaque="xyz"`},
		},
	}

	parts := digestParts(rs)
	if parts["qop"] != "auth" {
		t.Fatalf("got qop=%q, want auth", parts["qop"])
	}
	if parts["algorithm"] != "MD5-sess" {
		t.Fatalf("got algorithm=%q, want MD5-sess", parts["algorithm"])
	}
	if parts["nonce"] != "abc123" {
		t.Fatalf("got nonce=%q, want abc123", parts["nonce"])
	}
}

func TestConnectDigestAuthProtocolCompatibility(t *testing.T) {
	const (
		realm  = "testrealm@host.com"
		nonce  = "dcd98b7102dd2f0e8b11d0f600bfb0c093"
		opaque = "5ccc069c403ebaf9f0171e9517f40e41"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Digest ") {
			parts := parseDigestAuthorizationHeader(t, auth)
			expectedHA1 := getMD5(fmt.Sprintf("%s:%s:%s",
				getMD5("digestUser:"+realm+":digestPW"),
				nonce,
				parts["cnonce"],
			))
			expectedHA2 := getMD5(r.Method + ":" + parts["uri"])
			expectedResponse := getMD5(fmt.Sprintf("%s:%s:%s:%s:%s:%s",
				expectedHA1,
				nonce,
				parts["nc"],
				parts["cnonce"],
				"auth",
				expectedHA2,
			))

			if parts["qop"] == "auth" && parts["response"] == expectedResponse && parts["algorithm"] == "MD5-sess" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		w.Header().Set("Www-Authenticate", `Digest realm="`+realm+`", qop="auth-int,auth", nonce="`+nonce+`", algorithm=md5-sess, opaque="`+opaque+`"`)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	cli := NewClient(srv.URL, "digestUser", "digestPW")
	if err := cli.Connect(); err != nil {
		t.Fatalf("connect with digest challenge variants: %v", err)
	}
}

func TestProtocolConnectSendsOptionsDepthZero(t *testing.T) {
	srv := httptest.NewServer(basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			t.Fatalf("got method=%s, want OPTIONS", r.Method)
		}
		if got := r.Header.Get("Depth"); got != "0" {
			t.Fatalf("got Depth=%q, want 0", got)
		}
		w.WriteHeader(http.StatusOK)
	})))
	defer srv.Close()

	cli := NewClient(srv.URL, "user", "password")
	if err := cli.Connect(); err != nil {
		t.Fatalf("connect: %v", err)
	}
}

func TestProtocolPropfindHeaders(t *testing.T) {
	srv := httptest.NewServer(basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PROPFIND" {
			t.Fatalf("got method=%s, want PROPFIND", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/xml;charset=UTF-8" {
			t.Fatalf("got Content-Type=%q", got)
		}
		if got := r.Header.Get("Accept"); got != "application/xml,text/xml" {
			t.Fatalf("got Accept=%q", got)
		}

		w.WriteHeader(207)
		switch r.URL.Path {
		case "/hello.txt":
			if got := r.Header.Get("Depth"); got != "0" {
				t.Fatalf("stat got Depth=%q, want 0", got)
			}
			_, _ = io.WriteString(w, `<?xml version="1.0"?>
<d:multistatus xmlns:d="DAV:">
  <d:response>
    <d:href>/hello.txt</d:href>
    <d:propstat>
      <d:status>HTTP/1.1 200 OK</d:status>
      <d:prop>
        <d:displayname>hello.txt</d:displayname>
        <d:getcontentlength>15</d:getcontentlength>
        <d:getlastmodified>Mon, 02 Jan 2006 15:04:05 GMT</d:getlastmodified>
      </d:prop>
    </d:propstat>
  </d:response>
</d:multistatus>`)
		case "/":
			if got := r.Header.Get("Depth"); got != "1" {
				t.Fatalf("readdir got Depth=%q, want 1", got)
			}
			_, _ = io.WriteString(w, `<?xml version="1.0"?>
<d:multistatus xmlns:d="DAV:">
  <d:response>
    <d:href>/</d:href>
    <d:propstat>
      <d:status>HTTP/1.1 200 OK</d:status>
      <d:prop>
        <d:displayname>/</d:displayname>
        <d:resourcetype><d:collection/></d:resourcetype>
        <d:getlastmodified>Mon, 02 Jan 2006 15:04:05 GMT</d:getlastmodified>
      </d:prop>
    </d:propstat>
  </d:response>
  <d:response>
    <d:href>/hello.txt</d:href>
    <d:propstat>
      <d:status>HTTP/1.1 200 OK</d:status>
      <d:prop>
        <d:displayname>hello.txt</d:displayname>
        <d:getcontentlength>15</d:getcontentlength>
        <d:getlastmodified>Mon, 02 Jan 2006 15:04:05 GMT</d:getlastmodified>
      </d:prop>
    </d:propstat>
  </d:response>
</d:multistatus>`)
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})))
	defer srv.Close()

	cli := NewClient(srv.URL, "user", "password")
	if _, err := cli.Stat("/hello.txt"); err != nil {
		t.Fatalf("stat: %v", err)
	}
	if _, err := cli.ReadDir("/"); err != nil {
		t.Fatalf("readdir: %v", err)
	}
}

func TestProtocolCopySendsDestinationAndOverwriteHeaders(t *testing.T) {
	var (
		gotDestination string
		gotOverwrite   string
	)

	srv := httptest.NewServer(basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "COPY" {
			t.Fatalf("got method=%s, want COPY", r.Method)
		}
		gotDestination = r.Header.Get("Destination")
		gotOverwrite = r.Header.Get("Overwrite")
		w.WriteHeader(http.StatusCreated)
	})))
	defer srv.Close()

	cli := NewClient(srv.URL, "user", "password")
	if err := cli.Copy("/src file.txt", "/dst dir/file #1.txt", false); err != nil {
		t.Fatalf("copy: %v", err)
	}

	wantDestination := PathEscape(Join(FixSlash(srv.URL), "/dst dir/file #1.txt"))
	if gotDestination != wantDestination {
		t.Fatalf("got Destination=%q, want %q", gotDestination, wantDestination)
	}
	if gotOverwrite != "F" {
		t.Fatalf("got Overwrite=%q, want F", gotOverwrite)
	}
}

func TestProtocolReadStreamRangeSendsRangeHeader(t *testing.T) {
	srv := httptest.NewServer(basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("got method=%s, want GET", r.Method)
		}
		if got := r.Header.Get("Range"); got != "bytes=4-7" {
			t.Fatalf("got Range=%q, want bytes=4-7", got)
		}
		w.WriteHeader(http.StatusPartialContent)
		_, _ = io.WriteString(w, "o go")
	})))
	defer srv.Close()

	cli := NewClient(srv.URL, "user", "password")
	stream, err := cli.ReadStreamRange("/hello.txt", 4, 4)
	if err != nil {
		t.Fatalf("range read: %v", err)
	}
	defer stream.Close()

	body, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if string(body) != "o go" {
		t.Fatalf("got body=%q, want o go", string(body))
	}
}
