// searchAndReplace.go

/*
	Copyright ©2018-19 H.F.M - Search And Replace library
	This program comes with absolutely no warranty. See the
	The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package fileText

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	glco "github.com/hfmrow/genLib/crypto"
	glfs "github.com/hfmrow/genLib/files"

	glsg "github.com/hfmrow/genLib/strings"
	glsscc "github.com/hfmrow/genLib/strings/cClass"
)

// SearchAndReplaceFiles: is a structure that hold some methods
// to provide an efficient way to search, replace given pattern
// in text files. There is a lot of options to perform personalized
// research.
type SearchAndReplaceFiles struct {
	FileName      string
	SearchAndRepl SearchAndReplace
	NotTextFile   bool
	Occurrences   int
}

// SearchAndReplaceInMultipleFiles: Search in multiples text files.
// return a slice type []SearchAndReplaceInFiles that contain all
// informations about found patterns, indexes, lines position, file
// type, size and occurances count.
func SearchAndReplaceInFiles(filenames []string, toSearch, replaceWith string, /*thresholdLineEnd, thresholdOverChars float64,*/
	minSizeLimit, maxSizeLimit int64, caseSensitive, posixCharClass, posixStrictMode, regex, wildcard, useEscapeChar, wholeWord,
	doReplace, doSave, doBackup, acceptBinary, removeEmptyResult bool) (founds []SearchAndReplaceFiles, occurFound int, err error) {

	var stat os.FileInfo
	founds = make([]SearchAndReplaceFiles, len(filenames))

	for idxFile, file := range filenames {

		// check for File exist and is not a directory
		if stat, err = os.Stat(file); !os.IsNotExist(err) && !stat.IsDir() {

			founds[idxFile].FileName = file
			// Check for text file
			isTxt, gtSizeLimit, err := IsTextFile(
				file,
				minSizeLimit,
				maxSizeLimit, /*
					thresholdLineEnd,
					thresholdOverChars*/)

			if (!gtSizeLimit || (!acceptBinary && !isTxt)) && err == nil {
				// If it's a binary file & not allowed or size is lower than requested
				founds[idxFile].NotTextFile = true
				if !gtSizeLimit {
					founds[idxFile].SearchAndRepl.TextBytes = []byte(fmt.Sprintf("Files size < %d or > %d", minSizeLimit, maxSizeLimit))
				} else {
					founds[idxFile].SearchAndRepl.TextBytes = []byte("Binary content") // Put type of file in TextBytes field
				}
				// Adding a fake line to keep this entry
				founds[idxFile].SearchAndRepl.Pos.FoundLinesIdx = append(founds[idxFile].SearchAndRepl.Pos.FoundLinesIdx, lineIdxInf{Number: -1})
			} else {
				textBytes, err := ioutil.ReadFile(file)
				if err != nil {
					return founds, occurFound, err
				}

				founds[idxFile].SearchAndRepl = *SearchAndReplaceNew([]byte{}, "", "")
				founds[idxFile].SearchAndRepl.Init(
					textBytes,
					toSearch,
					replaceWith,
					caseSensitive,
					posixCharClass,
					posixStrictMode,
					regex,
					wildcard,
					useEscapeChar,
					wholeWord,
					doReplace)

				if err = founds[idxFile].SearchAndRepl.SearchAndReplace(); err != nil {
					return founds, occurFound, err
				}
				founds[idxFile].Occurrences = founds[idxFile].SearchAndRepl.Occurrences
				occurFound += founds[idxFile].Occurrences
				// Saving file if one or more modifications was done
				if doSave && doReplace && founds[idxFile].Occurrences > 0 /*&& founds[idxFile].SearchAndRepl.ReadyToReplace()*/ {
					err = glfs.WriteFile(founds[idxFile].FileName, founds[idxFile].SearchAndRepl.TextBytes, doBackup)
					if err != nil {
						return founds, occurFound, err
					}
				}
			}
		}
	}
	// Removing empty structures if requested
	if removeEmptyResult {
		for idx := len(founds) - 1; idx >= 0; idx-- {
			if len(founds[idx].SearchAndRepl.Pos.FoundLinesIdx) == 0 {
				founds = append(founds[:idx], founds[idx+1:]...)
			}
		}
	}
	return founds, occurFound, err
}

