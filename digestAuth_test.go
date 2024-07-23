package gowebdav

import (
	"errors"
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
	ex := `Digest username="user", realm="", nonce="", uri="/", nc=1, cnonce="`
	if strings.Index(rq.Header.Get("Authorization"), ex) != 0 {
		t.Error("got wrong Authorization header: " + rq.Header.Get("Authorization"))
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
