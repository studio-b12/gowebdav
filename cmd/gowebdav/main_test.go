package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	d "github.com/studio-b12/gowebdav"
)

func TestCmdPutUsesLocalBaseNameForRemoteDirectory(t *testing.T) {
	var putPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("Www-Authenticate", `Basic realm="x"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if user != "user" || pass != "password" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		switch r.Method {
		case "PROPFIND":
			w.WriteHeader(207)
			_, _ = io.WriteString(w, `<?xml version="1.0"?>
<d:multistatus xmlns:d="DAV:">
  <d:response>
    <d:href>/uploads</d:href>
    <d:propstat>
      <d:status>HTTP/1.1 200 OK</d:status>
      <d:prop>
        <d:displayname>uploads</d:displayname>
        <d:resourcetype><d:collection/></d:resourcetype>
        <d:getcontentlength>0</d:getcontentlength>
        <d:getlastmodified>Mon, 02 Jan 2006 15:04:05 GMT</d:getlastmodified>
      </d:prop>
    </d:propstat>
  </d:response>
</d:multistatus>`)
		case "MKCOL":
			w.WriteHeader(http.StatusMethodNotAllowed)
		case "PUT":
			putPath = r.URL.Path
			_, _ = io.Copy(io.Discard, r.Body)
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()

	dir := t.TempDir()
	localPath := filepath.Join(dir, "nested", "local.txt")
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		t.Fatalf("mkdir local dir: %v", err)
	}
	if err := os.WriteFile(localPath, []byte("hello"), 0644); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	cli := d.NewClient(srv.URL, "user", "password")
	if err := cmdPut(cli, "/uploads", localPath); err != nil {
		t.Fatalf("cmdPut: %v", err)
	}
	if putPath != "/uploads/local.txt" {
		t.Fatalf("got put path %q, want /uploads/local.txt", putPath)
	}
}