// SearchAndReplace: is a structure that hold some methods to
// provide an efficient way to search, replace given pattern
// in text. There is a lot of options to perform personalized
// research.
type SearchAndReplace struct {
	TextBytes       []byte
	ToSearch        string
	ToSearchRegexp  *regexp.Regexp
	ReplaceWith     string
	CaseSensitive   bool
	UseEscapeChar   bool
	PosixCharClass  bool
	PosixStrictMode bool
	Regex           bool
	Wildcard        bool
	WholeWord       bool
	DoReplace       bool
	TextBytesMd5    string
	Pos             LinesInfos
	Occurrences     int
	OnEachLine      func(idx, lineStart, lineEnd int)

	// Used to define if something has been changed in the
	// parameters, that permit to avoid useless computation
	readyToReplace bool // Access via method, read only ...

	// I use this variable to now whether a display have
	// been done in the parent application.
	hasBeenDisplayed bool // Access via method, read/write ...
}

// SearchAndReplaceNew: Cre	at new "SearchAndReplace" structure
// with short defaul parameters, for the case of single speed search.
func SearchAndReplaceNew(textBytes []byte, toSearch, replaceWith string) (s *SearchAndReplace) {

	s = new(SearchAndReplace)
	if len(textBytes) > 0 {
		s.Init(textBytes, toSearch, replaceWith, true, false,
			false, false, false, false, false, false)
	}
	return
}

// InitFull: do a complete initialization a "SearchAndReplace"
// structure with given parameters.
func (s *SearchAndReplace) Init(textBytes []byte, toSearch, replaceWith string, caseSensitive, posixCharClass,
	posixStrictMode, regex, wildcard, escapeChar, wholeWord, doReplace bool) (err error) {

	s.CaseSensitive = caseSensitive
	s.PosixCharClass = posixCharClass
	s.PosixStrictMode = posixStrictMode
	s.Regex = regex
	s.Wildcard = wildcard
	s.UseEscapeChar = escapeChar
	s.WholeWord = wholeWord
	s.DoReplace = doReplace

	// If something has changed from the last run, ReadyToReplace()
	// will return false.
	return s.compareEntries(textBytes, toSearch, replaceWith)
}

// ReadyToReplace: return variable content
func (s *SearchAndReplace) ReadyToReplace() bool {
	return s.readyToReplace
}

// HasBeenDisplayed: return/set variable content
func (s *SearchAndReplace) HasBeenDisplayed(set ...bool) bool {
	if len(set) > 0 {
		s.hasBeenDisplayed = set[0]
	}
	return s.hasBeenDisplayed
}

// ReadyToReplace: return variable content
func (s *SearchAndReplace) Reset() {
	s.readyToReplace = false
	s.hasBeenDisplayed = false
	s.Pos = LinesInfos{}
}

// replace: check for and do it if ok.
func (s *SearchAndReplace) replace() bool {
	if s.readyToReplace {
		s.TextBytes = []byte(s.ToSearchRegexp.ReplaceAllString(string(s.TextBytes), s.ReplaceWith))
		s.DoReplace = false
		return true
	}
	return false
}

func (s *SearchAndReplace) compareEntries(textBytes []byte, toSearch, replaceWith string) (err error) {

	var tmpReg *regexp.Regexp

	// Build regexp to compare with previous
	if tmpReg, err = BuildRegexp(toSearch, s.CaseSensitive, s.PosixCharClass,
		s.PosixStrictMode, s.Regex, s.Wildcard, s.UseEscapeChar, s.WholeWord); err != nil {
		s.readyToReplace = false
		return
	}
	// Compare with previous arguments
	currentMd5 := glco.Md5String(string(textBytes))

	if s.ToSearchRegexp != nil {
		switch {
		case tmpReg.String() != s.ToSearchRegexp.String():
			s.Reset()
		case s.TextBytesMd5 != currentMd5:
			s.Reset()
		case s.ReplaceWith != replaceWith:
			s.Reset()
		}
	} /*else {
		s.ToSearchRegexp = tmpReg
	}*/
	s.TextBytesMd5 = currentMd5
	s.ToSearchRegexp = tmpReg
	s.TextBytes = textBytes
	s.ToSearch = toSearch
	s.ReplaceWith = replaceWith
	return
}

