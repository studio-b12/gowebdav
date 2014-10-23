package gowebdav

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

	if !strings.HasSuffix(c.root, "/") {
		c.root += "/"
	}

	return c
}

func (c *Client) Connect() error {
	if rs, err := c.Options("/"); err == nil {
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

func (c *Client) List(path string) (*[]*File, error) {
	files := make([]*File, 0)
	parse := func(resp interface{}) {
		r := resp.(*response)
		if p := getProps(r, "200"); p != nil {
			var f File
			if p.Type.Local == "collection" {
				f = directory{p.Name}
			} else {
				f = file{p.Name, parseUint(&p.Size), parseModified(&p.Modified)}
			}

			files = append(files, &f)
			r.Props = nil
		}
	}

	err := c.Propfind(path, false,
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
	return &files, err
}

func (c *Client) Read(path string) {
	fmt.Println("Read " + path)
}
