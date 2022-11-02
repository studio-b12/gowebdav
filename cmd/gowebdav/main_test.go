package main

import (
	"context"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/net/webdav"
)

func basicAuth(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user, passwd, ok := r.BasicAuth(); ok {
			if user == "user" && passwd == "password" {
				h.ServeHTTP(w, r)
				return
			}

			http.Error(w, "not authorized", 403)
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="x"`)
			w.WriteHeader(401)
		}
	}
}

func newServer(t *testing.T) (*httptest.Server, webdav.FileSystem, context.Context) {
	mux := http.NewServeMux()
	fs := webdav.NewMemFS()
	ctx := fillFs(t, fs)
	mux.HandleFunc("/", basicAuth(&webdav.Handler{
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(),
	}))
	srv := httptest.NewServer(mux)
	os.Setenv("ROOT", srv.URL)
	os.Setenv("USER", "user")
	os.Setenv("PASSWORD", "password")
	return srv, fs, ctx
}

func fillFs(t *testing.T, fs webdav.FileSystem) context.Context {
	ctx := context.Background()
	f, err := fs.OpenFile(ctx, "hello.txt", os.O_CREATE, 0644)
	if err != nil {
		t.Errorf("fail to crate file: %v", err)
	}
	f.Write([]byte("hello gowebdav\n"))
	f.Close()
	err = fs.Mkdir(ctx, "/test", 0755)
	if err != nil {
		t.Errorf("fail to crate directory: %v", err)
	}
	f, err = fs.OpenFile(ctx, "/test/test.txt", os.O_CREATE, 0644)
	if err != nil {
		t.Errorf("fail to crate file: %v", err)
	}
	f.Write([]byte("test test gowebdav\n"))
	f.Close()
	return ctx
}

func TestLs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	srv, _, _ := newServer(t)

	defer srv.Close()

	flag.CommandLine = flag.NewFlagSet("ls", flag.ExitOnError)

	os.Args = []string{"ls", "-X", "ls", "/"}
	main()
}
