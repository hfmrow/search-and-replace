package fileText

import (
	"bytes"
	"fmt"
	"os"
	"unicode/utf8"
)

var binBytes map[byte]bool

// buildBinMap:  build binary bytes map
// https://tools.ietf.org/html/draft-ietf-websec-mime-sniff-03#page-9
func buildBinMap() (binBytes map[byte]bool) {
	bytesRange := [][]byte{[]byte{0x00, 0x08}, []byte{0x0B, 0x0B}, []byte{0x0E, 0x1A}, []byte{0x1C, 0x1F}}
	binBytes = make(map[byte]bool)
	for b := 0; b < 256; b++ {
		for _, rge := range bytesRange {
			if byte(b) >= rge[0] && byte(b) <= rge[1] {
				binBytes[byte(b)] = true
			}
		}
	}
	return binBytes
}

// IsTextFileSimple: Check for text file ... thresholdLineEnd work well with 0.6 and thresholdOverChars with 5
func IsTextFile(filename string, minSizeLimit, maxSizeLimit int64 /*, thresholdLineEnd, thresholdOverChars float64*/) (isText, gtLimit bool, err error) {
	//	var certain bool
	var stat os.FileInfo
	var file *os.File
	// var countedLineEnd float64
	var chunkLength int64 = 1024
	var data []byte
	// var bCR = []byte{0x0D}
	// var bLF = []byte{0x0A}
	// var bCRLF = []byte{0x0D, 0x0A}
	// var charSet string
	// var rx *regexp.Regexp
	// var overChars float64
	// var thres float64

	if binBytes == nil {
		binBytes = buildBinMap()
	}

	if file, err = os.Open(filename); err == nil {
		defer file.Close()
		if stat, err = file.Stat(); err == nil {
			size := stat.Size()
			if size >= minSizeLimit && size <= maxSizeLimit {
				if size < chunkLength { // File size is lower than buffer size, changing buffer size.
					chunkLength = size // to keep proportion in percent formula
				}
				gtLimit = true                   // File is greater than fixed limit
				data = make([]byte, chunkLength) // Reading datas
				if _, err = file.Read(data); err == nil {

					// Try to detect valid utf-8 using golang library unicode/utf8
					if utf8.Valid(data) {
						return true, gtLimit, err
					}

					// // Try to detect charset using golang library html/charset
					// if _, charSet, _ = charset.DetermineEncoding(data, ""); charSet == "utf-8" || charSet == "utf-16" {
					// 	// debug()
					// 	return true, gtLimit, err
					// }

					switch {
					// Text
					case bytes.Contains(data[0:2], []byte{0XFE, 0XFF}): // UTF-16BE BOM
						return true, gtLimit, err
					case bytes.Contains(data[0:2], []byte{0XFE, 0XFE}): // UTF-16LE BOM
						return true, gtLimit, err
					case bytes.Contains(data[0:3], []byte{0XFE, 0XBB, 0XBF}): // UTF-8 BOM
						return true, gtLimit, err
						// Web
					case bytes.Contains(data, []byte(`@charset "UTF-8";`)): // css
						return true, gtLimit, err
					case bytes.Contains(data, []byte(`"Content-Type: text/plain; charset=UTF-8\n"`)): // po-pot
						return true, gtLimit, err
					case bytes.Contains(data, []byte(`<meta charset="utf-8"/>`)): // html-utf-8
						return true, gtLimit, err
					case bytes.Contains(data, []byte(`<?xml version="`)) &&
						bytes.Contains(data, []byte(`" encoding="utf-8"?>`)): // xml-utf-8
						return true, gtLimit, err
					}

					// If previous test fails, we need check for binary Data Byte
					for _, inByte := range data {
						if binBytes[inByte] {
							return false, gtLimit, err
						}
					}
					return true, gtLimit, err

					// overChars = ((float64(overChar) * 100) / float64(chunkLength))
					// if overChars >= thresholdOverChars {
					// 	return false, gtLimit, err
					// } else {
					// 	return true, gtLimit, err
					// }

					// // If previous test fails, must check for line ends count
					// switch {
					// case bytes.Contains(data, bCRLF):
					// 	rx = regexp.MustCompile(".*(" + string(bCRLF) + ")*")
					// case bytes.Contains(data, bCR):
					// 	rx = regexp.MustCompile(".*(" + string(bCR) + ")*")
					// default:
					// 	rx = regexp.MustCompile(".*(" + string(bLF) + ").*")
					// }
					// countedLineEnd = float64(len(rx.FindAllIndex(data, -1)))
					// thres = ((countedLineEnd * 100) / float64(chunkLength))

					// if thres >= thresholdLineEnd {
					// 	return true, gtLimit, err
					// }
				}
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("Check for text file, Unable to open/read file: %s", err.Error())
	}
	return isText, gtLimit, err
}
