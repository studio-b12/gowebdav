package gowebdav

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"os"
	pathpkg "path"
	"strings"
	"time"
)

type Client struct {
	root    string
	headers http.Header
	c       *http.Client
}

func NewClient(uri string, user string, pw string) *Client {
	c := &Client{uri, make(http.Header), &http.Client{}}

	if len(user) > 0 && len(pw) > 0 {
		a := user + ":" + pw
		auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(a))
		c.headers.Add("Authorization", auth)
	}

	c.root = FixSlash(c.root)

	return c
}

func (c *Client) Connect() error {
	rs, err := c.options("/")
	if err == nil {
		rs.Body.Close()

		if rs.StatusCode != 200 || (rs.Header.Get("Dav") == "" && rs.Header.Get("DAV") == "") {
			return newPathError("Connect", c.root, rs.StatusCode)
		}

		_, err = c.ReadDir("/")
	}
	return err
}

type props struct {
	Status   string   `xml:"DAV: status"`
	Name     string   `xml:"DAV: prop>displayname,omitempty"`
	Type     xml.Name `xml:"DAV: prop>resourcetype>collection,omitempty"`
	Size     string   `xml:"DAV: prop>getcontentlength,omitempty"`
	Modified string   `xml:"DAV: prop>getlastmodified,omitempty"`
}
type response struct {
	Href  string  `xml:"DAV: href"`
	Props []props `xml:"DAV: propstat"`
}

func getProps(r *response, status string) *props {
	for _, prop := range r.Props {
		if strings.Index(prop.Status, status) != -1 {
			return &prop
		}
	}
	return nil
}

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

			if p.Type.Local == "collection" {
				f.path += "/"
				f.size = 0
				f.modified = time.Unix(0, 0)
				f.isdir = true
			} else {
				f.size = parseInt64(&p.Size)
				f.modified = parseModified(&p.Modified)
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
				<d:getlastmodified/>
			</d:prop>
		</d:propfind>`,
		&response{},
		parse)

	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			err = &os.PathError{"ReadDir", path, err}
		}
	}
	return files, err
}

func (c *Client) Stat(path string) (os.FileInfo, error) {
	var f *File = nil
	parse := func(resp interface{}) error {
		r := resp.(*response)
		if p := getProps(r, "200"); p != nil && f == nil {
			f = new(File)
			f.name = p.Name
			f.path = path

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
				<d:getlastmodified/>
			</d:prop>
		</d:propfind>`,
		&response{},
		parse)

	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			err = &os.PathError{"ReadDir", path, err}
		}
	}
	return f, err
}

func (c *Client) Remove(path string) error {
	return c.RemoveAll(path)
}

func (c *Client) RemoveAll(path string) error {
	rs, err := c.req("DELETE", path, nil, nil)
	if err != nil {
		return newPathError("Remove", path, 400)
	}
	rs.Body.Close()

	if rs.StatusCode == 200 || rs.StatusCode == 404 {
		return nil
	} else {
		return newPathError("Remove", path, rs.StatusCode)
	}
}

func (c *Client) Mkdir(path string, _ os.FileMode) error {
	path = FixSlashes(path)
	status := c.mkcol(path)
	if status == 201 {
		return nil
	} else {
		return newPathError("Mkdir", path, status)
	}
}

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

func (c *Client) Rename(oldpath string, newpath string, overwrite bool) error {
	return c.copymove("MOVE", oldpath, newpath, overwrite)
}

func (c *Client) Copy(oldpath string, newpath string, overwrite bool) error {
	return c.copymove("COPY", oldpath, newpath, overwrite)
}

func (c *Client) Read(path string) ([]byte, error) {
	if stream, err := c.ReadStream(path); err == nil {
		defer stream.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(stream)
		return buf.Bytes(), nil
	} else {
		return nil, err
	}
}

func (c *Client) ReadStream(path string) (io.ReadCloser, error) {
	rs, err := c.req("GET", path, nil, nil)
	if err != nil {
		return nil, newPathErrorErr("ReadStream", path, err)
	}
	return rs.Body, nil
}

func (c *Client) Write(path string, data []byte, _ os.FileMode) error {
	s := c.put(path, bytes.NewReader(data))
	switch s {

	case 200, 201:
		return nil

	case 409:
		if idx := strings.LastIndex(path, "/"); idx == -1 {
			// faulty root
			return newPathError("Write", path, 500)
		} else {
			if err := c.MkdirAll(path[0:idx+1], 0755); err == nil {
				s = c.put(path, bytes.NewReader(data))
				if s == 200 || s == 201 {
					return nil
				}
			}
		}
	}
	return newPathError("Write", path, s)
}

func (c *Client) WriteStream(path string, stream io.Reader, _ os.FileMode) error {
	// TODO check if parent collection exists
	s := c.put(path, stream)
	switch s {
	case 200, 201:
		return nil

	default:
		return newPathError("WriteStream", path, s)
	}
}
