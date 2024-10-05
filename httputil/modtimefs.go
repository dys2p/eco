package httputil

import (
	"embed"
	"io"
	"io/fs"
	"time"
)

// A ModTimeFS is an embed.FS with a given modification time. Else embed.FS returns an empty ModTime, resulting in http.FileServer not setting a Last-Modified header.
type ModTimeFS struct {
	embed.FS
	ModTime time.Time
}

func (fs ModTimeFS) Open(name string) (fs.File, error) {
	f, err := fs.FS.Open(name)
	if err != nil {
		return f, err
	}
	return file{f, fs.ModTime}, nil
}

type file struct {
	fs.File
	modTime time.Time
}

// embed.FS.Open implements io.ReaderAt
func (f file) ReadAt(p []byte, off int64) (n int, err error) {
	return f.File.(io.ReaderAt).ReadAt(p, off)
}

// embed.FS.Open implements io.Seeker
func (f file) Seek(offset int64, whence int) (int64, error) {
	return f.File.(io.Seeker).Seek(offset, whence)
}

func (f file) Stat() (fs.FileInfo, error) {
	info, err := f.File.Stat()
	if err != nil {
		return info, err
	}
	return fileInfo{info, f.modTime}, nil
}

type fileInfo struct {
	fs.FileInfo
	modTime time.Time
}

func (info fileInfo) ModTime() time.Time {
	return info.modTime
}
