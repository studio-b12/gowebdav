package gowebdav

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DigestAuth structure holds our credentials
type DigestAuth struct {
	user        string
	pw          string
	digestParts map[string]string
}

// NewDigestAuth creates a new instance of our Digest Authenticator
func NewDigestAuth(login, secret string, rs *http.Response) (Authenticator, error) {
	return &DigestAuth{user: login, pw: secret, digestParts: digestParts(rs)}, nil
}

// Authorize the current request
func (d *DigestAuth) Authorize(c *http.Client, rq *http.Request, path string) error {
	d.digestParts["uri"] = path
	d.digestParts["method"] = rq.Method
	d.digestParts["username"] = d.user
	d.digestParts["password"] = d.pw
	rq.Header.Set("Authorization", getDigestAuthorization(d.digestParts))
	return nil
}

// Verify checks for authentication issues and may trigger a re-authentication
func (d *DigestAuth) Verify(c *http.Client, rs *http.Response, path string) (redo bool, err error) {
	if rs.StatusCode == 401 {
		if isStaled(rs) {
			redo = true
			err = ErrAuthChanged
		} else {
			err = NewPathError("Authorize", path, rs.StatusCode)
		}
	}
	return
}

// Close cleans up all resources
func (d *DigestAuth) Close() error {
	return nil
}

// Clone creates a copy of itself
func (d *DigestAuth) Clone() Authenticator {
	parts := make(map[string]string, len(d.digestParts))
	for k, v := range d.digestParts {
		parts[k] = v
	}
	return &DigestAuth{user: d.user, pw: d.pw, digestParts: parts}
}

// String toString
func (d *DigestAuth) String() string {
	return fmt.Sprintf("DigestAuth login: %s", d.user)
}

func digestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	header := resp.Header.Get("Www-Authenticate")
	if header == "" {
		return result
	}

	if _, params, ok := strings.Cut(header, " "); ok {
		header = params
	}

	for _, directive := range splitHeaderDirectives(header) {
		key, value, ok := strings.Cut(directive, "=")
		if !ok {
			continue
		}

		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.Trim(strings.TrimSpace(value), `"`)

		switch key {
		case "nonce", "realm", "qop", "opaque", "algorithm", "entitybody":
			result[key] = value
		}
	}

	if qop, ok := result["qop"]; ok {
		result["qop"] = selectDigestQOP(qop)
	}
	if algorithm, ok := result["algorithm"]; ok {
		result["algorithm"] = normalizeDigestAlgorithm(algorithm)
	}

	return result
}

func splitHeaderDirectives(header string) []string {
	var (
		result  []string
		current strings.Builder
		quoted  bool
	)

	for _, r := range header {
		switch r {
		case '"':
			quoted = !quoted
			current.WriteRune(r)
		case ',':
			if quoted {
				current.WriteRune(r)
				continue
			}
			part := strings.TrimSpace(current.String())
			if part != "" {
				result = append(result, part)
			}
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}

	part := strings.TrimSpace(current.String())
	if part != "" {
		result = append(result, part)
	}

	return result
}

func selectDigestQOP(qop string) string {
	var firstSupported string
	for _, token := range strings.Split(qop, ",") {
		token = strings.ToLower(strings.TrimSpace(token))
		switch token {
		case "auth":
			return "auth"
		case "auth-int":
			if firstSupported == "" {
				firstSupported = token
			}
		}
	}
	return firstSupported
}

func normalizeDigestAlgorithm(algorithm string) string {
	switch strings.ToUpper(strings.TrimSpace(algorithm)) {
	case "MD5":
		return "MD5"
	case "MD5-SESS":
		return "MD5-sess"
	default:
		return strings.TrimSpace(algorithm)
	}
}

func getMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getCnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}

func getDigestAuthorization(digestParts map[string]string) string {
	d := digestParts
	// These are the correct ha1 and ha2 for qop=auth. We should probably check for other types of qop.

	var (
		ha1        string
		ha2        string
		nonceCount = "00000001"
		cnonce     = getCnonce()
		response   string
	)

	// 'ha1' value depends on value of "algorithm" field
	switch d["algorithm"] {
	case "MD5", "":
		ha1 = getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	case "MD5-sess":
		ha1 = getMD5(
			fmt.Sprintf("%s:%s:%s",
				getMD5(d["username"]+":"+d["realm"]+":"+d["password"]),
				d["nonce"],
				cnonce,
			),
		)
	}

	// 'ha2' value depends on value of "qop" field
	switch d["qop"] {
	case "auth", "":
		ha2 = getMD5(d["method"] + ":" + d["uri"])
	case "auth-int":
		if d["entityBody"] != "" {
			ha2 = getMD5(d["method"] + ":" + d["uri"] + ":" + getMD5(d["entityBody"]))
		}
	}

	// 'response' value depends on value of "qop" field
	switch d["qop"] {
	case "":
		response = getMD5(
			fmt.Sprintf("%s:%s:%s",
				ha1,
				d["nonce"],
				ha2,
			),
		)
	case "auth", "auth-int":
		response = getMD5(
			fmt.Sprintf("%s:%s:%s:%s:%s:%s",
				ha1,
				d["nonce"],
				nonceCount,
				cnonce,
				d["qop"],
				ha2,
			),
		)
	}

	authorization := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", nc=%s, cnonce="%s", response="%s"`,
		d["username"], d["realm"], d["nonce"], d["uri"], nonceCount, cnonce, response)

	if d["qop"] != "" {
		authorization += fmt.Sprintf(`, qop=%s`, d["qop"])
	}

	if d["opaque"] != "" {
		authorization += fmt.Sprintf(`, opaque="%s"`, d["opaque"])
	}

	if d["algorithm"] != "" {
		authorization += fmt.Sprintf(`, algorithm=%s`, d["algorithm"])
	}

	return authorization
}

func isStaled(rs *http.Response) bool {
	header := rs.Header.Get("Www-Authenticate")
	if len(header) > 0 {
		directives := strings.Split(header, ",")
		for i := range directives {
			name, value, _ := strings.Cut(strings.Trim(directives[i], " "), "=")
			if strings.EqualFold(name, "stale") {
				return strings.EqualFold(value, "true")
			}
		}
	}
	return false
}
