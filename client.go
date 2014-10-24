package gowebdav

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"
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
	if rs, err := c.options("/"); err == nil {
		defer rs.Body.Close()

		if rs.StatusCode != 200 || (rs.Header.Get("Dav") == "" && rs.Header.Get("DAV") == "") {
			return errors.New(fmt.Sprintf("Bad Request: %d - %s", rs.StatusCode, c.root))
		}

		// TODO check PROPFIND if path is collection

		return nil
	} else {
		return err
	}
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
	parse := func(resp interface{}) {
		r := resp.(*response)

		if skipSelf {
			skipSelf = false
			r.Props = nil
			return
		}

		if p := getProps(r, "200"); p != nil {
			f := new(File)
			f.name = p.Name
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
		err = &os.PathError{"ReadDir", path, err}
	}
	return files, err
}

func (c *Client) Remove(path string) error {
	return c.RemoveAll(path)
}

func (c *Client) RemoveAll(path string) error {
	rs, err := c.reqDo("DELETE", path, nil)
	if err != nil {
		return newPathError("Remove", path, 400)
	}
	defer rs.Body.Close()

	if rs.StatusCode == 200 {
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

func (c *Client) Read(path string) {
	fmt.Println("Read " + path)
}
