// miscFunctions.go

/*
*	©2019 H.F.M. MIT license
 */

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	glfsft "github.com/hfmrow/genLib/files/fileText"
	glsg "github.com/hfmrow/genLib/strings"
	gltsbh "github.com/hfmrow/genLib/tools/bench"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gipo "github.com/hfmrow/gtk3Import/pango"
)

/*****************************/
/* Scan directory functions */
/***************************/
// isSymlinkDir: File is a symlinked directory ?
func isSymlinkDir(slRoot string, slStat os.FileInfo, followSymlinkDir bool) (slIsDir bool, err error) {
	var fName string
	if slStat.Mode()&os.ModeSymlink != 0 && followSymlinkDir {
		if fName, err = os.Readlink(filepath.Join(slRoot, slStat.Name())); err == nil {
			if slStat, err = os.Stat(fName); err == nil {
				if slStat.IsDir() {
					return true, nil
				}
			}
		}
	}
	return false, err
}

// checkUnrecoverableErr: Error come from file not exist or file permission ?
func checkUnrecoverableErr(err error) error {
	if err != nil {
		if !(os.IsPermission(err) || os.IsNotExist(err)) {
			return errors.New(fmt.Sprintf("Got error: %s\n", err.Error()))
		}
	}
	return nil
}

// ScanDirDepth: retrieve files in a specific directory and his sub-directory depending on depth argument.
// depth = -1 mean infinite, depth = 0 mean no sub-dir. optParams: showDir, followSymlinks as bool.
func ScanDirDepth(root string, depth int, optParam ...bool) (files []string, err error) {
	var showDirs, followSl, isDir bool
	switch len(optParam) {
	case 1:
		showDirs = optParam[0]
	case 2:
		showDirs = optParam[0]
		followSl = optParam[1]
	}
	var depthRecurse int
	var tmpFiles []string
	var newFi []os.FileInfo
	var fRoot *os.File
	// Starting scannig directory
	if fRoot, err = os.Open(root); err == nil {
		defer fRoot.Close()
		if newFi, err = fRoot.Readdir(-1); err == nil {
			for _, file := range newFi {
				depthRecurse = depth
				if isDir, err = isSymlinkDir(root, file, followSl); err != nil {
					return files, checkUnrecoverableErr(err)
				}
				if isDir || file.IsDir() { // It's a Dir
					if showDirs {
						files = append(files, filepath.Join(root, file.Name()))
					}
					if depth != 0 {
						depthRecurse--
						if tmpFiles, err = ScanDirDepth(filepath.Join(root, file.Name()), depthRecurse, showDirs, followSl); err == nil {
							files = append(files, tmpFiles...)
						} else {
							return files, checkUnrecoverableErr(err)
						}
					}
				} else { // Not a Dir
					files = append(files, filepath.Join(root, file.Name()))
				}
			}
		} else {
			return files, checkUnrecoverableErr(err)
		}
	} else {
		return files, checkUnrecoverableErr(err)
	}
	return files, err
}

// scanForSubDir: In case where a display refresh is requiered from an existing files/dirs list.
func scanForSubDir(inFilesList []string) (err error) {
	var stat os.FileInfo
	var filesList []string
	var isDir bool

	toDispFileList = []string{}
	for idx := len(inFilesList) - 1; idx > -1; idx-- {
		file := inFilesList[idx]
		if len(file) != 0 {

			if stat, err = os.Stat(file); /*!(os.IsPermission(err) || os.IsNotExist(err)) &&*/ err != nil {
				return err
			} else if isDir, err = isSymlinkDir(file, stat,
				mainObjects.chkFollowSymlinkDir.GetActive()); !(os.IsPermission(err) || os.IsNotExist(err)) && err != nil {
				return err
			}

			if isDir || stat.IsDir() {
				if filesList, err = ScanDirDepth(file, mainObjects.spinButtonDepth.GetValueAsInt(),
					false, mainObjects.chkFollowSymlinkDir.GetActive()); os.IsPermission(err) || os.IsNotExist(err) || err == nil {
					// false, mainObjects.chkFollowSymlinkDir.GetActive()); !(os.IsPermission(err) || os.IsNotExist(err)) && err != nil {
					toDispFileList = append(toDispFileList, filesList...)
					// try remove symdir entry
				} else {
					return err
				}
			} else {
				toDispFileList = append(toDispFileList, file)
			}
		}
	}
	return err
}

