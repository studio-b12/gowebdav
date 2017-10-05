# GoWebDAV

[![Build Status](https://travis-ci.org/studio-b12/gowebdav.svg?branch=master)](https://travis-ci.org/studio-b12/gowebdav)
[![GoDoc](https://godoc.org/github.com/studio-b12/gowebdav?status.svg)](https://godoc.org/github.com/studio-b12/gowebdav)
[![Go Report Card](https://goreportcard.com/badge/github.com/studio-b12/gowebdav)](https://goreportcard.com/report/github.com/studio-b12/gowebdav)

A golang WebDAV client library and command line tool.

## Install command line tool

```sh
go get -u github.com/studio-b12/gowebdav/cmd/gowebdav
```

## Usage

```sh
$ gowebdav --help
Usage of gowebdav
  -X string
        Method:
                LS <PATH>
                STAT <PATH>

                MKDIR <PATH>
                MKDIRALL <PATH>

                GET <PATH> [<FILE>]
                PUT <PATH> [<FILE>]

                MV <OLD> <NEW>
                CP <OLD> <NEW>

                DEL <PATH>

  -netrc-file string
        read login from netrc file (default "~/.netrc")
  -pw string
        Password [ENV.PASSWORD]
  -root string
        WebDAV Endpoint [ENV.ROOT]
  -user string
        User [ENV.USER] (default "$USER")
```

*gowebdav wrapper script*

Create a wrapper script for example `$EDITOR ./dav && chmod a+x ./dav` for your
server and use [pass](https://www.passwordstore.org/ "the standard unix password manager")
or similar tools to retrieve the password.

```sh
#!/bin/sh

ROOT="https://my.dav.server/" \
USER="foo" \
PASSWORD="$(pass dav/foo@my.dav.server)" \
gowebdav $@
```

*Examples*

Using the `dav` wrapper:

```sh
$ ./dav -X LS /

$ echo hi dav! > hello && ./dav -X PUT /hello

$ ./dav -X STAT /hello

$ ./dav -X PUT /hello_dav hello

$ ./dav -X GET /hello_dav

$ ./dav -X GET /hello_dav hello.txt
```

## LINKS

* [RFC 2518 - HTTP Extensions for Distributed Authoring -- WEBDAV](http://www.faqs.org/rfcs/rfc2518.html "RFC 2518 - HTTP Extensions for Distributed Authoring -- WEBDAV")
* [RFC 2616 - HTTP/1.1 Status Code Definitions](http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html "HTTP/1.1 Status Code Definitions")
* [WebDav: Next Generation Collaborative Web Authoring By Lisa Dusseaul](https://books.google.de/books?isbn=0130652083 "WebDav: Next Generation Collaborative Web Authoring By Lisa Dusseault")

## API

`import "github.com/studio-b12/gowebdav"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)
* [Subdirectories](#pkg-subdirectories)

### <a name="pkg-overview">Overview</a>
Package gowebdav is a WebDAV client library with a command line tool
included.

### <a name="pkg-index">Index</a>
* [func FixSlash(s string) string](#FixSlash)
* [func FixSlashes(s string) string](#FixSlashes)
* [func Join(path0 string, path1 string) string](#Join)
* [func PathEscape(path string) string](#PathEscape)
* [func ReadConfig(uri, netrc string) (string, string)](#ReadConfig)
* [func String(r io.Reader) string](#String)
* [type Authenticator](#Authenticator)
* [type BasicAuth](#BasicAuth)
  * [func (b *BasicAuth) Authorize(c *Client, method string, path string)](#BasicAuth.Authorize)
  * [func (b *BasicAuth) Pass() string](#BasicAuth.Pass)
  * [func (b *BasicAuth) Type() string](#BasicAuth.Type)
  * [func (b *BasicAuth) User() string](#BasicAuth.User)
* [type Client](#Client)
  * [func NewClient(uri, user, pw string) *Client](#NewClient)
  * [func (c *Client) Connect() error](#Client.Connect)
  * [func (c *Client) Copy(oldpath, newpath string, overwrite bool) error](#Client.Copy)
  * [func (c *Client) Mkdir(path string, _ os.FileMode) error](#Client.Mkdir)
  * [func (c *Client) MkdirAll(path string, _ os.FileMode) error](#Client.MkdirAll)
  * [func (c *Client) Read(path string) ([]byte, error)](#Client.Read)
  * [func (c *Client) ReadDir(path string) ([]os.FileInfo, error)](#Client.ReadDir)
  * [func (c *Client) ReadStream(path string) (io.ReadCloser, error)](#Client.ReadStream)
  * [func (c *Client) Remove(path string) error](#Client.Remove)
  * [func (c *Client) RemoveAll(path string) error](#Client.RemoveAll)
  * [func (c *Client) Rename(oldpath, newpath string, overwrite bool) error](#Client.Rename)
  * [func (c *Client) SetHeader(key, value string)](#Client.SetHeader)
  * [func (c *Client) SetTimeout(timeout time.Duration)](#Client.SetTimeout)
  * [func (c *Client) SetTransport(transport http.RoundTripper)](#Client.SetTransport)
  * [func (c *Client) Stat(path string) (os.FileInfo, error)](#Client.Stat)
  * [func (c *Client) Write(path string, data []byte, _ os.FileMode) error](#Client.Write)
  * [func (c *Client) WriteStream(path string, stream io.Reader, _ os.FileMode) error](#Client.WriteStream)
* [type DigestAuth](#DigestAuth)
  * [func (d *DigestAuth) Authorize(c *Client, method string, path string)](#DigestAuth.Authorize)
  * [func (d *DigestAuth) Pass() string](#DigestAuth.Pass)
  * [func (d *DigestAuth) Type() string](#DigestAuth.Type)
  * [func (d *DigestAuth) User() string](#DigestAuth.User)
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
* [type NoAuth](#NoAuth)
  * [func (n *NoAuth) Authorize(c *Client, method string, path string)](#NoAuth.Authorize)
  * [func (n *NoAuth) Pass() string](#NoAuth.Pass)
  * [func (n *NoAuth) Type() string](#NoAuth.Type)
  * [func (n *NoAuth) User() string](#NoAuth.User)

##### <a name="pkg-examples">Examples</a>
* [PathEscape](#example_PathEscape)

##### <a name="pkg-files">Package files</a>
[basicAuth.go](https://github.com/studio-b12/gowebdav/blob/master/basicAuth.go) [client.go](https://github.com/studio-b12/gowebdav/blob/master/client.go) [digestAuth.go](https://github.com/studio-b12/gowebdav/blob/master/digestAuth.go) [doc.go](https://github.com/studio-b12/gowebdav/blob/master/doc.go) [file.go](https://github.com/studio-b12/gowebdav/blob/master/file.go) [netrc.go](https://github.com/studio-b12/gowebdav/blob/master/netrc.go) [requests.go](https://github.com/studio-b12/gowebdav/blob/master/requests.go) [utils.go](https://github.com/studio-b12/gowebdav/blob/master/utils.go) 

### <a name="FixSlash">func</a> [FixSlash](https://github.com/studio-b12/gowebdav/blob/master/utils.go?s=707:737#L45)
``` go
func FixSlash(s string) string
```
FixSlash appends a trailing / to our string

### <a name="FixSlashes">func</a> [FixSlashes](https://github.com/studio-b12/gowebdav/blob/master/utils.go?s=859:891#L53)
``` go
func FixSlashes(s string) string
```
FixSlashes appends and prepends a / if they are missing

### <a name="Join">func</a> [Join](https://github.com/studio-b12/gowebdav/blob/master/utils.go?s=976:1020#L61)
``` go
func Join(path0 string, path1 string) string
```
Join joins two paths

### <a name="PathEscape">func</a> [PathEscape](https://github.com/studio-b12/gowebdav/blob/master/utils.go?s=506:541#L36)
``` go
func PathEscape(path string) string
```
PathEscape escapes all segemnts of a given path

### <a name="ReadConfig">func</a> [ReadConfig](https://github.com/studio-b12/gowebdav/blob/master/netrc.go?s=428:479#L27)
``` go
func ReadConfig(uri, netrc string) (string, string)
```
ReadConfig reads login and password configuration from ~/.netrc
machine foo.com login username password 123456

### <a name="String">func</a> [String](https://github.com/studio-b12/gowebdav/blob/master/utils.go?s=1150:1181#L66)
``` go
func String(r io.Reader) string
```
String pulls a string out of our io.Reader

### <a name="Authenticator">type</a> [Authenticator](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=285:398#L24)
``` go
type Authenticator interface {
    Type() string
    User() string
    Pass() string
    Authorize(*Client, string, string)
}
```
Authenticator stub

### <a name="BasicAuth">type</a> [BasicAuth](https://github.com/studio-b12/gowebdav/blob/master/basicAuth.go?s=94:145#L8)
``` go
type BasicAuth struct {
    // contains filtered or unexported fields
}
```
BasicAuth structure holds our credentials

#### <a name="BasicAuth.Authorize">func</a> (\*BasicAuth) [Authorize](https://github.com/studio-b12/gowebdav/blob/master/basicAuth.go?s=461:529#L29)
``` go
func (b *BasicAuth) Authorize(c *Client, method string, path string)
```
Authorize the current request

#### <a name="BasicAuth.Pass">func</a> (\*BasicAuth) [Pass](https://github.com/studio-b12/gowebdav/blob/master/basicAuth.go?s=376:409#L24)
``` go
func (b *BasicAuth) Pass() string
```
Pass holds the BasicAuth password

#### <a name="BasicAuth.Type">func</a> (\*BasicAuth) [Type](https://github.com/studio-b12/gowebdav/blob/master/basicAuth.go?s=189:222#L14)
``` go
func (b *BasicAuth) Type() string
```
Type identifies the BasicAuthenticator

#### <a name="BasicAuth.User">func</a> (\*BasicAuth) [User](https://github.com/studio-b12/gowebdav/blob/master/basicAuth.go?s=285:318#L19)
``` go
func (b *BasicAuth) User() string
```
User holds the BasicAuth username

### <a name="Client">type</a> [Client](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=157:261#L16)
``` go
type Client struct {
    // contains filtered or unexported fields
}
```
Client defines our structure

#### <a name="NewClient">func</a> [NewClient](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=902:946#L57)
``` go
func NewClient(uri, user, pw string) *Client
```
NewClient creates a new instance of client

#### <a name="Client.Connect">func</a> (\*Client) [Connect](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=1516:1548#L77)
``` go
func (c *Client) Connect() error
```
Connect connects to our dav server

#### <a name="Client.Copy">func</a> (\*Client) [Copy](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=6960:7028#L314)
``` go
func (c *Client) Copy(oldpath, newpath string, overwrite bool) error
```
Copy copies a file from A to B

#### <a name="Client.Mkdir">func</a> (\*Client) [Mkdir](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=6051:6107#L273)
``` go
func (c *Client) Mkdir(path string, _ os.FileMode) error
```
Mkdir makes a directory

#### <a name="Client.MkdirAll">func</a> (\*Client) [MkdirAll](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=6286:6345#L284)
``` go
func (c *Client) MkdirAll(path string, _ os.FileMode) error
```
MkdirAll like mkdir -p, but for webdav

#### <a name="Client.Read">func</a> (\*Client) [Read](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=7134:7184#L319)
``` go
func (c *Client) Read(path string) ([]byte, error)
```
Read reads the contents of a remote file

#### <a name="Client.ReadDir">func</a> (\*Client) [ReadDir](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=3126:3186#L131)
``` go
func (c *Client) ReadDir(path string) ([]os.FileInfo, error)
```
ReadDir reads the contents of a remote directory

#### <a name="Client.ReadStream">func</a> (\*Client) [ReadStream](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=7495:7558#L337)
``` go
func (c *Client) ReadStream(path string) (io.ReadCloser, error)
```
ReadStream reads the stream for a given path

#### <a name="Client.Remove">func</a> (\*Client) [Remove](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=5557:5599#L250)
``` go
func (c *Client) Remove(path string) error
```
Remove removes a remote file

#### <a name="Client.RemoveAll">func</a> (\*Client) [RemoveAll](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=5665:5710#L255)
``` go
func (c *Client) RemoveAll(path string) error
```
RemoveAll removes remote files

#### <a name="Client.Rename">func</a> (\*Client) [Rename](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=6794:6864#L309)
``` go
func (c *Client) Rename(oldpath, newpath string, overwrite bool) error
```
Rename moves a file from A to B

#### <a name="Client.SetHeader">func</a> (\*Client) [SetHeader](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=1099:1144#L62)
``` go
func (c *Client) SetHeader(key, value string)
```
SetHeader lets us set arbitrary headers for a given client

#### <a name="Client.SetTimeout">func</a> (\*Client) [SetTimeout](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=1244:1294#L67)
``` go
func (c *Client) SetTimeout(timeout time.Duration)
```
SetTimeout exposes the ability to set a time limit for requests

#### <a name="Client.SetTransport">func</a> (\*Client) [SetTransport](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=1387:1445#L72)
``` go
func (c *Client) SetTransport(transport http.RoundTripper)
```
SetTransport exposes the ability to define custom transports

#### <a name="Client.Stat">func</a> (\*Client) [Stat](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=4513:4568#L198)
``` go
func (c *Client) Stat(path string) (os.FileInfo, error)
```
Stat returns the file stats for a specified path

#### <a name="Client.Write">func</a> (\*Client) [Write](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=7849:7918#L352)
``` go
func (c *Client) Write(path string, data []byte, _ os.FileMode) error
```
Write writes data to a given path

#### <a name="Client.WriteStream">func</a> (\*Client) [WriteStream](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=8320:8400#L374)
``` go
func (c *Client) WriteStream(path string, stream io.Reader, _ os.FileMode) error
```
WriteStream writes a stream

### <a name="DigestAuth">type</a> [DigestAuth](https://github.com/studio-b12/gowebdav/blob/master/digestAuth.go?s=157:254#L14)
``` go
type DigestAuth struct {
    // contains filtered or unexported fields
}
```
DigestAuth structure holds our credentials

#### <a name="DigestAuth.Authorize">func</a> (\*DigestAuth) [Authorize](https://github.com/studio-b12/gowebdav/blob/master/digestAuth.go?s=577:646#L36)
``` go
func (d *DigestAuth) Authorize(c *Client, method string, path string)
```
Authorize the current request

#### <a name="DigestAuth.Pass">func</a> (\*DigestAuth) [Pass](https://github.com/studio-b12/gowebdav/blob/master/digestAuth.go?s=491:525#L31)
``` go
func (d *DigestAuth) Pass() string
```
Pass holds the DigestAuth password

#### <a name="DigestAuth.Type">func</a> (\*DigestAuth) [Type](https://github.com/studio-b12/gowebdav/blob/master/digestAuth.go?s=299:333#L21)
``` go
func (d *DigestAuth) Type() string
```
Type identifies the DigestAuthenticator

#### <a name="DigestAuth.User">func</a> (\*DigestAuth) [User](https://github.com/studio-b12/gowebdav/blob/master/digestAuth.go?s=398:432#L26)
``` go
func (d *DigestAuth) User() string
```
User holds the DigestAuth username

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

### <a name="NoAuth">type</a> [NoAuth](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=442:490#L32)
``` go
type NoAuth struct {
    // contains filtered or unexported fields
}
```
NoAuth structure holds our credentials

#### <a name="NoAuth.Authorize">func</a> (\*NoAuth) [Authorize](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=785:850#L53)
``` go
func (n *NoAuth) Authorize(c *Client, method string, path string)
```
Authorize the current request

#### <a name="NoAuth.Pass">func</a> (\*NoAuth) [Pass](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=703:733#L48)
``` go
func (n *NoAuth) Pass() string
```
Pass returns the current password

#### <a name="NoAuth.Type">func</a> (\*NoAuth) [Type](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=529:559#L38)
``` go
func (n *NoAuth) Type() string
```
Type identifies the authenticator

#### <a name="NoAuth.User">func</a> (\*NoAuth) [User](https://github.com/studio-b12/gowebdav/blob/master/client.go?s=615:645#L43)
``` go
func (n *NoAuth) User() string
```
User returns the current user

- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
