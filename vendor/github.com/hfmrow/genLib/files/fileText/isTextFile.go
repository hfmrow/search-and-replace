package fileText

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"unicode/utf8"

	"golang.org/x/net/html/charset"
)

// IsTextFileSimple: Check for text file ... thresholdLineEnd work well with 0.6 and thresholdOverChars with 5
func IsTextFileSimple(filename string, sizeLimit int64, thresholdLineEnd, thresholdOverChars float64) (isText, gtLimit bool, err error) {
	//	var certain bool
	var stat os.FileInfo
	var file *os.File
	var overChar int
	var countedLineEnd float64
	var chunkLength int64 = 1024
	var data []byte
	var bOverChar = []byte{0x7F}
	var bCR = []byte{0x0D}
	var bLF = []byte{0x0A}
	var bCRLF = []byte{0x0D, 0x0A}
	var charSet string
	var rx *regexp.Regexp
	var overChars float64
	var thres float64

	/* DEBUG */
	// var debug = func() {

	// 	TmpCount++
	// 	fmt.Printf("%d - le:%f - oc:%f - %s\n",
	// 		TmpCount,
	// 		thres,
	// 		overChars,
	// 		filename)

	// }
	/* DEBUG */

	// debug()
	// File operations
	if file, err = os.Open(filename); err == nil {
		defer file.Close()
		if stat, err = file.Stat(); err == nil {
			size := stat.Size()
			if size >= sizeLimit {
				if size < chunkLength { // File size is lower than buffer size, changing buffer size.
					chunkLength = size // to keep proportion in percent formula
				}
				gtLimit = true                   // File is greater than fixed limit
				data = make([]byte, chunkLength) // Reading datas
				if _, err = file.Read(data); err == nil {

					// Try o detect charset using golang library unicode/utf8
					if utf8.Valid(data) {
						return true, gtLimit, err
					}

					// Try o detect charset using golang library html/charset
					if _, charSet, _ = charset.DetermineEncoding(data, ""); charSet == "utf-8" || charSet == "utf-16" {
						// debug()
						return true, gtLimit, err
					}

					// If previous test fails, we need to count for ascii > 127 or <10
					for _, inByte := range data {
						b := int(inByte)
						if b < int(bLF[0]) || b > int(bOverChar[0]) {
							overChar++
						}
					}
					overChars = ((float64(overChar) * 100) / float64(chunkLength))
					if overChars >= thresholdOverChars {
						// debug()
						return false, gtLimit, err
					} else {
						// debug()
						return true, gtLimit, err
					}

					// If previous test fails, must check for line ends count
					switch {
					case bytes.Contains(data, bCRLF):
						rx = regexp.MustCompile(".*(" + string(bCRLF) + ")*")
					case bytes.Contains(data, bCR):
						rx = regexp.MustCompile(".*(" + string(bCR) + ")*")
					default:
						rx = regexp.MustCompile(".*(" + string(bLF) + ").*")
					}
					countedLineEnd = float64(len(rx.FindAllIndex(data, -1)))
					thres = ((countedLineEnd * 100) / float64(chunkLength))

					if thres >= thresholdLineEnd {
						// debug()
						return true, gtLimit, err
					}
				}
			}
		}
	}
	if err != nil {
		return isText, gtLimit, errors.New(fmt.Sprintf("Check for text file, Unable to open/read file: %s", err.Error()))
	}
	return isText, gtLimit, err
}
