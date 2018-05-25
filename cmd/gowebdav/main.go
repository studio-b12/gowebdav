package main

import (
	"errors"
	"flag"
	"fmt"
	d "github.com/studio-b12/gowebdav"
	"io"
	"os"
	"strings"
)

func main() {
	root := flag.String("root", os.Getenv("ROOT"), "WebDAV Endpoint [ENV.ROOT]")
	usr := flag.String("user", os.Getenv("USER"), "User [ENV.USER]")
	pw := flag.String("pw", os.Getenv("PASSWORD"), "Password [ENV.PASSWORD]")
	m := flag.String("X", "GET", "Method")
	flag.Parse()

	if *root == "" {
		fail("Set WebDAV ROOT")
	}

	var path0, path1 string
	switch len(flag.Args()) {
	case 1:
		path0 = flag.Args()[0]
	case 2:
		path1 = flag.Args()[1]
	default:
		fail("Unsupported arguments")
	}

	c := d.NewClient(*root, *usr, *pw)
	if err := c.Connect(); err != nil {
		fail(fmt.Sprintf("Failed to connect due to: %s", err.Error()))
	}

	cmd := getCmd(strings.ToUpper(*m))

	if e := cmd(c, path0, path1); e != nil {
		fail(e)
	}
}

func fail(err interface{}) {
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(-1)
}

func getCmd(method string) func(c *d.Client, p0, p1 string) error {
	switch method {
	case "LS", "LIST", "PROPFIND":
		return cmdLs

	case "STAT":
		return cmdStat

	case "GET", "PULL", "READ":
		return cmdGet

	case "DELETE", "RM", "DEL":
		return cmdRm

	case "MKCOL", "MKDIR":
		return cmdMkdir

	case "MKCOLALL", "MKDIRALL":
		return cmdMkdirAll

	case "RENAME", "MV", "MOVE":
		return cmdMv

	case "COPY", "CP":
		return cmdCp

	case "PUT", "PUSH", "WRITE":
		return cmdPut

	default:
		return func(c *d.Client, p0, p1 string) (err error) {
			return errors.New("Unsupported method: " + method)
		}
	}
}

func cmdLs(c *d.Client, p0, p1 string) (err error) {
	files, err := c.ReadDir(p0)
	if err == nil {
		fmt.Println(fmt.Sprintf("ReadDir: '%s' entries: %d ", p0, len(files)))
		for _, f := range files {
			fmt.Println(f)
		}
	}
	return
}

func cmdStat(c *d.Client, p0, p1 string) (err error) {
	file, err := c.Stat(p0)
	if err == nil {
		fmt.Println(file)
	}
	return
}

func cmdGet(c *d.Client, p0, p1 string) (err error) {
	bytes, err := c.Read(p0)
	if err == nil {
		if err = writeFile(p1, bytes, 0644); err == nil {
			fmt.Println(fmt.Sprintf("Written %d bytes to: %s", len(bytes), p1))
		}
	}
	return
}

func cmdRm(c *d.Client, p0, p1 string) (err error) {
	if err = c.Remove(p0); err == nil {
		fmt.Println("RM: " + p0)
	}
	return
}

func cmdMkdir(c *d.Client, p0, p1 string) (err error) {
	if err = c.Mkdir(p0, 0755); err == nil {
		fmt.Println("MkDir: " + p0)
	}
	return
}

func cmdMkdirAll(c *d.Client, p0, p1 string) (err error) {
	if err = c.MkdirAll(p0, 0755); err == nil {
		fmt.Println("MkDirAll: " + p0)
	}
	return
}

func cmdMv(c *d.Client, p0, p1 string) (err error) {
	if err = c.Rename(p0, p1, true); err == nil {
		fmt.Println("Rename: " + p0 + " -> " + p1)
	}
	return
}

func cmdCp(c *d.Client, p0, p1 string) (err error) {
	if err = c.Copy(p0, p1, true); err == nil {
		fmt.Println("Copy: " + p0 + " -> " + p1)
	}
	return
}

func cmdPut(c *d.Client, p0, p1 string) (err error) {
	stream, err := getStream(p1)
	if err == nil {
		if err = c.WriteStream(p0, stream, 0644); err == nil {
			fmt.Println(fmt.Sprintf("Put: '%s' -> %s", p1, p0))
		}
	}
	return
}

func writeFile(path string, bytes []byte, mode os.FileMode) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	return err
}

func getStream(pathOrString string) (io.ReadCloser, error) {
	fi, err := os.Stat(pathOrString)
	if err == nil {
		if fi.IsDir() {
			return nil, &os.PathError{
				Op:   "Open",
				Path: pathOrString,
				Err:  errors.New("Path: '" + pathOrString + "' is a directory"),
			}
		}
		f, err := os.Open(pathOrString)
		if err == nil {
			return f, nil
		}
		return nil, &os.PathError{
			Op:   "Open",
			Path: pathOrString,
			Err:  err,
		}
	}
	return nopCloser{strings.NewReader(pathOrString)}, nil
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}
