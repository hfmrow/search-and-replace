// eol.go

/*
* Some end of line functions used by different OS
 */

package strings

import (
	"bytes"
	"errors"
	"io/ioutil"
	"runtime"
	"strings"
)

var platforms = [][]string{
	{"darwin", "386", "\n"},
	{"darwin", "amd64", "\n"},
	{"dragonfly", "amd64", "\n"},
	{"freebsd", "386", "\n"},
	{"freebsd", "amd64", "\n"},
	{"freebsd", "arm", "\n"},
	{"linux", "386", "\n"},
	{"linux", "amd64", "\n"},
	{"linux", "arm", "\n"},
	{"linux", "arm64", "\n"},
	{"linux", "ppc64", "\n"},
	{"linux", "ppc64le", "\n"},
	{"linux", "mips", "\n"},
	{"linux", "mipsle", "\n"},
	{"linux", "mips64", "\n"},
	{"linux", "mips64le", "\n"},
	{"linux", "s390x", "\n"},
	{"nacl", "386", "\n"},
	{"nacl", "amd64p32", "\n"},
	{"nacl", "arm", "\n"},
	{"netbsd", "386", "\n"},
	{"netbsd", "amd64", "\n"},
	{"netbsd", "arm", "\n"},
	{"openbsd", "386", "\n"},
	{"openbsd", "amd64", "\n"},
	{"openbsd", "arm", "\n"},
	{"plan9", "386", "\n"},
	{"plan9", "amd64", "\n"},
	{"plan9", "arm", "\n"},
	{"solaris", "amd64", "\n"},
	{"windows", "386", "\r\n"},
	{"windows", "amd64", "\r\n"}}

// GetOsLineEnd: Get current OS line-feed
func GetOsLineEnd() string {
	for _, row := range platforms {
		if row[0] == runtime.GOOS {
			return row[2]
		}
	}
	return "\n"
}

// GetTextEOL: Get EOL from text bytes (CR, LF, CRLF)
func GetTextEOL(inTextBytes []byte) (outString string) {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	switch {
	case bytes.Contains(inTextBytes, bCRLF):
		outString = string(bCRLF)
	case bytes.Contains(inTextBytes, bCR):
		outString = string(bCR)
	default:
		outString = string(bLF)
	}
	return
}

// SetTextEOL: Get EOL from text bytes and convert it to another EOL (CR, LF, CRLF)
func SetTextEOL(inTextBytes []byte, eol string) (outTextBytes []byte, err error) {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	var outEol []byte
	switch eol {
	case "CR":
		outEol = bCR
	case "LF":
		outEol = bLF
	case "CRLF":
		outEol = bCRLF
	default:
		return outTextBytes, errors.New("EOL convert error: Undefined end of line")
	}
	// Handle end of line
	outTextBytes = bytes.Replace(inTextBytes, []byte(GetTextEOL(inTextBytes)), outEol, -1)
	return outTextBytes, nil
}

// SplitAtEOL: split data to slice
func SplitAtEOL(data []byte) (outSlice []string) {
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	data = bytes.ReplaceAll(data, bCRLF, bLF)
	return strings.Split(string(data), string(bLF))
}

// GetFileEOL: Open file and get (CR, LF, CRLF) > string or get OS line end.
func GetFileEOL(filename string) (outString string, err error) {
	textFileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return outString, err
	}
	return GetTextEOL(textFileBytes), nil
}

// SetFileEOL: Open file and convert EOL (CR, LF, CRLF) then write it back.
func SetFileEOL(filename, eol string) error {
	textFileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	// Handle end of line
	textFileBytes, err = SetTextEOL(textFileBytes, eol)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, textFileBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