// updateTreeViewFilesDisplay: Display files to treeView
func updateTreeViewFilesDisplay() {
	var err error
	var stat os.FileInfo

	//	signalsIntercept(true)

	// TODO /* BENCH Update Display*/
	bench := new(gltsbh.Bench)
	bench.Lapse("Start")
	/* BENCH */

	if !fromDnD {
		currentInFilesList = currentInFilesList[:0]
		currentInFilesList = append(currentInFilesList, mainObjects.fileChooserBtn.GetFilename())
	} else {
		if len(currentInFilesList) == 1 {
			if stat, err = os.Stat(currentInFilesList[0]); err == nil {
				if stat.IsDir() {
					mainObjects.fileChooserBtn.SetFilename(currentInFilesList[0])
					mainObjects.switchFileChooserButton.SetActive(true)
					fromDnD = false
				}
			}
		}
	}
	if err != nil {
		gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
	} else {
		err = scanForSubDir(currentInFilesList)
		bench.Stop()
		scanTime = fmt.Sprintf("%dm %ds %dms", bench.NumTime[0].Min, bench.NumTime[0].Sec, bench.NumTime[0].Ms)
		displayFiles()
	}
	//	signalsIntercept(false)
}

func signalsIntercept(intercept bool) {
	if intercept {
		mainObjects.fileChooserBtn.HandlerBlock(fileChooserSigHandlerSelChanged)
		mainObjects.spinButtonDepth.HandlerBlock(spinButtonDepthSigHandlerChanged)
	} else {
		mainObjects.fileChooserBtn.HandlerUnblock(fileChooserSigHandlerSelChanged)
		mainObjects.spinButtonDepth.HandlerUnblock(spinButtonDepthSigHandlerChanged)
	}
}

// displayFiles: display formatted files to treeview.
func displayFiles() {
	mainOptions.UpdateOptions()
	if tvsList != nil {

		// TODO /* BENCH Display ListView*/
		bench := new(gltsbh.Bench)
		bench.Lapse("Start")
		/* BENCH */

		// Detach listStore & clean before fill
		tvsList.StoreDetach()
		tvsList.ListStore.Clear()

		formattedFilesToDisp := formatFilesToDisp(toDispFileList)
		filesCount = len(formattedFilesToDisp)

		for idx := 0; idx < filesCount; idx++ {
			tvsList.AddRow(nil, tvsList.ColValuesStringSliceToIfaceSlice(formattedFilesToDisp[idx]...))
		}

		// Attach listStore
		tvsList.StoreAttach()

		// TODO /* DEBUG */
		// fmt.Printf("From ListStore: %d, From scandir: %d\n", tvsList.CountRows(), filesCount)
		// g.WriteFile("/home/syndicate/Downloads/filesS&Replace.txt", []byte(strings.Join(toDispFileList, "\n")))
		/* DEBUG */
		bench.Stop()
		dispTime = fmt.Sprintf("%dm %ds %dms", bench.NumTime[0].Min, bench.NumTime[0].Sec, bench.NumTime[0].Ms)

		updateStatusBar()
	}
}

// formatFilesToDisp: convert filenames to be displayed into TreeView
func formatFilesToDisp(inFiles []string) (outFiles [][]string) {
	var err error
	var file os.FileInfo
	var ok bool
	var size, time string
	ExtSliceToOpt()
	for _, rowFile := range inFiles {
		if file, err = os.Stat(rowFile); err != nil {
			markedup := "<span strikethrough=\"true\" strikethrough_color=\"#FF0000\">" + strings.ReplaceAll(rowFile, "&", "&amp;") + "</span>"
			if os.IsNotExist(err) {
				outFiles = append(outFiles, []string{filepath.Base(rowFile), "Err", "", strings.Replace(sts["file-LinkNotExist"], "(s)", "", -1) + " " + markedup})
			} else if os.IsPermission(err) {
				outFiles = append(outFiles, []string{filepath.Base(rowFile), "Err", "", sts["file-perm"] + " " + markedup})
			} else {
				gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"]+"_Disp", "\n\n"+err.Error(), "", "Ok")
			}
		} else {
			// get information to be displayed.
			size = fmt.Sprintf("%s", humanize.Bytes(uint64(file.Size())))
			time = fmt.Sprintf("%s", humanize.Time(file.ModTime()))
			// Check for ext matching well.
			for _, msk := range mainOptions.ExtMask {
				if ok, err = filepath.Match(msk, file.Name()); err == nil {
					if ok {
						outFiles = append(outFiles, []string{file.Name(), size, time, rowFile})
					}
				} else {
					gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"]+"_MatchExt", "\n\n"+err.Error(), "", "Ok")
				}
			}
		}
	}
	return outFiles
}

