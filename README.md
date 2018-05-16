# GoWebDAV

[![Build Status](https://travis-ci.org/studio-b12/gowebdav.svg?branch=master)](https://travis-ci.org/studio-b12/gowebdav)
[![Go Report Card](https://goreportcard.com/badge/github.com/studio-b12/gowebdav)](https://goreportcard.com/report/github.com/studio-b12/gowebdav)

A WebDAV client and library for golang.

## Install

```sh
go get -u github.com/studio-b12/gowebdav/cmd/gowebdav
```

## Usage

```sh
$ gowebdav
Usage: gowebdav FLAGS ARGS
Flags:
  -X string
        Method ... (default "GET")
  -pw string
        password
  -root string
        WebDAV Endpoint (default "URL")
  -user string
        user
Method <ARGS>
 LS | LIST | PROPFIND <PATH>
 RM | DELETE | DEL <PATH>
 MKDIR | MKCOL <PATH>
 MKDIRALL | MKCOLALL <PATH>
 MV | MOVE | RENAME <OLD_PATH> <NEW_PATH>
 CP | COPY <OLD_PATH> <NEW_PATH>
 GET | PULL | READ <PATH>
 PUT | PUSH | WRITE <PATH> <FILE>
 STAT <PATH>
```

## LINKS

* [RFC 2518 - HTTP Extensions for Distributed Authoring -- WEBDAV](http://www.faqs.org/rfcs/rfc2518.html "RFC 2518 - HTTP Extensions for Distributed Authoring -- WEBDAV")
* [RFC 2616 - HTTP/1.1 Status Code Definitions](http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html "HTTP/1.1 Status Code Definitions")
* [WebDav: Next Generation Collaborative Web Authoring By Lisa Dusseaul](https://books.google.de/books?isbn=0130652083 "WebDav: Next Generation Collaborative Web Authoring By Lisa Dusseault")

## API

`import "github.com/studio-b12/gowebdav"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

### <a name="pkg-overview">Overview</a>
Package gowebdav A golang WebDAV library

### <a name="pkg-index">Index</a>
* [func FixSlash(s string) string](#FixSlash)
* [func FixSlashes(s string) string](#FixSlashes)
* [func Join(path0 string, path1 string) string](#Join)
* [func String(r io.Reader) string](#String)
* [type Client](#Client)
  * [func NewClient(uri string, user string, pw string) *Client](#NewClient)
  * [func (c *Client) Connect() error](#Client.Connect)
  * [func (c *Client) Copy(oldpath string, newpath string, overwrite bool) error](#Client.Copy)
  * [func (c *Client) Mkdir(path string, _ os.FileMode) error](#Client.Mkdir)
  * [func (c *Client) MkdirAll(path string, _ os.FileMode) error](#Client.MkdirAll)
  * [func (c *Client) Read(path string) ([]byte, error)](#Client.Read)
  * [func (c *Client) ReadDir(path string) ([]os.FileInfo, error)](#Client.ReadDir)
  * [func (c *Client) ReadStream(path string) (io.ReadCloser, error)](#Client.ReadStream)
  * [func (c *Client) Remove(path string) error](#Client.Remove)
  * [func (c *Client) RemoveAll(path string) error](#Client.RemoveAll)
  * [func (c *Client) Rename(oldpath string, newpath string, overwrite bool) error](#Client.Rename)
  * [func (c *Client) SetHeader(key, value string)](#Client.SetHeader)
  * [func (c *Client) SetTimeout(timeout time.Duration)](#Client.SetTimeout)
  * [func (c *Client) SetTransport(transport http.RoundTripper)](#Client.SetTransport)
  * [func (c *Client) Stat(path string) (os.FileInfo, error)](#Client.Stat)
  * [func (c *Client) Write(path string, data []byte, _ os.FileMode) error](#Client.Write)
  * [func (c *Client) WriteStream(path string, stream io.Reader, _ os.FileMode) error](#Client.WriteStream)
* [type File](#File)
  * [func (f File) ContentType() string](#File.ContentType)
  * [func (f File) ETag() string](#File.ETag)
  * [func (f File) IsDir() bool](#File.IsDir)
  * [func (f File) ModTime() time.Time](#File.ModTime)
  * [func (f File) Mode() os.FileMode](#File.Mode)
  * [func (f File) Name() string](#File.Name)
  * [func (f File) Size() int64](#File.Size)
  * [func (f File) String() string](#File.String)
  * [func (f File) Sys() interface{}](#File.Sys)

##### <a name="pkg-files">Package files</a>
[client.go](https://github.com/studio-b12/gowebdav/blob/master/client.go) [file.go](https://github.com/studio-b12/gowebdav/blob/master/file.go) [requests.go](https://github.com/studio-b12/gowebdav/blob/master/requests.go) [utils.go](https://github.com/studio-b12/gowebdav/blob/master/utils.go) 

### <a name="FixSlash">func</a> [FixSlash](https://github.com/studio-b12/gowebdav/blob/master/utils.go?s=491:521#L35)
``` go
func FixSlash(s string) string
```
FixSlash appends a trailing / to our string

### <a name="FixSlashes">func</a> [FixSlashes](https://github.com/studio-b12/gowebdav/blob/master/utils.go?s=643:675#L43)
``` go
func FixSlashes(s string) string
```
FixSlashes appends and prepends a / if they are missing

### <a name="Join">func</a> [Join](https://github.com/studio-b12/gowebdav/blob/master/utils.go?s=760:804#L51)
``` go
func Join(path0 string, path1 string) string
```
Join joins two paths

### <a name="String">func</a> [String](https://github.com/studio-b12/gowebdav/blob/master/utils.go?s=934:965#L56)
``` go
func String(r io.Reader) string
```
String pulls a string out of our io.Reader

### <a name="Client">type</a> [Client](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=220:301#L18)
``` go
type Client struct {
    // contains filtered or unexported fields
}
```
Client defines our structure

#### <a name="NewClient">func</a> [NewClient](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=349:407#L25)
``` go
func NewClient(uri string, user string, pw string) *Client
```
NewClient creates a new instance of client

#### <a name="Client.Connect">func</a> (\*Client) [Connect](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=1152:1184#L55)
``` go
func (c *Client) Connect() error
```
Connect connects to our dav server

#### <a name="Client.Copy">func</a> (\*Client) [Copy](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=6149:6224#L289)
``` go
func (c *Client) Copy(oldpath string, newpath string, overwrite bool) error
```
Copy copies a file from A to B

#### <a name="Client.Mkdir">func</a> (\*Client) [Mkdir](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=5233:5289#L248)
``` go
func (c *Client) Mkdir(path string, _ os.FileMode) error
```
Mkdir makes a directory

#### <a name="Client.MkdirAll">func</a> (\*Client) [MkdirAll](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=5468:5527#L259)
``` go
func (c *Client) MkdirAll(path string, _ os.FileMode) error
```
MkdirAll like mkdir -p, but for webdav

#### <a name="Client.Read">func</a> (\*Client) [Read](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=6330:6380#L294)
``` go
func (c *Client) Read(path string) ([]byte, error)
```
Read reads the contents of a remote file

#### <a name="Client.ReadDir">func</a> (\*Client) [ReadDir](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=2240:2300#L98)
``` go
func (c *Client) ReadDir(path string) ([]os.FileInfo, error)
```
ReadDir reads the contents of a remote directory

#### <a name="Client.ReadStream">func</a> (\*Client) [ReadStream](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=6691:6754#L312)
``` go
func (c *Client) ReadStream(path string) (io.ReadCloser, error)
```
ReadStream reads the stream for a given path

#### <a name="Client.Remove">func</a> (\*Client) [Remove](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=4739:4781#L225)
``` go
func (c *Client) Remove(path string) error
```
Remove removes a remote file

#### <a name="Client.RemoveAll">func</a> (\*Client) [RemoveAll](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=4847:4892#L230)
``` go
func (c *Client) RemoveAll(path string) error
```
RemoveAll removes remote files

#### <a name="Client.Rename">func</a> (\*Client) [Rename](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=5976:6053#L284)
``` go
func (c *Client) Rename(oldpath string, newpath string, overwrite bool) error
```
Rename moves a file from A to B

#### <a name="Client.SetHeader">func</a> (\*Client) [SetHeader](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=735:780#L40)
``` go
func (c *Client) SetHeader(key, value string)
```
SetHeader lets us set arbitrary headers for a given client

#### <a name="Client.SetTimeout">func</a> (\*Client) [SetTimeout](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=880:930#L45)
``` go
func (c *Client) SetTimeout(timeout time.Duration)
```
SetTimeout exposes the ability to set a time limit for requests

#### <a name="Client.SetTransport">func</a> (\*Client) [SetTransport](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=1023:1081#L50)
``` go
func (c *Client) SetTransport(transport http.RoundTripper)
```
SetTransport exposes the ability to define custom transports

#### <a name="Client.Stat">func</a> (\*Client) [Stat](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=3661:3716#L169)
``` go
func (c *Client) Stat(path string) (os.FileInfo, error)
```
Stat returns the file stats for a specified path

#### <a name="Client.Write">func</a> (\*Client) [Write](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=7045:7114#L327)
``` go
func (c *Client) Write(path string, data []byte, _ os.FileMode) error
```
Write writes data to a given path

#### <a name="Client.WriteStream">func</a> (\*Client) [WriteStream](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=7516:7596#L349)
``` go
func (c *Client) WriteStream(path string, stream io.Reader, _ os.FileMode) error
```
WriteStream writes a stream

### <a name="File">type</a> [File](https://github.com/studio-b12/gowebdav/blob/master/file.go?s=93:253#L10)
``` go
type File struct {
    // contains filtered or unexported fields
}
```
File is our structure for a given file

#### <a name="File.ContentType">func</a> (File) [ContentType](https://github.com/studio-b12/gowebdav/blob/master/file.go?s=388:422#L26)
``` go
func (f File) ContentType() string
```
ContentType returns the content type of a file

#### <a name="File.ETag">func</a> (File) [ETag](https://github.com/studio-b12/gowebdav/blob/master/file.go?s=841:868#L51)
``` go
func (f File) ETag() string
```
ETag returns the ETag of a file

#### <a name="File.IsDir">func</a> (File) [IsDir](https://github.com/studio-b12/gowebdav/blob/master/file.go?s=947:973#L56)
``` go
func (f File) IsDir() bool
```
IsDir let us see if a given file is a directory or not

#### <a name="File.ModTime">func</a> (File) [ModTime](https://github.com/studio-b12/gowebdav/blob/master/file.go?s=748:781#L46)
``` go
func (f File) ModTime() time.Time
```
ModTime returns the modified time of a file

#### <a name="File.Mode">func</a> (File) [Mode](https://github.com/studio-b12/gowebdav/blob/master/file.go?s=577:609#L36)
``` go
func (f File) Mode() os.FileMode
```
Mode will return the mode of a given file

#### <a name="File.Name">func</a> (File) [Name](https://github.com/studio-b12/gowebdav/blob/master/file.go?s=290:317#L21)
``` go
func (f File) Name() string
```
Name returns the name of a file

#### <a name="File.Size">func</a> (File) [Size](https://github.com/studio-b12/gowebdav/blob/master/file.go?s=485:511#L31)
``` go
func (f File) Size() int64
```
Size returns the size of a file

#### <a name="File.String">func</a> (File) [String](https://github.com/studio-b12/gowebdav/blob/master/file.go?s=1095:1124#L66)
``` go
func (f File) String() string
```
String lets us see file information

#### <a name="File.Sys">func</a> (File) [Sys](https://github.com/studio-b12/gowebdav/blob/master/file.go?s=1007:1038#L61)
``` go
func (f File) Sys() interface{}
```
Sys ????

- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
