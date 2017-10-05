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

func fail(err interface{}) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Usage: client FLAGS ARGS")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println("Method <ARGS>")
		fmt.Println(" LS | LIST | PROPFIND <PATH>")
		fmt.Println(" RM | DELETE | DEL <PATH>")
		fmt.Println(" MKDIR | MKCOL <PATH>")
		fmt.Println(" MKDIRALL | MKCOLALL <PATH>")
		fmt.Println(" MV | MOVE | RENAME <OLD_PATH> <NEW_PATH>")
		fmt.Println(" CP | COPY <OLD_PATH> <NEW_PATH>")
		fmt.Println(" GET | PULL | READ <PATH>")
		fmt.Println(" PUT | PUSH | WRITE <PATH> <FILE>")
		fmt.Println(" STAT <PATH>")
	}
	os.Exit(-1)
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

func main() {
	root := flag.String("root", "URL", "WebDAV Endpoint")
	usr := flag.String("user", "", "user")
	pw := flag.String("pw", "", "password")
	m := flag.String("X", "GET", "Method ...")
	flag.Parse()

	if *root == "URL" {
		fail(nil)
	}

	M := strings.ToUpper(*m)
	m = &M

	c := d.NewClient(*root, *usr, *pw)
	if err := c.Connect(); err != nil {
		fail(fmt.Sprintf("Failed to connect due to: %s", err.Error()))
	}
	alen := len(flag.Args())
	if alen == 1 {
		path := flag.Args()[0]
		switch *m {
		case "LS", "LIST", "PROPFIND":
			if files, err := c.ReadDir(path); err == nil {
				fmt.Println(fmt.Sprintf("ReadDir: '%s' entries: %d ", path, len(files)))
				for _, f := range files {
					fmt.Println(f)
				}
			} else {
				fmt.Println(err)
			}

		case "STAT":
			if file, err := c.Stat(path); err == nil {
				fmt.Println(file)
			} else {
				fmt.Println(err)
			}

		case "GET", "PULL", "READ":
			if bytes, err := c.Read(path); err == nil {
				if lidx := strings.LastIndex(path, "/"); lidx != -1 {
					path = path[lidx+1:]
				}
				if err := writeFile(path, bytes, 0644); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(fmt.Sprintf("Written %d bytes to: %s", len(bytes), path))
				}
			} else {
				fmt.Println(err)
			}

		case "DELETE", "RM", "DEL":
			if err := c.Remove(path); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Remove: " + path)
			}

		case "MKCOL", "MKDIR":
			if err := c.Mkdir(path, 0); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("MkDir: " + path)
			}

		case "MKCOLALL", "MKDIRALL":
			if err := c.MkdirAll(path, 0); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("MkDirAll: " + path)
			}

		default:
			fail(nil)
		}

	} else if alen == 2 {
		a0 := flag.Args()[0]
		a1 := flag.Args()[1]
		switch *m {
		case "RENAME", "MV", "MOVE":
			if err := c.Rename(a0, a1, true); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Rename: " + a0 + " -> " + a1)
			}

		case "COPY", "CP":
			if err := c.Copy(a0, a1, true); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Copy: " + a0 + " -> " + a1)
			}

		case "PUT", "PUSH", "WRITE":
			stream, err := getStream(a1)
			if err != nil {
				fail(err)
			}
			if err := c.WriteStream(a0, stream, 0644); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(fmt.Sprintf("Written: '%s' -> %s", a1, a0))
			}

		default:
			fail(nil)
		}
	} else {
		fail(nil)
	}
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