// Search in plain text, use "Init" to fill needed
// informations about search preferences before using ...
func (s *SearchAndReplace) SearchAndReplace() (err error) {

	if !s.readyToReplace {

		if len(s.ToSearch) > 0 {

			// Set on_each_line function
			s.Pos.onEachLine = s.OnEachLine

			if !s.readyToReplace {
				// Do the search/Replace job
				if patternPos := s.ToSearchRegexp.FindAllStringIndex(string(s.TextBytes), -1); len(patternPos) > 0 {
					s.Occurrences = len(patternPos)
					if s.Occurrences > 0 {
						s.readyToReplace = true
					} else {
						return
					}
					// Whether the choice is to search and replace at once, try it
					if s.DoReplace {
						s.replace()
						return
					} else {
						// Only search ... and store positions
						s.Pos = LinesInfosBuild(&s.TextBytes, patternPos)
						if s.Pos.Count > 0 {
							s.hasBeenDisplayed = false
						}
					}
				}
			}
		}
	} else {
		if s.DoReplace {
			s.replace()
		}
	}
	return
}

// LinesInfos: This structure hold some methods to get indexes
// referenced by lines from a list of indexes that was found
// using regexp' functions.
type LinesInfos struct {
	FoundLinesIdx []lineIdxInf // Indexes by lines of all found patterns
	WordsPosIdx   [][]int      // Indexes of all found patterns
	Count         int
	Eol           string

	linesIndexes [][]int // Position of lines from their offsets
	textByte     *[]byte
	onEachLine   func(idx, lineStart, lineEnd int)
}

// lineIdxInf: hold line informations (pattern matched corresponding
// to line number).
type lineIdxInf struct {
	Number   int   // Line number
	FoundIdx []int // Indexes of all found patterns in this line
}

// LinesInfosBuild: Create a structure to hold indexes and positions
// of them by line number. "patternPos" must contain the results of
// a previous call to "FindAllStringIndex".
func LinesInfosBuild(textByte *[]byte, patternPos [][]int) (li LinesInfos) {

	li.WordsPosIdx = patternPos
	li.textByte = textByte
	li.init()
	li.buildLinesIndexes()
	return
}

// initLinesIndexes: Initialize structure. In case where the
// "BuildLinesIndexes" function is not used to create the structure,
// the "textByte" and "patternPos" variables need to be be filled
// before using "Init" function as well.
func (li *LinesInfos) init() {

	// Init variables' struct
	li.Count = len(li.WordsPosIdx)
	li.FoundLinesIdx = make([]lineIdxInf, li.Count)
	li.Eol = glsg.GetTextEOL(*li.textByte)

	// Build EOL indexes
	eolRegx := regexp.MustCompile(li.Eol)
	if eolPositions := eolRegx.FindAllIndex(*li.textByte, -1); len(eolPositions) > 0 {

		// Add a fake Eol position to avoid index issue that appear
		// in the rare cases - multi-lines pattern at EOF.
		eolPositions = append(eolPositions, []int{eolPositions[len(eolPositions)-1][0] + 1, eolPositions[len(eolPositions)-1][1] + 1})

		// Define and prepare slice of line indexes
		li.linesIndexes = make([][]int, len(eolPositions))
		li.linesIndexes[0] = []int{0, eolPositions[0][0]}

		// Creating lines indexes
		for idx := 1; idx < len(eolPositions); idx++ {
			li.linesIndexes[idx] = []int{eolPositions[idx-1][1], eolPositions[idx][0]}
		}
		return
	}
	// ther is only one line ...
	li.linesIndexes = append(li.linesIndexes, []int{0, 0})
}

