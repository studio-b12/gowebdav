package gowebdav

import (
	"fmt"
	"time"
)

type File interface {
	Name() string
	Size() uint
	Modified() time.Time
	IsDirectory() bool
	String() string
}

type file struct {
	name     string
	size     uint
	modified time.Time
}

func (_ file) IsDirectory() bool {
	return false
}

func (f file) Modified() time.Time {
	return f.modified
}

func (f file) Name() string {
	return f.name
}

func (f file) Size() uint {
	return f.size
}

func (f file) String() string {
	return fmt.Sprintf("FILE: %s SIZE: %d MODIFIED: %s", f.name, f.size, f.modified.String())
}
