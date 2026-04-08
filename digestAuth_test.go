package gowebdav

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestNewDigestAuth(t *testing.T) {
	a := &DigestAuth{user: "user", pw: "password", digestParts: make(map[string]string, 0)}

	ex := "DigestAuth login: user"
	if a.String() != ex {
		t.Error("expected: " + ex + " got: " + a.String())
	}

	if a.Clone() == a {
		t.Error("expected a different instance")
	}

	if a.Close() != nil {
		t.Error("expected close without errors")
	}
}

func TestDigestAuthAuthorize(t *testing.T) {
	a := &DigestAuth{user: "user", pw: "password", digestParts: make(map[string]string, 0)}
	rq, _ := http.NewRequest("GET", "http://localhost/", nil)
	a.Authorize(nil, rq, "/")
	// TODO this is a very lazy test it cuts of cnonce
	ex := `Digest username="user", realm="", nonce="", uri="/", nc=00000001, cnonce="`
	if strings.Index(rq.Header.Get("Authorization"), ex) != 0 {
		t.Error("got wrong Authorization header: " + rq.Header.Get("Authorization"))
	}
}

func TestDigestAuthAuthorizeMD5Sess(t *testing.T) {
	a := &DigestAuth{
		user: "user",
		pw:   "password",
		digestParts: map[string]string{
			"realm":     "testrealm@host.com",
			"nonce":     "dcd98b7102dd2f0e8b11d0f600bfb0c093",
			"algorithm": "MD5-sess",
			"qop":       "auth",
		},
	}
	rq, _ := http.NewRequest("GET", "http://localhost/", nil)
	if err := a.Authorize(nil, rq, "/dir/index.html"); err != nil {
		t.Fatalf("authorize: %v", err)
	}

	parts := parseDigestAuthorizationHeader(t, rq.Header.Get("Authorization"))
	expectedHA1 := getMD5(fmt.Sprintf("%s:%s:%s",
		getMD5("user:testrealm@host.com:password"),
		"dcd98b7102dd2f0e8b11d0f600bfb0c093",
		parts["cnonce"],
	))
	expectedHA2 := getMD5("GET:/dir/index.html")
	expectedResponse := getMD5(fmt.Sprintf("%s:%s:%s:%s:%s:%s",
		expectedHA1,
		"dcd98b7102dd2f0e8b11d0f600bfb0c093",
		parts["nc"],
		parts["cnonce"],
		"auth",
		expectedHA2,
	))

	if parts["response"] != expectedResponse {
		t.Fatalf("got response=%q, want %q", parts["response"], expectedResponse)
	}
	if parts["algorithm"] != "MD5-sess" {
		t.Fatalf("got algorithm=%q, want MD5-sess", parts["algorithm"])
	}
}

func TestDigestAuthVerify(t *testing.T) {
	a := &DigestAuth{user: "user", pw: "password", digestParts: make(map[string]string, 0)}

	// Nominal test: 200 OK response
	rs := &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
	}

	redo, err := a.Verify(nil, rs, "/")

	if err != nil {
		t.Errorf("got error: %v, want nil", err)
	}

	if redo {
		t.Errorf("got redo: %t, want false", redo)
	}

	// Digest expiration test: 401 Unauthorized response with stale directive in WWW-Authenticate header
	rs = &http.Response{
		Status:     "401 Unauthorized",
		StatusCode: http.StatusUnauthorized,
		Header: http.Header{
			"Www-Authenticate": []string{"Digest realm=\"webdav\", nonce=\"YVvALpkdBgA=931bbf2b6fa9dda227361dba38a735f005fd9f97\", algorithm=MD5, qop=\"auth\", stale=true"},
		},
	}

	redo, err = a.Verify(nil, rs, "/")

	if !errors.Is(err, ErrAuthChanged) {
		t.Errorf("got error: %v, want ErrAuthChanged", err)
	}

	if !redo {
		t.Errorf("got redo: %t, want true", redo)
	}
}

func parseDigestAuthorizationHeader(t *testing.T, header string) map[string]string {
	t.Helper()

	header = strings.TrimPrefix(header, "Digest ")
	result := make(map[string]string)
	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			t.Fatalf("malformed digest header part %q", part)
		}
		result[kv[0]] = strings.Trim(kv[1], `"`)
	}
	return result
}
