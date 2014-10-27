package gowebdav

import (
	"fmt"
	"os"
	"time"
)

type File struct {
	path     string
	name     string
	size     int64
	modified time.Time
	isdir    bool
}

func (f File) Name() string {
	return f.name
}

func (f File) Size() int64 {
	return f.size
}

func (f File) Mode() os.FileMode {
	// TODO check webdav perms
	if f.isdir {
		return 0775 | os.ModeDir
	} else {
		return 0664
	}
}

func (f File) ModTime() time.Time {
	return f.modified
}

func (f File) IsDir() bool {
	return f.isdir
}

func (f File) Sys() interface{} {
	return nil
}

func (f File) String() string {
	if f.isdir {
		return fmt.Sprintf("Dir : '%s' - '%s'", f.path, f.name)
	} else {
		return fmt.Sprintf("File: '%s' SIZE: %d MODIFIED: %s", f.path, f.size, f.modified.String())
	}
}
