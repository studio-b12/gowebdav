package main

import (
	d "gowebdav"
	"flag"
	"os"
	"fmt"
)

func Fail(err interface{}) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Usage: client FLAGS ARGS")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println("Method:")
		fmt.Println(" LIST, PROPFIND:")
		fmt.Println(" RM, DELETE, DEL:")
		fmt.Println(" MKDIR, MKCOL:")
	}
	os.Exit(-1)
}

func main() {
	root := flag.String("root", "URL", "WebDAV Endpoint")
	usr := flag.String("user", "", "user")
	pw := flag.String("pw", "", "password")
	m := flag.String("X", "GET", "Method ...")
	flag.Parse()

	if *root == "URL" {
		Fail(nil)
	}

	c := d.NewClient(*root, *usr, *pw)
	if err := c.Connect(); err != nil {
		Fail(err)
	}

	if len(flag.Args()) > 0 {
		path := flag.Args()[0]
		switch *m {
			case "LIST", "PROPFIND":
				if files, err := c.ReadDir(path); err == nil {
					fmt.Println(fmt.Sprintf("ReadDir: '%s' entries: %d ", path, len(files)))
					for _, f := range files {
						fmt.Println(f)
					}
				} else {
					fmt.Println(err)
				}

			case "GET": c.Read(path)

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

			default: Fail(nil)
		}
	} else {
		Fail(nil)
	}
}

