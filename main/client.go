package main

import (
	"flag"
	"fmt"
	d "gowebdav"
	"os"
	"strings"
)

func Fail(err interface{}) {
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
	}
	os.Exit(-1)
}

func main() {
	root := flag.String("root", "URL", "WebDAV Endpoint")
	usr := flag.String("user", "", "user")
	pw := flag.String("pw", "", "password")
	mm := strings.ToUpper(*(flag.String("X", "GET", "Method ...")))
	m := &mm
	flag.Parse()

	if *root == "URL" {
		Fail(nil)
	}

	c := d.NewClient(*root, *usr, *pw)
	if err := c.Connect(); err != nil {
		Fail(fmt.Sprintf("Failed to connect due to: %s", err.Error()))
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

		case "GET":
			c.Read(path)

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
			Fail(nil)
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

		default:
			Fail(nil)
		}
	} else {
		Fail(nil)
	}
}
