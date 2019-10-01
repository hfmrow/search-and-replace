package files

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// IsDirEmpty:
func IsDirEmpty(name string) (empty bool, err error) {
	var f *os.File
	if f, err = os.Open(name); err == nil {
		defer f.Close()
		if _, err = f.Readdirnames(1); err == io.EOF {
			return true, nil
		}
	}

	return false, err
}

// GetCurrentDir: Get current directory
func GetCurrentDir() (dir string, err error) {
	return os.Getwd()
}

// TempMake: Make temporary directory
func TempMake(prefix string) string {
	dir, err := ioutil.TempDir("", prefix+"-")
	if err != nil {
		log.Fatalf("Unablme to create temp directory: %s\n", err.Error())
	}
	return dir + string(os.PathSeparator)
}

// TempRemove: Remove directory recursively
func TempRemove(fName string) (err error) {
	if err = os.RemoveAll(fName); err != nil {
		return (err)
	}
	return nil
}

// ExtEnsure: ensure the filename have desired extension
func ExtEnsure(filename, ext string) (outFilename string) {
	outFilename = filename
	if !strings.HasSuffix(filename, ext) {
		currExt := path.Ext(filename)
		outFilename = filename[:len(filename)-len(currExt)] + ext
	}
	return outFilename
}

// BaseNoExt: get only the name without ext.
func BaseNoExt(filename string) (outFilename string) {
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filepath.Base(filename)))
}

// magic number mime detection
var magicTable = map[string]string{
	"\x37\x7A\xBC\xAF\x27\x1C\x00\x04": "7zip",
	"\xFD\x37\x7A\x58\x5A\x00\x00":     "xz",
	"\x1F\x8B\x08\x00\x00\x09\x6E\x88": "gzip",
	"\x75\x73\x74\x61\x72":             "tar",
}

// GetFileMime: scan first bytes to detect mime type of file.
func GetFileMime(filename string) string {
	if file, err := os.Open(filename); err == nil {
		defer file.Close()
		buffReader := bufio.NewReader(file)
		for magic, mime := range magicTable {
			if peeked, err := buffReader.Peek(len([]byte(magic))); err == nil {
				tmpMagic := []byte(magic)
				if bytes.Index(peeked, tmpMagic) == 0 {
					return mime
				}
			}
		}
	}
	return "Unknown"
}