// showResults: Get results and fill (treestore)
func showResults(mainFound *[]glfsft.Find_s) (countFiles int) {
	var err error
	var lineNbr int
	var iter *gtk.TreeIter
	var total, binary, text int

	// Detach & clear treeStore
	tvsTree.StoreDetach()
	tvsTree.TreeStore.Clear()

	for _, result := range *mainFound {
		total++
		if result.Occurrences > 0 {
			// if len(result.Positions.WordsPos) != 0 {
			text++
			/* Add parent */
			if iter, err = tvsTree.AddRow(nil, tvsTree.ColValuesStringSliceToIfaceSlice(result.FileName)); err == nil {

				// Make markup
				pm := gipo.PangoMarkup{}
				pm.Init(string(result.TextBytes)) // Text
				pm.AddPosition(result.Positions.WordsPos...)
				pm.AddTypes([][]string{{"bgc", pm.Colors.Lightgreen}}...)
				text := pm.MarkupAtPos()

				/* Add Childs */
				linesStr := strings.Split(text, result.LineEnd)
				for idxLine := len(result.Positions.Line) - 1; idxLine >= 0; idxLine-- {
					lineNbr = result.Positions.Line[idxLine]
					pm.Init(fmt.Sprintf(" %v ", lineNbr+1)) // Line number
					pm.AddTypes([][]string{{"bgc", pm.Colors.Lightgray}}...)
					lineNbStr := pm.Markup()

					tvsTree.AddRow(iter, tvsTree.ColValuesStringSliceToIfaceSlice(lineNbStr+"\t"+string(linesStr[lineNbr])))
				}
				countFiles++

			}
		} else if mainObjects.findWinChkDispForbFiles.GetActive() && result.NotTextFile {
			binary++
			// Markup bad file: File size too short, Binary files warning ...
			pm := gipo.PangoMarkup{}
			pm.Init(string(result.FileName)) // Filename
			pm.AddPosition(result.Positions.WordsPos...)
			pm.AddTypes([][]string{{"stc", pm.Colors.Red}}...)
			filename := pm.Markup()

			/* Add parent */
			if iter, err = tvsTree.AddRow(nil, tvsTree.ColValuesStringSliceToIfaceSlice(filename)); err == nil {
				pm.Init(string(result.TextBytes)) // text with reason why it's not displayed
				pm.AddTypes([][]string{{"fgc", pm.Colors.Red}}...)
				desciption := pm.Markup()

				/* Add Child */
				tvsTree.AddRow(iter, tvsTree.ColValuesStringSliceToIfaceSlice(desciption))
			}
		}
	}
	// Attach treeStore
	tvsTree.StoreAttach()
	return countFiles
}

// Show text edit window with colored text (preview)
func showTextWin(text string, filename string, line int) {
	var occurrences int
	var buff *gtk.TextBuffer
	var err error

	mainObjects.findWin.SetModal(false)
	mainObjects.textWin.SetKeepAbove(false)
	mainObjects.textWinTextview.SetEditable(true)
	x, y := mainObjects.mainWin.GetPosition()
	mainObjects.textWin.Move(x+(mainOptions.CascadeDepth*2), y+(mainOptions.CascadeDepth*2))
	mainObjects.textWin.Resize(mainObjects.mainWin.GetSize())

	// Get buffer from textview
	if buff, err = mainObjects.textWinTextview.GetBuffer(); err == nil {
		// Set text
		if occurrences, err = highlightText(buff, text, line); err == nil {
			textWinTitle.MainTitle = truncatePath(filename, mainOptions.FilePathLength)
			textWinTitle.Update([]string{fmt.Sprintf("%s %d", sts["totalOccurrences"], occurrences)})
			mainObjects.textWin.Show()
		}
	}
	if err != nil {
		gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
	}
}