// getLineFromOffsets: get the line number corresponding to offsets.
// Notice: line number start at 0.
func (li *LinesInfos) getLineFromOffsets(sOfst, eOfst int) (lStart int, lEnd int) {
	for lineNb, lineIdxs := range li.linesIndexes {
		switch {
		case sOfst >= lineIdxs[0] && sOfst <= lineIdxs[1]:
			lStart = lineNb
			if eOfst <= lineIdxs[1] { // only one line
				lEnd = lineNb
				return
			}
		case eOfst >= lineIdxs[0] && eOfst <= lineIdxs[1]:
			lEnd = lineNb
			return
		}
	}
	return
}

// buildLinesIndexes: building the line indexes structure.
func (li *LinesInfos) buildLinesIndexes() {
	var lineStart, lineEnd int
	var lineIdxs, ofst []int

	for idx := 0; idx < li.Count; idx++ {

		ofst = (li.WordsPosIdx[idx])

		lineStart, lineEnd = li.getLineFromOffsets(ofst[0], ofst[1])
		li.FoundLinesIdx[idx].Number = lineStart

		lineIdxs = li.linesIndexes[lineStart]

		if li.onEachLine != nil { // TODO remove if not really usefull
			li.onEachLine(idx, lineStart, lineEnd)
		}

		startIdx := ofst[0] - lineIdxs[0]
		endIdx := ofst[1] - lineIdxs[0]

		if lineStart == lineEnd {
			// found pattern end at the same line.
			li.FoundLinesIdx[idx].FoundIdx = append(li.FoundLinesIdx[idx].FoundIdx,
				[]int{startIdx, endIdx}...)
		} else {
			// found pattern end at another line.
			li.FoundLinesIdx[idx].FoundIdx = append(
				li.FoundLinesIdx[idx].FoundIdx,
				[]int{startIdx, lineIdxs[1]}...) // fill current line to the end.

			var currLineBytesCount int
			var lineIdxsSubValues []int
			currFoundLinesIdx := idx
			restSize := ofst[1] - lineIdxs[1] - len(li.Eol)

			for lineIdxSub := lineStart + 1; lineIdxSub <= lineEnd; lineIdxSub++ {

				currFoundLinesIdx++
				lineIdxsSubValues = li.linesIndexes[lineIdxSub]
				currLineBytesCount = lineIdxsSubValues[1] - lineIdxsSubValues[0]

				li.FoundLinesIdx = append(li.FoundLinesIdx[:currFoundLinesIdx], append([]lineIdxInf{{
					Number:   lineIdxSub,          // Line number
					FoundIdx: []int{0, restSize}}, // Indexes of all found patterns in this line
				}, li.FoundLinesIdx[currFoundLinesIdx:]...)...)

				restSize = restSize - currLineBytesCount - len(li.Eol)
			}
		}
	}
}

// BuildRegexp: get regular expression from given pattern
// taking into account the parameters provided.
func BuildRegexp(search string, caseSensitive, POSIXcharClass, POSIXstrictMode,
	regExp, wildcard, useEscapeChar, wholeWord bool) (regX *regexp.Regexp, err error) {
	if !regExp {
		switch {
		case POSIXcharClass:
			search = glsscc.StringToCharacterClasses(search, caseSensitive, POSIXstrictMode)
		case wildcard:
			if !useEscapeChar {
				search = strings.Replace(search, "?", "¤_¤_¤_¤_¤_¤", -1)
				search = strings.Replace(search, "*", "_¤_¤_¤_¤_¤_", -1)
				search = regexp.QuoteMeta(search)
				search = strings.Replace(search, "¤_¤_¤_¤_¤_¤", "?", -1)
				search = strings.Replace(search, "_¤_¤_¤_¤_¤_", "*", -1)
			}
			search = strings.Replace(search, "*", `.+`, -1)
			search = strings.Replace(search, "?", `.{1}`, -1)
		case useEscapeChar:
			search = strings.Replace(search, `\`, "_¤_¤_¤_¤_¤_", -1)
			search = regexp.QuoteMeta(search)
			search = strings.Replace(search, "_¤_¤_¤_¤_¤_", `\`, -1)
		default:
			search = regexp.QuoteMeta(search)
		}
		search = `(` + search + `)`

		if wholeWord {
			search = `\b` + search + `\b`
		}
		if !caseSensitive && !POSIXcharClass {
			search = `(?i)` + search
		}
	}
	return regexp.Compile(search)
}
