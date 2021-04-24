// miscFunctions.go

/*
	Source file auto-generated on Sat, 24 Apr 2021 04:33:11 using Gotk3 Objects Handler v1.7.8
	©2018-21 hfmrow https://hfmrow.github.io

	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 hfmrow - Search And Replace v1.10 github.com/hfmrow/search-and-replace

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/glib"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	humanize "github.com/dustin/go-humanize"

	glfsft "github.com/hfmrow/genLib/files/fileText"

	gipo "github.com/hfmrow/gtk3Import/pango"
)

// displayProgressBar: Hide or show some controls during a long search process.
func displayProgressBar(toggle bool) {

	obj.mainBox.SetSensitive(!toggle)
	obj.MainButtonOptions.SetSensitive(!toggle)
}

// updateTreeViewFilesDisplay: Display files to treeView
func updateTreeViewFilesDisplay() {
	var err error
	var stat os.FileInfo

	bench := BenchNew(false)
	bench.Lapse()

	if !fromDnD {
		currentInFilesList = currentInFilesList[:0]
		currentInFilesList = append(currentInFilesList, obj.fileChooserBtn.GetFilename())
	} else {
		if len(currentInFilesList) == 1 {
			if stat, err = os.Stat(currentInFilesList[0]); err == nil {
				if stat.IsDir() {
					obj.fileChooserBtn.SetFilename(currentInFilesList[0])
					fromDnD = false
				}
			}
		}
	}
	if err != nil {
		DialogMessage(obj.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
	} else {
		err = scanForSubDir(currentInFilesList)
		bench.Stop()
		scanTime = bench.Lapses[0].StringShort
		displayFiles()
	}
}

// displayFiles: display formatted files to treeview.
func displayFiles() {
	opt.UpdateOptions()

	if tvsList != nil {
		bench := BenchNew(false)
		bench.Lapse()

		// Detach listStore & clean before fill it
		tvsList.StoreDetach()
		defer tvsList.StoreAttach()
		tvsList.Clear()

		var ok bool
		var err error

		for _, rowFile := range toDispFileList {
			for _, msk := range opt.ExtMask { // Check for ext matching well.
				if ok, err = filepath.Match(msk, rowFile.FileInfo.Name()); err == nil {
					if ok {
						tvsList.AddRow(nil,
							html.EscapeString(rowFile.FileInfo.Name()),
							HumanReadableSize(rowFile.FileInfo.Size()),
							// TODO make a human readable Time function
							humanize.Time(rowFile.FileInfo.ModTime()),
							html.EscapeString(rowFile.Filename),
							rowFile.Filename, // Unmodified filename
							rowFile.FileInfo.Size(),
							rowFile.FileInfo.ModTime().Unix())
					}
				} else {
					DlgErr("_MatchExt", err)
				}
			}
		}

		// Attach listStore

		bench.Stop()

		dispTime = bench.Lapses[0].StringShort

		updateStatusBar()
	}
}

// showResults: Get results and fill (treestore)
func showResults(mainFound *[]glfsft.SearchAndReplaceFiles) (countFiles int) {
	var err error
	var lineNbr, prevLineNbr int
	var iter *gtk.TreeIter
	var total, binary, text int

	pc := PangoColorNew()

	// Detach & clear treeStore
	tvsTree.StoreDetach()
	defer tvsTree.StoreAttach()
	tvsTree.Clear()

	for fileIdx, result := range *mainFound {
		prevLineNbr = -1
		total++
		if result.Occurrences > 0 {
			text++
			/* Add parent */
			if iter, err = tvsTree.AddRow(nil, true, result.FileName, fileIdx, -1); err == nil {

				// Make markup
				pm := gipo.PangoMarkup{}
				pm.Init(string(result.SearchAndRepl.TextBytes)) // Text
				pm.AddPosition(result.SearchAndRepl.Pos.WordsPosIdx...)
				pm.AddTypes([][]string{{"bgc", "#AAFFAA"}}...) // Light green
				text := pm.MarkupAtPos()

				/* Add Childs */
				linesStr := strings.Split(text, result.SearchAndRepl.Pos.Eol)
				for lineIdx := 0; lineIdx < len(result.SearchAndRepl.Pos.FoundLinesIdx); lineIdx++ {

					lineNbr = result.SearchAndRepl.Pos.FoundLinesIdx[lineIdx].Number
					if prevLineNbr != lineNbr { // Avoid displaying duplicated rows
						pm.Init(fmt.Sprintf(" %v ", lineNbr+1))        // Line number
						pm.AddTypes([][]string{{"bgc", "#E4DDDD"}}...) // Light grey
						lineNbStr := pm.Markup()

						tvsTree.AddRow(iter, true, lineNbStr+"\t"+string(linesStr[lineNbr]), fileIdx, lineIdx)
					}
					prevLineNbr = lineNbr
				}
				countFiles++
			}
		} else if obj.findWinChkDispForbFiles.GetActive() && result.NotTextFile {
			binary++
			// Markup bad file: File size too short or Binary files warning ...
			pm := gipo.PangoMarkup{}
			pm.Init(string(result.FileName)) // Filename
			pm.AddPosition(result.SearchAndRepl.Pos.WordsPosIdx...)
			pm.AddTypes([][]string{{"stc", pc.Red}}...)
			filename := pm.Markup()

			/* Add parent */
			if iter, err = tvsTree.AddRow(nil, false, filename, fileIdx, -1); err == nil {
				pm.Init(string(result.SearchAndRepl.TextBytes)) // text with reason why it's not displayed
				pm.AddTypes([][]string{{"fgc", pc.Red}}...)
				desciption := pm.Markup()

				/* Add Child */
				tvsTree.AddRow(iter, false, desciption, fileIdx, -1)
			}
		}
	}
	glib.IdleAdd(func() {
		tvsTree.ExpandAll(!opt.ExpandAll)
	})

	return countFiles
}

