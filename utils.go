package gowebdav

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func log(msg interface{}) {
	fmt.Println(msg)
}

func newPathError(op string, path string, statusCode int) error {
	return &os.PathError{
		Op:   op,
		Path: path,
		Err:  fmt.Errorf("%d", statusCode),
	}
}

func newPathErrorErr(op string, path string, err error) error {
	return &os.PathError{
		Op:   op,
		Path: path,
		Err:  err,
	}
}

// FixSlash appends a trailing / to our string
func FixSlash(s string) string {
	if !strings.HasSuffix(s, "/") {
		s += "/"
	}
	return s
}

// FixSlashes appends and prepends a / if they are missing
func FixSlashes(s string) string {
	if s[0] != '/' {
		s = "/" + s
	}
	return FixSlash(s)
}

// Join joins two paths
func Join(path0 string, path1 string) string {
	return strings.TrimSuffix(path0, "/") + "/" + strings.TrimPrefix(path1, "/")
}

// String pulls a string out of our io.Reader
func String(r io.Reader) string {
	buf := new(bytes.Buffer)
	// TODO - mkae String return an error as well
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func parseUint(s *string) uint {
	if n, e := strconv.ParseUint(*s, 10, 32); e == nil {
		return uint(n)
	}
	return 0
}

func parseInt64(s *string) int64 {
	if n, e := strconv.ParseInt(*s, 10, 64); e == nil {
		return n
	}
	return 0
}

func parseModified(s *string) time.Time {
	if t, e := time.Parse(time.RFC1123, *s); e == nil {
		return t
	}
	return time.Unix(0, 0)
}

func parseXML(data io.Reader, resp interface{}, parse func(resp interface{}) error) error {
	decoder := xml.NewDecoder(data)
	for t, _ := decoder.Token(); t != nil; t, _ = decoder.Token() {
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "response" {
				if e := decoder.DecodeElement(resp, &se); e == nil {
					if err := parse(resp); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