// Find ...
func Find(entrySearchText, entryReplaceText string, removeEmptyResult bool) (mainFound []glfsft.Find_s, occurrences int, err error) {
	if tvsList.Selection.CountSelectedRows() == 0 {
		return mainFound, occurrences, errors.New("\n\n" + sts["noFileSel"] + "\n")
	}
	gimc.Notify("Information", "Please wait while processing ...")

	// Retrieving filenames
	var value *glib.Value
	var str string
	var iters []*gtk.TreeIter
	var treeviewSelectedRows []string

	if iters, err = tvsList.GetSelectedIters(); err == nil {
		for _, iter := range iters {
			if value, err = tvsList.ListStore.GetValue(iter, 3); err == nil { // Field 3: get full path
				if str, err = value.GetString(); err == nil {
					treeviewSelectedRows = append(treeviewSelectedRows, str)
				} else {
					break
				}
			}
		}
	}

	if err != nil {
		return mainFound, occurrences, errors.New("\n" + sts["alert"] + "\n\n" + err.Error())
	}

	// g.TmpCount = 0

	mainFound, occurrences, err = glfsft.SearchAndReplaceInMultipleFiles(
		treeviewSelectedRows,
		entrySearchText,
		entryReplaceText,
		mainOptions.LineEndThreshold,   //thresholdLineEnd
		mainOptions.OverCharsThreshold, // thresholdOverChars
		mainOptions.FileSizeLimit,      // Size limit
		mainObjects.chkCaseSensitive.GetActive(),
		mainObjects.chkCharacterClass.GetActive(),
		mainObjects.chkCharacterClassStrictMode.GetActive(),
		mainObjects.chkRegex.GetActive(),
		mainObjects.chkWildcard.GetActive(),
		mainObjects.chkUseEscapeChar.GetActive(),
		mainObjects.chkWoleWord.GetActive(),
		applyChanges, // DoReplace
		applyChanges, // DoSave
		mainOptions.MakeBackup,
		acceptBinary,
		removeEmptyResult)

	// Errors handling ...
	if os.IsNotExist(err) {
		return mainFound, occurrences, errors.New("\n\n" + sts["removed"] + "\n" + err.Error())
	}
	if err != nil {
		return mainFound, occurrences, err
	}
	return mainFound, occurrences, nil // all is ok, lets return ...
}

// onTheFlyReplace: Simply replace
func onTheFlyReplace(inTextBytes []byte, outText *[]byte, entrySearchText string) (err error) {
	if len(inTextBytes) != 0 {
		entryReplaceText := getEntryText(mainObjects.entryReplace)
		mainFound := glfsft.Find_s{}
		mainFound.FileName = "From clipboard"
		mainFound.TextBytes = inTextBytes
		mainFound.ReplaceWith = entryReplaceText
		mainFound.ToSearch = entrySearchText
		mainFound.DoReplace = true
		mainFound.CaseSensitive = mainObjects.chkCaseSensitive.GetActive()
		mainFound.LineEnd = getTextEOL(mainFound.TextBytes)
		mainFound.POSIXcharClass = mainObjects.chkCharacterClass.GetActive()
		mainFound.POSIXstrictMode = mainObjects.chkCharacterClassStrictMode.GetActive()
		mainFound.Regex = mainObjects.chkRegex.GetActive()
		mainFound.Wildcard = mainObjects.chkWildcard.GetActive()
		mainFound.UseEscapeChar = mainObjects.chkUseEscapeChar.GetActive()
		mainFound.WholeWord = mainObjects.chkWoleWord.GetActive()

		gimc.Notify("Information", "Please wait while processing ...")
		err = mainFound.SearchAndReplace()
		*outText = mainFound.TextBytes
	}
	return err
}

/*
* Highlight preview with pango markup
 */
var pangoEscapeChar = [][]string{{"<", "&lt;", string([]byte{0x15})}, {"&", "&amp;", string([]byte{0x16})}}

// Prepare string with special characters to be marked ("<", "&")
func PreparePango(inString string) string {
	inString = strings.Replace(inString, pangoEscapeChar[1][0], pangoEscapeChar[1][2], -1)
	return strings.Replace(inString, pangoEscapeChar[0][0], pangoEscapeChar[0][2], -1)
}

// Escape special characters after marking ("<", "&")
func FinalizePango(inString string) string {
	inString = strings.Replace(inString, pangoEscapeChar[1][2], pangoEscapeChar[1][1], -1)
	return strings.Replace(inString, pangoEscapeChar[0][2], pangoEscapeChar[0][1], -1)
}

