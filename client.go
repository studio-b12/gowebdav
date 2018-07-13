package gowebdav

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"os"
	pathpkg "path"
	"strings"
	"time"
)

// Client defines our structure
type Client struct {
	root    string
	headers http.Header
	c       *http.Client
	auth    Authenticator
}

// Authenticator stub
type Authenticator interface {
	Type() string
	User() string
	Pass() string
	Authorize(*Client, string, string)
}

// NoAuth structure holds our credentials
type NoAuth struct {
	user string
	pw   string
}

// Type identifies the authenticator
func (n *NoAuth) Type() string {
	return "NoAuth"
}

// User returns the current user
func (n *NoAuth) User() string {
	return n.user
}

// Pass returns the current password
func (n *NoAuth) Pass() string {
	return n.pw
}

// Authorize the current request
func (n *NoAuth) Authorize(c *Client, method string, path string) {
}

// NewClient creates a new instance of client
func NewClient(uri, user, pw string) *Client {
	return &Client{FixSlash(uri), make(http.Header), &http.Client{}, &NoAuth{user, pw}}
}

// SetHeader lets us set arbitrary headers for a given client
func (c *Client) SetHeader(key, value string) {
	c.headers.Add(key, value)
}

// SetTimeout exposes the ability to set a time limit for requests
func (c *Client) SetTimeout(timeout time.Duration) {
	c.c.Timeout = timeout
}

// SetTransport exposes the ability to define custom transports
func (c *Client) SetTransport(transport http.RoundTripper) {
	c.c.Transport = transport
}

// Connect connects to our dav server
func (c *Client) Connect() error {
	rs, err := c.options("/")
	if err != nil {
		return err
	}

	err = rs.Body.Close()
	if err != nil {
		return err
	}

	if rs.StatusCode != 200 || (rs.Header.Get("Dav") == "" && rs.Header.Get("DAV") == "") {
		return newPathError("Connect", c.root, rs.StatusCode)
	}

	return nil
}

type props struct {
	Status      string   `xml:"DAV: status"`
	Name        string   `xml:"DAV: prop>displayname,omitempty"`
	Type        xml.Name `xml:"DAV: prop>resourcetype>collection,omitempty"`
	Size        string   `xml:"DAV: prop>getcontentlength,omitempty"`
	ContentType string   `xml:"DAV: prop>getcontenttype,omitempty"`
	ETag        string   `xml:"DAV: prop>getetag,omitempty"`
	Modified    string   `xml:"DAV: prop>getlastmodified,omitempty"`
}

type response struct {
	Href  string  `xml:"DAV: href"`
	Props []props `xml:"DAV: propstat"`
}

func getProps(r *response, status string) *props {
	for _, prop := range r.Props {
		if strings.Contains(prop.Status, status) {
			return &prop
		}
	}
	return nil
}

// ReadDir reads the contents of a remote directory
func (c *Client) ReadDir(path string) ([]os.FileInfo, error) {
	path = FixSlashes(path)
	files := make([]os.FileInfo, 0)
	skipSelf := true
	parse := func(resp interface{}) error {
		r := resp.(*response)

		if skipSelf {
			skipSelf = false
			if p := getProps(r, "200"); p != nil && p.Type.Local == "collection" {
				r.Props = nil
				return nil
			}
			return newPathError("ReadDir", path, 405)
		}

		if p := getProps(r, "200"); p != nil {
			f := new(File)
			if ps, err := url.QueryUnescape(r.Href); err == nil {
				f.name = pathpkg.Base(ps)
			} else {
				f.name = p.Name
			}
			f.path = path + f.name
			f.modified = parseModified(&p.Modified)
			f.etag = p.ETag
			f.contentType = p.ContentType

			if p.Type.Local == "collection" {
				f.path += "/"
				f.size = 0
				f.isdir = true
			} else {
				f.size = parseInt64(&p.Size)
				f.isdir = false
			}

			files = append(files, *f)
		}

		r.Props = nil
		return nil
	}

	err := c.propfind(path, false,
		`<d:propfind xmlns:d='DAV:'>
			<d:prop>
				<d:displayname/>
				<d:resourcetype/>
				<d:getcontentlength/>
				<d:getcontenttype/>
				<d:getetag/>
				<d:getlastmodified/>
			</d:prop>
		</d:propfind>`,
		&response{},
		parse)

	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			err = newPathErrorErr("ReadDir", path, err)
		}
	}
	return files, err
}

