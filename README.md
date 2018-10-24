# GoWebDAV

[![Build Status](https://travis-ci.org/studio-b12/gowebdav.svg?branch=master)](https://travis-ci.org/studio-b12/gowebdav)
[![GoDoc](https://godoc.org/github.com/studio-b12/gowebdav?status.svg)](https://godoc.org/github.com/studio-b12/gowebdav)
[![Go Report Card](https://goreportcard.com/badge/github.com/studio-b12/gowebdav)](https://goreportcard.com/report/github.com/studio-b12/gowebdav)

A golang WebDAV client library.

## Main features
`gowebdav` library allows to perform following actions on the remote WebDAV server:
* [create path](#create-path-on-a-webdav-server)
* [get files list](#get-files-list)
* [download file](#download-file-to-byte-array)
* [upload file](#upload-file-from-byte-array)
* [get information about specified file/folder](#get-information-about-specified-filefolder)
* [move file to another location](#move-file-to-another-location)
* [copy file to another location](#copy-file-to-another-location)
* [delete file](#delete-file)

## Usage

First of all you should create `Client` instance using `NewClient()` function:

```go
root := "https://webdav.mydomain.me"
user := "user"
password := "password"

c := gowebdav.NewClient(root, user, password)
```

After you can use this `Client` to perform actions, described below.

**NOTICE:** we will not check errors in examples, to focus you on the `gowebdav` library's code, but you should do it in your code!

### Create path on a WebDAV server
```go
err := c.Mkdir("folder", 0644)
```
In case you want to create several folders you can use `c.MkdirAll()`:
```go
err := c.MkdirAll("folder/subfolder/subfolder2", 0644)
```

### Get files list
```go
files, _ := c.ReadDir("folder/subfolder")
for _, file := range files {
    //notice that [file] has os.FileInfo type
    fmt.Println(file.Name())
}
```

### Download file to byte array
```go
webdavFilePath := "folder/subfolder/file.txt"
localFilePath := "/tmp/webdav/file.txt"

bytes, _ := c.Read(webdavFilePath)
ioutil.WriteFile(localFilePath, bytes, 0644)
```

### Download file via reader
Also you can use `c.ReadStream()` method:
```go
webdavFilePath := "folder/subfolder/file.txt"
localFilePath := "/tmp/webdav/file.txt"

reader, _ := c.ReadStream(webdavFilePath)

file, _ := os.Create(localFilePath)
defer file.Close()

io.Copy(file, reader)
```

### Upload file from byte array
```go
webdavFilePath := "folder/subfolder/file.txt"
localFilePath := "/tmp/webdav/file.txt"

bytes, _ := ioutil.ReadFile(localFilePath)

c.Write(webdavFilePath, bytes, 0644)
```

### Upload file via writer
```go
webdavFilePath := "folder/subfolder/file.txt"
localFilePath := "/tmp/webdav/file.txt"

file, _ := os.Open(localFilePath)
defer file.Close()

c.WriteStream(webdavFilePath, file, 0644)
```

### Get information about specified file/folder
```go
webdavFilePath := "folder/subfolder/file.txt"

info := c.Stat(webdavFilePath)
//notice that [info] has os.FileInfo type
fmt.Println(info)
```

### Move file to another location
```go
oldPath := "folder/subfolder/file.txt"
newPath := "folder/subfolder/moved.txt"
isOverwrite := true

c.Rename(oldPath, newPath, isOverwrite)
```

### Copy file to another location
```go
oldPath := "folder/subfolder/file.txt"
newPath := "folder/subfolder/file-copy.txt"
isOverwrite := true

c.Copy(oldPath, newPath, isOverwrite)
```

### Delete file
```go
webdavFilePath := "folder/subfolder/file.txt"

c.Remove(webdavFilePath)
```

## Links

More details about WebDAV server you can read from following resources:

* [RFC 4918 - HTTP Extensions for Web Distributed Authoring and Versioning (WebDAV)](https://tools.ietf.org/html/rfc4918)
* [RFC 5689 - Extended MKCOL for Web Distributed Authoring and Versioning (WebDAV)](https://tools.ietf.org/html/rfc5689)
* [RFC 2616 - HTTP/1.1 Status Code Definitions](http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html "HTTP/1.1 Status Code Definitions")
* [WebDav: Next Generation Collaborative Web Authoring By Lisa Dusseaul](https://books.google.de/books?isbn=0130652083 "WebDav: Next Generation Collaborative Web Authoring By Lisa Dusseault")

**NOTICE**: RFC 2518 is obsoleted by RFC 4918 in June 2007

## Contributing
All contributing are welcome. If you have any suggestions or find some bug - please create an Issue to let us make this project better. We appreciate your help!

## License
This library is distributed under the BSD 3-Clause license found in the [LICENSE](https://github.com/studio-b12/gowebdav/blob/master/LICENSE) file.
