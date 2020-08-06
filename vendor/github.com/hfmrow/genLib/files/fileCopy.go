/// fileCopy.go

package files

import (
	"io"
	"os"
	"path/filepath"
)

// FileCopy:
func FileCopy(src, dest string) (err error) {
	var info os.FileInfo
	var f, s *os.File
	if info, err = os.Lstat(src); err == nil {
		if err = os.MkdirAll(filepath.Dir(dest), os.ModePerm); err == nil {
			if f, err = os.Create(dest); err == nil {
				defer f.Close()
				if err = os.Chmod(f.Name(), info.Mode()); err == nil {
					if s, err = os.Open(src); err == nil {
						defer s.Close()
						_, err = io.Copy(f, s)
					}
				}
			}
		}
	}
	return err
}