// Stat returns the file stats for a specified path
func (c *Client) Stat(path string) (os.FileInfo, error) {
	var f *File
	parse := func(resp interface{}) error {
		r := resp.(*response)
		if p := getProps(r, "200"); p != nil && f == nil {
			f = new(File)
			f.name = p.Name
			f.path = path
			f.etag = p.ETag
			f.contentType = p.ContentType

			if p.Type.Local == "collection" {
				if !strings.HasSuffix(f.path, "/") {
					f.path += "/"
				}
				f.size = 0
				f.modified = time.Unix(0, 0)
				f.isdir = true
			} else {
				f.size = parseInt64(&p.Size)
				f.modified = parseModified(&p.Modified)
				f.isdir = false
			}
		}

		r.Props = nil
		return nil
	}

	err := c.propfind(path, true,
		`<d:propfind xmlns:d='DAV:'>
			<d:prop>
				<d:displayname/>
				<d:resourcetype/>
				<d:getcontentlength/>
				<d:getcontenttype/>
				<d:getetag/>
				<d:getlastmodified/>
			</d:prop>
		</d:propfind>`,
		&response{},
		parse)

	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			err = newPathErrorErr("ReadDir", path, err)
		}
	}
	return f, err
}

// Remove removes a remote file
func (c *Client) Remove(path string) error {
	return c.RemoveAll(path)
}

// RemoveAll removes remote files
func (c *Client) RemoveAll(path string) error {
	rs, err := c.req("DELETE", path, nil, nil)
	if err != nil {
		return newPathError("Remove", path, 400)
	}
	err = rs.Body.Close()
	if err != nil {
		return err
	}

	if rs.StatusCode == 200 || rs.StatusCode == 204 || rs.StatusCode == 404 {
		return nil
	}

	return newPathError("Remove", path, rs.StatusCode)
}

// Mkdir makes a directory
func (c *Client) Mkdir(path string, _ os.FileMode) error {
	path = FixSlashes(path)
	status := c.mkcol(path)
	if status == 201 {
		return nil
	}

	return newPathError("Mkdir", path, status)
}

// MkdirAll like mkdir -p, but for webdav
func (c *Client) MkdirAll(path string, _ os.FileMode) error {
	path = FixSlashes(path)
	status := c.mkcol(path)
	if status == 201 {
		return nil
	} else if status == 409 {
		paths := strings.Split(path, "/")
		sub := "/"
		for _, e := range paths {
			if e == "" {
				continue
			}
			sub += e + "/"
			status = c.mkcol(sub)
			if status != 201 {
				return newPathError("MkdirAll", sub, status)
			}
		}
		return nil
	}

	return newPathError("MkdirAll", path, status)
}

// Rename moves a file from A to B
func (c *Client) Rename(oldpath, newpath string, overwrite bool) error {
	return c.copymove("MOVE", oldpath, newpath, overwrite)
}

// Copy copies a file from A to B
func (c *Client) Copy(oldpath, newpath string, overwrite bool) error {
	return c.copymove("COPY", oldpath, newpath, overwrite)
}

// Read reads the contents of a remote file
func (c *Client) Read(path string) ([]byte, error) {
	var stream io.ReadCloser
	var err error

	if stream, err = c.ReadStream(path); err != nil {
		return nil, err
	}
	defer stream.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(stream)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ReadStream reads the stream for a given path
func (c *Client) ReadStream(path string) (io.ReadCloser, error) {
	rs, err := c.req("GET", path, nil, nil)
	if err != nil {
		return nil, newPathErrorErr("ReadStream", path, err)
	}

	if rs.StatusCode == 200 {
		return rs.Body, nil
	}

	rs.Body.Close()
	return nil, newPathError("ReadStream", path, rs.StatusCode)
}

// Write writes data to a given path
func (c *Client) Write(path string, data []byte, _ os.FileMode) error {
	s := c.put(path, bytes.NewReader(data))
	switch s {

	case 200, 201, 204:
		return nil

	case 409:
		err := c.createParentCollection(path)
		if err != nil {
			return err
		}

		s = c.put(path, bytes.NewReader(data))
		if s == 200 || s == 201 || s == 204 {
			return nil
		}
	}

	return newPathError("Write", path, s)
}

// WriteStream writes a stream
func (c *Client) WriteStream(path string, stream io.Reader, _ os.FileMode) error {

	err := c.createParentCollection(path)
	if err != nil {
		return err
	}

	s := c.put(path, stream)

	switch s {
	case 200, 201, 204:
		return nil

	default:
		return newPathError("WriteStream", path, s)
	}
}
