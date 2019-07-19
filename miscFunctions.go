// miscFunctions.go

/*
*	©2019 H.F.M. MIT license
 */

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	g "github.com/hfmrow/sAndReplace/genLib"
	gi "github.com/hfmrow/sAndReplace/gtk3Import"
)

// Check: Display error messages in HR version with onClickJump enabled in
// my favourite Golang IDE editor. Return true if error exist.
func Check(err error, message ...string) (state bool) {
	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`) //	to match 2 or more whitespace symbols inside a string
	var msgs string
	if err != nil {
		state = true
		if len(message) != 0 { // Make string with messages if exists
			for _, mess := range message {
				msgs += `[` + mess + `]`
			}
		}
		pc, file, line, ok := runtime.Caller(1) //	(pc uintptr, file string, line int, ok bool)
		if ok == false {                        // Remove "== false" if needed
			fName := runtime.FuncForPC(pc).Name()
			fmt.Printf("[%s][%s][File: %s][Func: %s][Line: %d]\n", msgs, err.Error(), file, fName, line)
		} else {
			stack := strings.Split(fmt.Sprintf("%s", debug.Stack()), "\n")
			for idx := 5; idx < len(stack)-1; idx = idx + 2 {
				//	To match 2 or more whitespace leading/ending/inside a string (include \t, \n)
				mess1 := strings.Join(strings.Fields(stack[idx]), " ")
				mess2 := strings.TrimSpace(remInside.ReplaceAllString(stack[idx+1], " "))
				fmt.Printf("%s[%s][%s]\n", msgs, err.Error(), strings.Join([]string{mess1, mess2}, "]["))
			}
		}
	}
	return state
}

// FindDir retrieve file in a specific directory with more options.
func FindDir(dir string, masks []string, returnedStrSlice *[][]string, scanSub, showDir, followSymlinkDir bool) (err error) {
	var ok bool
	var fName, time, size string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("%s\nDir-Name: %s\n", err.Error(), dir))
	}
	for _, file := range files {
		fName = filepath.Join(dir, file.Name())
		if followSymlinkDir { // Check for symlink ..
			file, err = os.Lstat(fName)
			if err != nil {
				return errors.New(fmt.Sprintf("%s\nFilename: %s\n", err.Error(), fName))
			}
			if file.Mode()&os.ModeSymlink != 0 { // Is a symlink ?
				fName, err := os.Readlink(fName) // Then read it...
				if err != nil {
					return errors.New(fmt.Sprintf("%s\nFilename: %s\n", err.Error(), fName))
				}
				file, err = os.Stat(fName) // Get symlink infos.
				if err != nil {
					return errors.New(fmt.Sprintf("%s\nFilename: %s\n", err.Error(), fName))
				}
				fName = filepath.Join(dir, file.Name())
			}
		}
		// Recursive play if it's a directory
		if file.IsDir() && scanSub {
			tmpFileList := new([][]string)
			err = FindDir(fName, masks, tmpFileList, scanSub, showDir, followSymlinkDir)
			if err != nil {
				return errors.New(fmt.Sprintf("%s\nFilename: %s\n", err.Error(), fName))
			}
			*returnedStrSlice = append(*returnedStrSlice, *tmpFileList...)
		}
		// get information to be displayed.
		size = fmt.Sprintf("%s", humanize.Bytes(uint64(file.Size())))
		time = fmt.Sprintf("%s.", humanize.Time(file.ModTime()))
		// Check for ext matching well.
		for _, msk := range masks {
			ok, err = filepath.Match(msk, file.Name())
			if err != nil {
				return err
			}
			if ok {
				break
			}
		}
		if ok {
			if showDir { // Limit display directories if requested
				*returnedStrSlice = append(*returnedStrSlice, []string{file.Name(), size, time, fName})
			} else {
				_, err = ioutil.ReadDir(fName)
				if err != nil {
					*returnedStrSlice = append(*returnedStrSlice, []string{file.Name(), size, time, fName})
				}
			}
		}
	}
	return nil
}

// BuildExtSlice
func BuildExtSlice() {
	var extSep = ";"
	mainOptions.ExtMask = []string{}

	tmpSliceStrings := strings.Split(getEntryText(mainObjects.entryExtMask), extSep)
	for _, str := range tmpSliceStrings {
		str = strings.TrimSpace(str)
		if len(str) != 0 {
			mainOptions.ExtMask = append(mainOptions.ExtMask, str)
		}
	}
	mainObjects.entryExtMask.SetText(strings.Join(mainOptions.ExtMask, extSep+" "))
}

// scanFilesAndDisp:
func scanFilesAndDisp() {
	var err error
	mainOptions.UpdateOptions()
	if !cmdLineArg {
		// Check if directory exist ...
		if _, err = os.Stat(mainOptions.Directory); os.IsNotExist(err) {
			gi.DlgMessage(mainObjects.mainWin, "error", mainOptions.TxtAlert, "\n"+sts["dir-rem"]+"\n\n"+err.Error(), "", "Ok")
			mainOptions.Directory = filepath.Dir(os.Args[0])
		}

		// Set control value to be sure is always displayed (case of cmdline input.)
		mainObjects.fileChooserBtn.SetCurrentFolder(mainOptions.Directory)
		filesList := new([][]string)

		err = FindDir(mainOptions.Directory, mainOptions.ExtMask, filesList,
			mainObjects.chkSubDir.ToggleButton.GetActive(),
			false,
			mainObjects.chkFollowSymlinkDir.ToggleButton.GetActive())
		g.Check(err, "fileChooserBtnClicked", "FindDir")

		if mainObjects.listStore != nil {
			filesCount = len(*filesList)
			// Clean before fill
			mainObjects.listStore.Clear()
			mainOptions.currentInFilesList = mainOptions.currentInFilesList[:0]
			for idx := 0; idx < filesCount; idx++ {
				mainOptions.currentInFilesList = append(mainOptions.currentInFilesList, (*filesList)[idx][3])
				gi.ListStoreAddRow(mainObjects.listStore, (*filesList)[idx]) // dereferencing pointer to index slice ...
			}
			updateStatusBar()
		}
	}
}

// getFilesSelection: Format selected files to be displayed into treeview
func getFilesSelection(inStrings []string) (err error) {
	var tmpErr string
	var file os.FileInfo
	var filesList [][]string
	var errorStop, ok bool
	var errorCount int

	if mainObjects.listStore != nil {
		for _, fullName := range inStrings {

			// Check for ext matching well.
			if len(mainOptions.ExtMask) > 0 {
				for _, msk := range mainOptions.ExtMask {
					ok, err = filepath.Match(msk, filepath.Base(fullName))
					if err != nil {
						return err
					}
					if ok {
						break
					}
				}
			} else {
				ok = true
			}
			if ok {
				file, err = os.Stat(fullName)
				if err == nil {
					// get informations to be displayed.
					filesList = append(filesList, []string{file.Name(),
						fmt.Sprintf("%s", humanize.Bytes(uint64(file.Size()))),
						fmt.Sprintf("%s.", humanize.Time(file.ModTime())),
						fullName})
				} else {
					if errorCount < 5 {
						tmpErr += err.Error() + "\n"
						errorCount++
					} else if !errorStop {
						tmpErr += "\nAnd more ..."
						errorStop = true
					}
				}
			}
			ok = false
		}
		filesCount = len(filesList)
		// Clean before fill
		mainObjects.listStore.Clear()
		if filesCount > 0 {
			for idx := 0; idx < filesCount; idx++ {
				gi.ListStoreAddRow(mainObjects.listStore, (filesList[idx]))
			}
		}
		updateStatusBar()
	}
	if len(tmpErr) != 0 {
		err = errors.New("\n" + tmpErr)
	}
	return err
}

// getEntryText: retrieve value of an entry control.
func getEntryText(e *gtk.Entry) (outTxt string) {
	outTxt, err = e.GetText()
	Check(err, "getEntryText")
	return outTxt
}

// Clipboard handling ...
func clipboardInit() {
	var err error
	mainObjects.clipboard, err = gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
	Check(err, "clipboardInit")
}

func clipboardCopyFromTextWin() {
	txtBuffer, err := mainObjects.textWinTextview.GetBuffer()
	Check(err, "clipboardCopyFromTextWin")
	iterStart := txtBuffer.GetStartIter()
	iterEnd := txtBuffer.GetEndIter()
	text, err := txtBuffer.GetText(iterStart, iterEnd, true)
	Check(err, "clipboardCopyFromTextWin")

	mainObjects.clipboard.SetText(text)
}

func clipboardPastToTextWin() {
	txtBuffer, err := mainObjects.textWinTextview.GetBuffer()
	Check(err, "clipboardPastToTextWin")
	txtBuffer.SetText(clipboardGet())
}

func clipboardGet() (clipboardContent string) {
	clipboardContent, err := mainObjects.clipboard.WaitForText()
	Check(err, "clipboardGet")
	return clipboardContent
}

func clipboardSet(clipboardContent string) {
	mainObjects.clipboard.SetText(clipboardContent)
}

// Get results and fill treestore
func showResults(mainFound *[]g.Find_s) {
	var lineNbr int
	mainObjects.treeStore.Clear()
	for _, result := range *mainFound {
		if len(result.Positions.WordsPos) != 0 {
			iter := gi.TreeStoreAddRow(mainObjects.treeStore, result.FileName)
			// Make markup
			pm := gi.PangoMarkup{}
			pm.Init(string(result.TextBytes)) // Text
			pm.AddPosition(result.Positions.WordsPos...)
			pm.AddTypes([][]string{{"bgc", pm.Colors.Lightgreen}}...)
			text := pm.MarkupAtPos()
			linesStr := strings.Split(text, result.LineEnd)
			for idxLine := len(result.Positions.Line) - 1; idxLine >= 0; idxLine-- {
				lineNbr = result.Positions.Line[idxLine]
				pm.Init(fmt.Sprintf(" %v ", lineNbr+1)) // Line number
				pm.AddTypes([][]string{{"bgc", pm.Colors.Lightgray}}...)
				lineNbStr := pm.Markup()
				gi.TreeStoreAddSubRow(mainObjects.treeStore,
					iter,
					lineNbStr+"\t"+linesStr[lineNbr])
			}
		} else {
			// Markup bad file: File size too short, Binary files warning ...
			pm := gi.PangoMarkup{}
			pm.Init(result.FileName) // Filename
			pm.AddPosition(result.Positions.WordsPos...)
			pm.AddTypes([][]string{{"stc", pm.Colors.Red}}...)
			filename := pm.Markup()
			iter := gi.TreeStoreAddRow(mainObjects.treeStore, filename)
			pm.Init(string(result.TextBytes)) // text with reason why is not displayed
			pm.AddTypes([][]string{{"fgc", pm.Colors.Red}}...)
			desciption := pm.Markup()
			gi.TreeStoreAddSubRow(mainObjects.treeStore, iter, desciption)
		}
	}
}

// Show text edit window with colored text
func showTextWin(text string) {
	if !mainObjects.textWin.GetVisible() {
		mainObjects.textWinTextview.SetEditable(true)
		margin := 4

		mainObjects.textWinTextview.SetLeftMargin(margin)
		mainObjects.textWinTextview.SetRightMargin(margin)
		mainObjects.textWinTextview.SetMarginTop(margin)
		mainObjects.textWinTextview.SetMarginBottom(margin)
		// Get buffer from textview
		buff, err := mainObjects.textWinTextview.GetBuffer()
		Check(err, "showTextWin")
		// Set text
		buff.SetText(text)
		mainObjects.textWin.Show()
	}
}

// StatusBar function ...
func updateStatusBar() {
	wordFile := " Files"
	wordSel := " Files selected"
	if filesCount < 2 {
		wordFile = " File"
	}
	if filesSelected < 2 {
		wordSel = " File selected"
	}
	displayStatusBar(fmt.Sprint(filesCount, wordFile), fmt.Sprint(filesSelected, wordSel))
}

// StatusBar function ...
func displayStatusBar(str ...string) {
	var outText []string
	if len(str) != 0 {
		for _, toDisp := range str {
			outText = append(outText, toDisp)
		}
		contextId1 := mainObjects.statusbar.GetContextId("part1")
		mainObjects.statusbar.Push(contextId1, strings.Join(outText, " | "))
	}
}

// Find ...
func Find(entrySearchText, entryReplaceText string, removeEmptyResult bool) (mainFound []g.Find_s, err error) {
	if len(treeviewSelectedRows) != 0 {
		gi.Notify("Information", "Please wait while processing ...")
		mainFound, err := g.SearchAndReplaceInMultipleFiles(
			treeviewSelectedRows,
			entrySearchText,
			entryReplaceText,
			mainObjects.chkCaseSensitive.GetActive(),
			mainObjects.chkCharacterClass.GetActive(),
			mainObjects.chkCharacterClassStrictMode.GetActive(),
			mainObjects.chkRegex.GetActive(),
			mainObjects.chkWildcard.GetActive(),
			mainObjects.chkUseEscapeChar.GetActive(),
			mainObjects.chkWoleWord.GetActive(),
			applyChanges,
			applyChanges,
			mainOptions.MakeBackup,
			acceptBinary,
			removeEmptyResult)
		//		gi.Notify("Information", "Processing Done ...")

		// Errors handling ...
		if err != nil {
			var display string
			if os.IsNotExist(err) {
				display = mainOptions.TxtFileRemovedBeforeProcess
			} else {
				display = mainOptions.TextRegexError
			}
			return mainFound, errors.New("\n\n" + display + "\n" + g.FormatText(err.Error(), 69, true))
		} else {
			return mainFound, nil // all is ok, lets return ...
		}
	} else {
		return mainFound, errors.New("\n\n" + mainOptions.TxtNoFileToSearch + "\n")
	}
	return mainFound, err
}

// Simply replace
func onTheFlyReplace(inTextBytes []byte, outText *[]byte, entrySearchText string) (err error) {
	if len(inTextBytes) != 0 {
		entryReplaceText := getEntryText(mainObjects.entryReplace1)
		mainFound := g.Find_s{}
		mainFound.FileName = "From clipboard"
		mainFound.TextBytes = inTextBytes
		mainFound.ReplaceWith = entryReplaceText
		mainFound.ToSearch = entrySearchText
		mainFound.DoReplace = true
		mainFound.CaseSensitive = mainObjects.chkCaseSensitive.GetActive()
		mainFound.LineEnd = g.GetTextEOL(mainFound.TextBytes)
		mainFound.POSIXcharClass = mainObjects.chkCharacterClass.GetActive()
		mainFound.POSIXstrictMode = mainObjects.chkCharacterClassStrictMode.GetActive()
		mainFound.Regex = mainObjects.chkRegex.GetActive()
		mainFound.Wildcard = mainObjects.chkWildcard.GetActive()
		mainFound.UseEscapeChar = mainObjects.chkUseEscapeChar.GetActive()
		mainFound.WholeWord = mainObjects.chkWoleWord.GetActive()

		gi.Notify("Information", "Please wait while processing ...")
		err = mainFound.SearchAndReplace()
		*outText = mainFound.TextBytes
	}
	return err
}

// initDropSets: configure controls to receive dndcontent.
func initDropSets() {
	var targets []gtk.TargetEntry // Build dnd context
	te, err := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_OTHER_APP, 0)
	if err != nil {
		log.Fatal(err)
	}
	targets = append(targets, *te)

	mainObjects.treeviewFiles.DragDestSet(
		gtk.DEST_DEFAULT_ALL,
		targets,
		gdk.ACTION_COPY)
}
