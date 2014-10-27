package gowebdav

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) req(method string, path string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, Join(c.root, path), body)
	if err != nil {
		return nil, err
	}
	for k, vals := range c.headers {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}
	return req, nil
}

func (c *Client) reqDo(method string, path string, body io.Reader) (*http.Response, error) {
	rq, err := c.req(method, path, body)
	if err != nil {
		return nil, err
	}

	return c.c.Do(rq)
}

func (c *Client) mkcol(path string) int {
	rs, err := c.reqDo("MKCOL", path, nil)
	if err != nil {
		return 400
	}
	defer rs.Body.Close()

	if rs.StatusCode == 201 || rs.StatusCode == 405 {
		return 201
	}

	return rs.StatusCode
}

func (c *Client) options(path string) (*http.Response, error) {
	rq, err := c.req("OPTIONS", path, nil)
	if err != nil {
		return nil, err
	}

	rq.Header.Add("Depth", "0")

	return c.c.Do(rq)
}

func (c *Client) propfind(path string, self bool, body string, resp interface{}, parse func(resp interface{}) error) error {
	rq, err := c.req("PROPFIND", path, strings.NewReader(body))
	if err != nil {
		return err
	}

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

	rs, err := c.c.Do(rq)
	if err != nil {
		return err
	}
	defer rs.Body.Close()

	if rs.StatusCode != 207 {
		return errors.New(fmt.Sprintf("%s - %s %s", rs.Status, rq.Method, rq.URL.String()))
	}

	return parseXML(rs.Body, resp, parse)
}

func (c *Client) copymove(method string, oldpath string, newpath string, overwrite bool) error {
	rq, err := c.req(method, oldpath, nil)
	if err != nil {
		return newPathErrorErr(method, oldpath, err)
	}

	rq.Header.Add("Destination", Join(c.root, newpath))
	if overwrite {
		rq.Header.Add("Overwrite", "T")
	} else {
		rq.Header.Add("Overwrite", "F")
	}

	rs, err := c.c.Do(rq)
	if err != nil {
		return newPathErrorErr(method, oldpath, err)
	}
	defer rs.Body.Close()

	switch rs.StatusCode {
	case 201, 204:
		return nil

	case 207:
		// TODO handle multistat errors, worst case ...
		log(String(rs.Body))

	case 409:
		// TODO create dst path
	}

	return newPathError(method, oldpath, rs.StatusCode)
}

func (c *Client) put(path string, stream io.Reader) int {
	rs, err := c.reqDo("PUT", path, stream)
	if err != nil {
		return 400
	}
	defer rs.Body.Close()
	return rs.StatusCode
}
