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

func (c *Client) propfind(path string, self bool, body string, resp interface{}, parse func(resp interface{})) error {
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

	parseXML(rs.Body, resp, parse)

	return nil
}
