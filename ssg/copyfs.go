package ssg

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// like https://github.com/golang/go/issues/62484#issue-1884498794 but with error handling, custom walk root, and follows symlinks
func CopyFS(dst string, fsys fs.FS, fspath string) error {
	return fs.WalkDir(fsys, fspath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// follow symlink
		var isDir = d.IsDir()
		if d.Type()&fs.ModeSymlink != 0 {
			// get symlink target FileInfo with fs.Stat
			info, err := fs.Stat(fsys, filepath.Join(fspath, d.Name()))
			if err == nil {
				if info.Mode()&fs.ModeDir != 0 {
					isDir = true
				}
			}
		}

		targ := filepath.Join(dst, filepath.FromSlash(path))
		if isDir {
			if err := os.MkdirAll(targ, 0777); err != nil {
				return err
			}
			return nil
		}
		r, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()
		info, err := r.Stat()
		if err != nil {
			return err
		}
		w, err := os.OpenFile(targ, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666|info.Mode()&0777)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, r); err != nil {
			w.Close()
			return fmt.Errorf("copying %s: %v", path, err)
		}
		if err := w.Close(); err != nil {
			return err
		}
		return nil
	})
}
