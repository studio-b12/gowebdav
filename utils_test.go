package gowebdav

import "testing"

func TestJoin(t *testing.T) {
	eq(t, "/", "", "")
	eq(t, "/", "/", "/")
	eq(t, "/foo", "", "/foo")
	eq(t, "foo/foo", "foo/", "/foo")
	eq(t, "foo/foo", "foo/", "foo")
}

func eq(t *testing.T, expected string, s0 string, s1 string) {
	s := Join(s0, s1)
	if s != expected {
		t.Error("For", "'"+s0+"','"+s1+"'", "expeted", "'"+expected+"'", "got", "'"+s+"'")
	}
}
