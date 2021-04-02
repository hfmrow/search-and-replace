// miscFunctions.go

/*
	Source file auto-generated using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-19 H.F.M - Search And Replace
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

	humanize "github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	glfsft "github.com/hfmrow/genLib/files/fileText"
	glsg "github.com/hfmrow/genLib/strings"
	gltsbh "github.com/hfmrow/genLib/tools/bench"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gipo "github.com/hfmrow/gtk3Import/pango"
	gipops "github.com/hfmrow/gtk3Import/pango/pangoSimple"
	gitl "github.com/hfmrow/gtk3Import/tools"
)

// displayProgressBar: Hide or show some controls during a long search process.
func displayProgressBar(toggle bool) {

	mainObjects.MainGrid.SetSensitive(!toggle)
	mainObjects.MainTopGrig.SetSensitive(!toggle)
	mainObjects.listViewFiles.SetSensitive(!toggle)
}

// updateTreeViewFilesDisplay: Display files to treeView
func updateTreeViewFilesDisplay() {
	var err error
	var stat os.FileInfo

	bench := gltsbh.BenchNew(false)
	bench.Lapse()

	if !fromDnD {
		currentInFilesList = currentInFilesList[:0]
		currentInFilesList = append(currentInFilesList, mainObjects.fileChooserBtn.GetFilename())
	} else {
		if len(currentInFilesList) == 1 {
			if stat, err = os.Stat(currentInFilesList[0]); err == nil {
				if stat.IsDir() {
					mainObjects.fileChooserBtn.SetFilename(currentInFilesList[0])
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
		scanTime = bench.Lapses[0].StringShort
		displayFiles()
	}
}

// displayFiles: display formatted files to treeview.
func displayFiles() {
	mainOptions.UpdateOptions()

	if tvsList != nil {
		bench := gltsbh.BenchNew(false)
		bench.Lapse()

		// Detach listStore & clean before fill it
		tvsList.StoreDetach()
		tvsList.Clear()

		var ok bool
		var err error

		for _, rowFile := range toDispFileList {
			for _, msk := range mainOptions.ExtMask { // Check for ext matching well.
				if ok, err = filepath.Match(msk, rowFile.FileInfo.Name()); err == nil {
					if ok {
						tvsList.AddRow(nil,
							html.EscapeString(rowFile.FileInfo.Name()),
							humanize.Bytes(uint64(rowFile.FileInfo.Size())),
							humanize.Time(rowFile.FileInfo.ModTime()),
							html.EscapeString(rowFile.Filename),
							rowFile.FileInfo.Size(),
							rowFile.FileInfo.ModTime().Unix())
					}
				} else {
					DlgErr("_MatchExt", err)
				}
			}
		}

		// Attach listStore
		tvsList.StoreAttach()

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

	pc := gipops.PangoColorNew()

	// Detach & clear treeStore
	tvsTree.StoreDetach()
	tvsTree.Clear()

	for _, result := range *mainFound {
		prevLineNbr = -1
		total++
		if result.Occurrences > 0 {
			text++
			/* Add parent */
			if iter, err = tvsTree.AddRow(nil, result.FileName); err == nil {

				// Make markup
				pm := gipo.PangoMarkup{}
				pm.Init(string(result.SearchAndRepl.TextBytes)) // Text
				pm.AddPosition(result.SearchAndRepl.Pos.WordsPosIdx...)
				pm.AddTypes([][]string{{"bgc", "#AAFFAA"}}...) // Light green
				text := pm.MarkupAtPos()

				/* Add Childs */
				linesStr := strings.Split(text, result.SearchAndRepl.Pos.Eol)
				for idxLine := 0; idxLine < len(result.SearchAndRepl.Pos.FoundLinesIdx); idxLine++ {

					lineNbr = result.SearchAndRepl.Pos.FoundLinesIdx[idxLine].Number
					if prevLineNbr != lineNbr { // Avoid displaying duplicated rows
						pm.Init(fmt.Sprintf(" %v ", lineNbr+1))        // Line number
						pm.AddTypes([][]string{{"bgc", "#E4DDDD"}}...) // Light grey
						lineNbStr := pm.Markup()

						tvsTree.AddRow(iter, lineNbStr+"\t"+string(linesStr[lineNbr]))
					}
					prevLineNbr = lineNbr
				}
				countFiles++
			}
		} else if mainObjects.findWinChkDispForbFiles.GetActive() && result.NotTextFile {
			binary++
			// Markup bad file: File size too short or Binary files warning ...
			pm := gipo.PangoMarkup{}
			pm.Init(string(result.FileName)) // Filename
			pm.AddPosition(result.SearchAndRepl.Pos.WordsPosIdx...)
			pm.AddTypes([][]string{{"stc", pc.Red}}...)
			filename := pm.Markup()

			/* Add parent */
			if iter, err = tvsTree.AddRow(nil, filename); err == nil {
				pm.Init(string(result.SearchAndRepl.TextBytes)) // text with reason why it's not displayed
				pm.AddTypes([][]string{{"fgc", pc.Red}}...)
				desciption := pm.Markup()

				/* Add Child */
				tvsTree.AddRow(iter, desciption)
			}
		}
	}
	// Attach treeStore
	tvsTree.StoreAttach()
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
		mainObjects.findWin.SetModal(false)
		mainObjects.textWin.SetKeepAbove(false)
		x, y := mainObjects.mainWin.GetPosition()
		mainObjects.textWin.Move(x+(mainOptions.CascadeDepth*2), y+(mainOptions.CascadeDepth*2))
		mainObjects.textWin.Resize(mainObjects.mainWin.GetSize())
		alreadyPlacedPrevWin = true
	}
	// Set text
	if occurrences, err = highlightText(text, mainObjects.textWinChkShowModifications.GetActive()); err == nil {
		textWinTitle.MainTitle = glsg.TruncatePath(filename, mainOptions.FilePathLength)
		textWinTitle.Update([]string{fmt.Sprintf("%s %d", sts["totalOccurrences"], occurrences)})

		svs.BringToFront()
	}

	DlgErr("showTextWin:highlightText", err)
}

