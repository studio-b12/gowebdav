package gowebdav

import (
	"net/url"
	"os"
	"strings"
)

func parseLine(s string) (login, pass string) {
	fields := strings.Fields(s)
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "login":
			if i+1 >= len(fields) {
				return login, pass
			}
			login = fields[i+1]
			i++
		case "password":
			if i+1 >= len(fields) {
				return login, pass
			}
			pass = fields[i+1]
			i++
		}
	}
	return login, pass
}

func hostMatches(u *url.URL, machine string) bool {
	return strings.EqualFold(machine, u.Host) || strings.EqualFold(machine, u.Hostname())
}

// ReadConfig reads login and password configuration from ~/.netrc
// machine foo.com login username password 123456
func ReadConfig(uri, netrc string) (string, string) {
	u, err := url.Parse(uri)
	if err != nil || u.Host == "" {
		return "", ""
	}

	data, err := os.ReadFile(netrc)
	if err != nil {
		return "", ""
	}

	tokens := strings.Fields(string(data))
	for i := 0; i < len(tokens); {
		if tokens[i] != "machine" {
			i++
			continue
		}
		if i+1 >= len(tokens) {
			break
		}

		machine := tokens[i+1]
		i += 2

		login, pass := "", ""
		for i < len(tokens) && tokens[i] != "machine" {
			switch tokens[i] {
			case "login":
				if i+1 < len(tokens) {
					login = tokens[i+1]
					i += 2
					continue
				}
			case "password":
				if i+1 < len(tokens) {
					pass = tokens[i+1]
					i += 2
					continue
				}
			}
			i++
		}

		if hostMatches(u, machine) && login != "" && pass != "" {
			return login, pass
		}
	}

	return "", ""
}
