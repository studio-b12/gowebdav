package gowebdav

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) req(method, path string, body io.Reader, intercept func(*http.Request)) (req *http.Response, err error) {
	r, err := http.NewRequest(method, Join(c.root, path), body)
	if err != nil {
		return nil, err
	}
	for k, vals := range c.headers {
		for _, v := range vals {
			r.Header.Add(k, v)
		}
	}

	if intercept != nil {
		intercept(r)
	}

	return c.c.Do(r)
}

func (c *Client) mkcol(path string) int {
	rs, err := c.req("MKCOL", path, nil, nil)
	defer rs.Body.Close()
	if err != nil {
		return 400
	}

	if rs.StatusCode == 201 || rs.StatusCode == 405 {
		return 201
	}

	return rs.StatusCode
}

func (c *Client) options(path string) (*http.Response, error) {
	return c.req("OPTIONS", path, nil, func(rq *http.Request) {
		rq.Header.Add("Depth", "0")
	})
}

func (c *Client) propfind(path string, self bool, body string, resp interface{}, parse func(resp interface{}) error) error {
	rs, err := c.req("PROPFIND", path, strings.NewReader(body), func(rq *http.Request) {
		if self {
			rq.Header.Add("Depth", "0")
		} else {
			rq.Header.Add("Depth", "1")
		}
		rq.Header.Add("Content-Type", "text/xml;charset=UTF-8")
		rq.Header.Add("Accept", "application/xml,text/xml")
		rq.Header.Add("Accept-Charset", "utf-8")
		// TODO add support for 'gzip,deflate;q=0.8,q=0.7'
		rq.Header.Add("Accept-Encoding", "")
	})
	defer rs.Body.Close()
	if err != nil {
		return err
	}

	if rs.StatusCode != 207 {
		return fmt.Errorf("%s - %s %s", rs.Status, "PROPFIND", path)
	}

	return parseXML(rs.Body, resp, parse)
}

func (c *Client) doCopyMove(method string, oldpath string, newpath string, overwrite bool) (int, io.ReadCloser) {
	rs, err := c.req(method, oldpath, nil, func(rq *http.Request) {
		rq.Header.Add("Destination", Join(c.root, newpath))
		if overwrite {
			rq.Header.Add("Overwrite", "T")
		} else {
			rq.Header.Add("Overwrite", "F")
		}
	})
	if err != nil {
		return 400, nil
	}
	return rs.StatusCode, rs.Body
}

func (c *Client) copymove(method string, oldpath string, newpath string, overwrite bool) error {
	s, data := c.doCopyMove(method, oldpath, newpath, overwrite)
	defer data.Close()

	switch s {
	case 201, 204:
		return nil

	case 207:
		// TODO handle multistat errors, worst case ...
		log(fmt.Sprintf(" TODO handle %s - %s multistatus result %s", method, oldpath, String(data)))

	case 409:
		// TODO create dst path
	}

	return newPathError(method, oldpath, s)
}

func (c *Client) put(path string, stream io.Reader) int {
	rs, err := c.req("PUT", path, stream, nil)
	defer rs.Body.Close()
	if err != nil {
		return 400
	}

	return rs.StatusCode
}