// highlightText: Highlight (found) and syntax highlight operations.
func highlightText(txtStr string, doReplace bool) (occurrences int, err error) {

	// Set colors for GtkSourceView
	svs.TxtBgCol = mainOptions.TxtBgCol
	svs.TxtFgCol = mainOptions.TxtFgCol

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
			svs.TextSearch = gitl.GetEntryText(mainObjects.entryReplace)
			svs.Search(svs.Buffer.GetStartIter(), false)
			svs.TextSearch = tmpTextSearch

		} else {
			// SourceView search part: used to highlight found or replaced patterns
			svs.TextSearch = gitl.GetEntryText(mainObjects.entrySearch)
			svs.UseRegexp = mainObjects.chkRegex.GetActive()
			svs.WordBoundaries = mainObjects.chkWholeWord.GetActive()
			svs.CaseSensitive = mainObjects.chkCaseSensitive.GetActive()
			svs.Search(svs.Buffer.GetStartIter(), false)
		}

		textWinTitle.Update([]string{fmt.Sprintf("%s %d", sts["totalOccurrences"], fileFoundSingle.Occurrences)})

		if currentLine > -1 {
			svs.RunAfterEvents(func() {
				svs.ScrollToLine(currentLine)
			})
		}
	}

	// // Set colors for GtkSourceView
	// svs.TxtBgCol = mainOptions.TxtBgCol
	// svs.TxtFgCol = mainOptions.TxtFgCol

	// // store current text
	// if currentText != txtStr || (!doReplace && currentTextChanged) {
	// 	currentText = txtStr
	// 	svs.SetText(txtStr)
	// 	currentTextChanged = false
	// }

	// if len(currentText) != 0 {

	// 	svs.TextSearch = gitl.GetEntryText(mainObjects.entrySearch)
	// 	svs.UseRegexp = mainObjects.chkRegex.GetActive()
	// 	svs.WordBoundaries = mainObjects.chkWholeWord.GetActive()
	// 	svs.CaseSensitive = mainObjects.chkCaseSensitive.GetActive()

	// 	if !doReplace {
	// 		svs.Search(svs.Buffer.GetStartIter(), false)
	// 	} else {
	// 		svs.TextReplace = gitl.GetEntryText(mainObjects.entryReplace)
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
	if mainObjects.clipboard, err = gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD); err != nil {
		log.Fatalf("clipboardInit: %s\n", err.Error())
	}
}

func clipboardCopyFromTextWin() {
	mainObjects.clipboard.SetText(svs.GetText())
}

func clipboardPastToTextWin() {
	svs.Buffer.SetText(clipboardGet())
}

func clipboardGet() (clipboardContent string) {
	var err error
	if clipboardContent, err = mainObjects.clipboard.WaitForText(); err != nil {
		log.Fatalf("clipboardGet: %s\n", err.Error())
	}
	return clipboardContent
}

func clipboardSet(clipboardContent string) {
	mainObjects.clipboard.SetText(clipboardContent)
}

// Convert Entry to list of extensions.
func ExtSliceToOpt() {
	mainOptions.ExtMask = []string{}

	tmpSliceStrings := strings.Split(gitl.GetEntryText(mainObjects.entryExtMask), mainOptions.ExtSep)
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