// Highlight
func highlightText(buffer *gtk.TextBuffer, txtStr string, line int) (occurrences int, err error) {
	if len(txtStr) != 0 {
		pangoBfr := []byte(`<span background="#87FF87">`) // foreground Lightgreen color
		pangoAfr := []byte(`</span>`)
		byteTxt := []byte(PreparePango(txtStr))

		mainFound := glfsft.Find_s{}
		mainFound.FileName = "From clipboard"
		mainFound.TextBytes = byteTxt
		mainFound.ReplaceWith = getEntryText(mainObjects.entryReplace)
		mainFound.ToSearch = getEntryText(mainObjects.entrySearch)
		mainFound.DoReplace = false
		mainFound.CaseSensitive = mainObjects.chkCaseSensitive.GetActive()
		mainFound.LineEnd = glsg.GetTextEOL(mainFound.TextBytes)
		mainFound.POSIXcharClass = mainObjects.chkCharacterClass.GetActive()
		mainFound.POSIXstrictMode = mainObjects.chkCharacterClassStrictMode.GetActive()
		mainFound.Regex = mainObjects.chkRegex.GetActive()
		mainFound.Wildcard = mainObjects.chkWildcard.GetActive()
		mainFound.UseEscapeChar = mainObjects.chkUseEscapeChar.GetActive()
		mainFound.WholeWord = mainObjects.chkWoleWord.GetActive()

		gimc.Notify("Information", "Please wait while processing ...")
		if err = mainFound.SearchAndReplace(); err == nil {
			occurrences = mainFound.Occurrences
			for posIdx := len(mainFound.Positions.WordsPos) - 1; posIdx >= 0; posIdx-- {
				pos := mainFound.Positions.WordsPos[posIdx]
				byteTxt = append(byteTxt[:pos[1]], append(pangoAfr, byteTxt[pos[1]:]...)...)
				byteTxt = append(byteTxt[:pos[0]], append(pangoBfr, byteTxt[pos[0]:]...)...)
			}
			buffer.Delete(buffer.GetStartIter(), buffer.GetEndIter())
			buffer.InsertMarkup(buffer.GetStartIter(), FinalizePango(string(byteTxt)))

			if !mainObjects.textWinChkWrap.GetActive() {
				mainObjects.textWinTextview.SetProperty("left-margin", 2)

				mainOptions.textViewRowNumber.UpdateTextViewNumbers()
			}
			// Scroll textview at desired line.
			mainOptions.textViewRowNumber.ScrollToLine(line)
		}
	}
	return occurrences, err
}

/*
* StatusBar function ...
 */
func updateStatusBar() {
	wordFile := sts["sbFile"]
	wordSel := sts["sbFileSel"]

	if tvsList.CountRows() > 1 {
		wordFile = sts["sbFiles"]
	}
	if tvsList.Selection.CountSelectedRows() > 1 {
		wordSel = sts["sbFilesSel"]
	}
	statusbar.Prefix[0] = wordFile
	statusbar.Prefix[1] = wordSel

	statusbar.CleanAll()
	statusbar.Set(fmt.Sprint(tvsList.CountRows()), 0)
	statusbar.Set(fmt.Sprint(tvsList.Selection.CountSelectedRows()), 1)
	statusbar.Set(fmt.Sprint(scanTime), 2)
	statusbar.Set(fmt.Sprint(searchTime), 3)
	statusbar.Set(fmt.Sprint(dispTime), 4)
}

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

/*
* Clipboard handling ...
 */
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

// getEntryText: retrieve value of an entry control.
func getEntryText(e *gtk.Entry) (outTxt string) {
	var err error

	outTxt, err = e.GetText()
	Check(err, "getEntryText")
	return outTxt
}

/*
* Text functions
 */
// Convert Entry to list of extensions.
func ExtSliceToOpt() {
	mainOptions.ExtMask = []string{}

	tmpSliceStrings := strings.Split(getEntryText(mainObjects.entryExtMask), mainOptions.ExtSep)
	for _, str := range tmpSliceStrings {
		str = strings.TrimSpace(str)
		if len(str) != 0 {
			mainOptions.ExtMask = append(mainOptions.ExtMask, str)
		}
	}
	mainObjects.entryExtMask.SetText(strings.Join(mainOptions.ExtMask, mainOptions.ExtSep+" "))
}

func OptToExtSlice() {
	mainObjects.entryExtMask.SetText(strings.Join(mainOptions.ExtMask, mainOptions.ExtSep+" "))
}

// ReducePath: Reduce path length by preserving count element from the end
func truncatePath(fullpath string, count ...int) (reduced string) {
	elemCnt := 2
	if len(count) != 0 {
		elemCnt = count[0]
	}
	splited := strings.Split(fullpath, string(os.PathSeparator))
	if len(splited) > elemCnt+1 {
		return "..." + string(os.PathSeparator) + filepath.Join(splited[len(splited)-elemCnt:]...)
	}
	return fullpath
}

// getTextEOL: Get EOL from text bytes (CR, LF, CRLF)
func getTextEOL(inTextBytes []byte) (outString string) {
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
