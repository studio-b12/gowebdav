package gowebdav

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// PathEscape escapes all segments of a given path
func PathEscape(path string) string {
	s := strings.Split(path, "/")
	for i, e := range s {
		s[i] = url.PathEscape(e)
	}
	return strings.Join(s, "/")
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
	if !strings.HasPrefix(s, "/") {
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
	// TODO - make String return an error as well
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

// limitedReadCloser wraps a io.ReadCloser and limits the number of bytes that can be read from it.
type limitedReadCloser struct {
	rc        io.ReadCloser
	remaining int
}

func (l *limitedReadCloser) Read(buf []byte) (int, error) {
	if l.remaining <= 0 {
		return 0, io.EOF
	}

	if len(buf) > l.remaining {
		buf = buf[0:l.remaining]
	}

	n, err := l.rc.Read(buf)
	l.remaining -= n

	return n, err
}

func (l *limitedReadCloser) Close() error {
	return l.rc.Close()
}

func getContentLength(reader io.Reader) (int64, error) {
	contentLength := int64(-1)
	switch reader := reader.(type) {
	case *bytes.Buffer:
		contentLength = int64(reader.Len())
	case io.Seeker:
		currentPos, err := reader.Seek(0, io.SeekCurrent)
		if err != nil {
			return -1, err
		}

		endPos, err := reader.Seek(0, io.SeekEnd)
		if err != nil {
			return -1, err
		}

		contentLength = endPos - currentPos

		_, err = reader.Seek(currentPos, io.SeekStart)
		if err != nil {
			return -1, err
		}
	}

	return contentLength, nil
}