// BringToFront: Set window position to be over all others windows
// without staying on top whether another window come to be selected.
func BringToFront(win *gtk.Window) {
	win.Deiconify()
	win.ShowAll()
	win.GrabFocus()
}

// Show text edit window with colored text (preview)
func showTextWin(text string, filename string) {
	var occurrences int
	var err error

	if !alreadyPlacedPrevWin {
		obj.findWin.SetModal(false)
		obj.textWin.SetKeepAbove(false)
		x, y := obj.mainWin.GetPosition()
		obj.textWin.Move(x+(opt.CascadeDepth*2), y+(opt.CascadeDepth*2))
		obj.textWin.Resize(obj.mainWin.GetSize())
		alreadyPlacedPrevWin = true
	}
	// Set text
	if occurrences, err = highlightText(text, obj.textWinChkShowModifications.GetActive()); err == nil {
		textWinTitle.MainTitle = TruncatePath(filename, opt.FilePathLength)
		textWinTitle.Update([]string{fmt.Sprintf("%s %d", sts["totalOccurrences"], occurrences)})
		svs.BringToFront()
	}

	DlgErr("showTextWin:highlightText", err)
}

// highlightText: Highlight (found) and syntax highlight operations.
func highlightText(txtStr string, doReplace bool) (occurrences int, err error) {

	// Set colors for GtkSourceView
	svs.TxtBgCol = opt.TxtBgCol
	svs.TxtFgCol = opt.TxtFgCol

	// store current text
	if currentText != txtStr || (!doReplace && currentTextChanged) {
		currentText = txtStr
		svs.SetText(txtStr)
		currentTextChanged = false
	}

	if len(currentText) != 0 {

		tmpTxtStr, _ := onTheFlySearch([]byte(txtStr), doReplace)

		svs.SetText(string(tmpTxtStr))

		if doReplace {
			currentTextChanged = true

			tmpTextSearch := svs.TextSearch
			svs.TextSearch = GetEntryText(obj.entryReplace)
			svs.Search(svs.Buffer.GetStartIter(), false)
			svs.TextSearch = tmpTextSearch

		} else {
			// SourceView search part: used to highlight found or replaced patterns
			svs.TextSearch = UnEscapedStr(GetEntryText(obj.entrySearch))
			svs.UseRegexp = obj.chkRegex.GetActive()
			svs.WordBoundaries = obj.chkWholeWord.GetActive()
			svs.CaseSensitive = obj.chkCaseSensitive.GetActive()
			svs.Search(svs.Buffer.GetStartIter(), false)
		}
		occurrences = fileFoundSingle.Occurrences
		// textWinTitle.Update([]string{fmt.Sprintf("%s %d", sts["totalOccurrences"], fileFoundSingle.Occurrences)})

		if currentLine > -1 {
			svs.RunAfterEvents(func() {
				svs.ScrollToLine(currentLine)
			})
		}
	}

	// // Set colors for GtkSourceView
	// svs.TxtBgCol = opt.TxtBgCol
	// svs.TxtFgCol = opt.TxtFgCol

	// // store current text
	// if currentText != txtStr || (!doReplace && currentTextChanged) {
	// 	currentText = txtStr
	// 	svs.SetText(txtStr)
	// 	currentTextChanged = false
	// }

	// if len(currentText) != 0 {

	// 	svs.TextSearch = GetEntryText(obj.entrySearch)
	// 	svs.UseRegexp = obj.chkRegex.GetActive()
	// 	svs.WordBoundaries = obj.chkWholeWord.GetActive()
	// 	svs.CaseSensitive = obj.chkCaseSensitive.GetActive()

	// 	if !doReplace {
	// 		svs.Search(svs.Buffer.GetStartIter(), false)
	// 	} else {
	// 		svs.TextReplace = GetEntryText(obj.entryReplace)
	// 		occurrences, err = svs.SearchCtx.ReplaceAll(svs.TextReplace)
	// 		currentTextChanged = true
	// 		tmpTextSearch := svs.TextSearch
	// 		svs.TextSearch = svs.TextReplace
	// 		svs.Search(svs.Buffer.GetStartIter(), false)
	// 		svs.TextSearch = tmpTextSearch
	// 	}

	// 	svs.GetOccurences(func(occ int) {
	// 		textWinTitle.Update([]string{fmt.Sprintf("%s %d", sts["totalOccurrences"], occ)})
	// 	})

	// 	if currentLine > -1 {
	// 		svs.RunAfterEvents(func() {
	// 			svs.ScrollToLine(currentLine)
	// 		})
	// 	}
	// }

	return
}

