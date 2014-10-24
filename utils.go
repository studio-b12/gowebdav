package gowebdav

import (
	"bytes"
	"encoding/xml"
	"errors"
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
	return &os.PathError{op, path, errors.New(fmt.Sprintf("%d", statusCode))}
}

func newPathErrorErr(op string, path string, err error) error {
	return &os.PathError{op, path, err}
}

func FixSlash(s string) string {
	if !strings.HasSuffix(s, "/") {
		s += "/"
	}
	return s
}

func FixSlashes(s string) string {
	if s[0] != '/' {
		s = "/" + s
	}
	return FixSlash(s)
}

func Join(path0 string, path1 string) string {
	return strings.TrimSuffix(path0, "/") + "/" + strings.TrimPrefix(path1, "/")
}

func String(r io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
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
