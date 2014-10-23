package gowebdav

import (
	"fmt"
	"time"
)

type Directory interface {
}

type directory struct {
	name string
}

func (d directory) Name() string {
	return d.name
}

func (_ directory) Size() uint {
	return 0
}

func (_ directory) IsDirectory() bool {
	return true
}

func (_ directory) Modified() time.Time {
	return time.Unix(0, 9)
}
func (d directory) String() string {
	return fmt.Sprintf("DIRECTORY: %s", d.name)
}