// updateStatusBar:
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

/*
* Clipboard handling ...
 */

func clipboardInit() {
	var err error
	if obj.clipboard, err = gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD); err != nil {
		log.Fatalf("clipboardInit: %s\n", err.Error())
	}
}

func clipboardCopyFromTextWin() {
	obj.clipboard.SetText(svs.GetText())
}

func clipboardPastToTextWin() {
	svs.Buffer.SetText(clipboardGet())
}

func clipboardGet() (clipboardContent string) {
	var err error
	if clipboardContent, err = obj.clipboard.WaitForText(); err != nil {
		log.Fatalf("clipboardGet: %s\n", err.Error())
	}
	return clipboardContent
}

func clipboardSet(clipboardContent string) {
	obj.clipboard.SetText(clipboardContent)
}

// Convert Entry to list of extensions.
func ExtSliceToOpt() {
	opt.ExtMask = []string{}

	tmpSliceStrings := strings.Split(GetEntryText(obj.entryExtMask), opt.ExtSep)
	for _, str := range tmpSliceStrings {
		str = strings.TrimSpace(str)
		if len(str) != 0 {
			opt.ExtMask = append(opt.ExtMask, str)
		}
	}
	obj.entryExtMask.SetText(strings.Join(opt.ExtMask, opt.ExtSep+" "))
}

func OptToExtSlice() {
	obj.entryExtMask.SetText(strings.Join(opt.ExtMask, opt.ExtSep+" "))
}
