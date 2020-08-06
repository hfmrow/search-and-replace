// miscFunctions.go

/*
	Source file auto-generated using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M

	This software use:
	- gotk3 that is licensed under the ISC License:
	  https://github.com/gotk3/gotk3/blob/master/LICENSE

	- Chroma — A general purpose syntax highlighter in pure Go, under the MIT License:
	  https://github.com/alecthomas/chroma/LICENSE

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
	"reflect"
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
	gitvtt "github.com/hfmrow/gtk3Import/textView/textTag"
	gitl "github.com/hfmrow/gtk3Import/tools"
)

// displayProgressBar: Hide or show some controls during a long search process.
func displayProgressBar(toggle bool) {
	// Finally i prefere blocking whole boxes except the bottom (exit button)
	// since aplying modifications depend on actual selected files and entries
	// when proceeded from found window
	// if this solution becomes unfriendly to use

	// TODO Keep previous selected files and entries, options ... to avoid this.

	mainObjects.MainGrid.SetSensitive(!toggle)
	mainObjects.MainTopGrig.SetSensitive(!toggle)
	mainObjects.listViewFiles.SetSensitive(!toggle)
	// mainObjects.btnFind.SetSensitive(!toggle)
	// mainObjects.fileChooserBtn.SetSensitive(!toggle)
	// mainObjects.btnScan.SetSensitive(!toggle)
	// mainObjects.spinButtonDepth.SetSensitive(!toggle)
	// mainObjects.btnReplaceInClipboard.SetSensitive(!toggle)
	// mainObjects.btnShowClipboard.SetSensitive(!toggle)
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
func showTextWin(text string, filename string, line int) {
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
	if occurrences, err = highlightText(text, line, mainObjects.textWinChkShowModifications.GetActive()); err == nil {
		textWinTitle.MainTitle = glsg.TruncatePath(filename, mainOptions.FilePathLength)
		textWinTitle.Update([]string{fmt.Sprintf("%s %d", sts["totalOccurrences"], occurrences)})
		BringToFront(mainObjects.textWin)
	}

	DlgErr("showTextWin:highlightText", err)
}

// highlightText: Highlight (found) and syntax highlight operations.
func highlightText(txtStr string, line int, doReplace bool) (occurrences int, err error) {

	// Set colors for TextView
	textViewRowNumber.TxtBgCol = mainOptions.TxtBgCol
	textViewRowNumber.TxtFgCol = mainOptions.TxtFgCol
	textViewRowNumber.NumFgCol = mainOptions.NumFgCol

	// store current text
	if currentText != txtStr {
		currentText = txtStr
	}
	if len(currentText) != 0 {

		var outTextBytes []byte

		if outTextBytes, err = onTheFlySearch([]byte(currentText), doReplace); err != nil {
			return
		}

		// Only compute new display if not already done.
		if /*true*/ !fileFoundSingle.HasBeenDisplayed() {

			// Remove found_background_color_lightgreen tag
			gitvtt.TagRemoveIfExists(textViewRowNumber.BuffTxt, found_background_color_lightgreen)

			textViewRowNumber.WaitForEventPending()

			if mainObjects.textWinChkSyntxHighlight.GetActive() {

				// Chroma syntax highlighting initialisation
				if !highlighter.Initialised() {
					if highlighter, err = ChromaHighlightNew(textViewRowNumber.BuffTxt, 1); err != nil {
						return
					}
					fileFoundSingle.HasBeenDisplayed(true)
				}
				textViewRowNumber.BufferDetach()

				// Let there be more light ... as a pig on the wings
				err = highlighter.Highlight(string(outTextBytes),
					mainObjects.textWinComboBoxLanguage.GetActiveText(),
					mainObjects.textWinComboBoxTextStyleChooser.GetActiveText())
				DlgErr("Highlight:gtkDirectToTextBuffer", err)

				switch highlighter.Formatter {
				case 1, 2:
					err = highlighter.ToTextBuff()
					DlgErr("Highlight:gtkTextBuffer/pango", err)
				default:
					err = highlighter.ToFile("out.php")
					DlgErr("Highlight:ToFile", err)
				}

				textViewRowNumber.BufferAttach()

			} else {
				textViewRowNumber.BufferDetach()
				highlighter.RemoveTags()
				textViewRowNumber.BuffTxt.SetText(string(outTextBytes))
				textViewRowNumber.BufferAttach()
			}
		}

		if !reflect.DeepEqual(fileFoundSingle.Pos.WordsPosIdx, previsouWordsPos) || doReplace {

			occurrences = fileFoundSingle.Occurrences
			if !doReplace {

				prop := map[string]interface{}{"background": "#AAFFAA"}
				tag := gitvtt.TagCreateIfNotExists(textViewRowNumber.BuffTxt, found_background_color_lightgreen, prop)
				DlgErr("highlightText:CreateTag", err)

				for _, line := range fileFoundSingle.Pos.FoundLinesIdx {
					for pIdx := 0; pIdx < len(line.FoundIdx); pIdx = pIdx + 2 {
						// Using line Index instead the Offsets to avoid issue where characters lenght > 1 byte (unicode)
						s := textViewRowNumber.BuffTxt.GetIterAtLineIndex(line.Number, line.FoundIdx[pIdx])
						e := textViewRowNumber.BuffTxt.GetIterAtLineIndex(line.Number, line.FoundIdx[pIdx+1])

						textViewRowNumber.BuffTxt.ApplyTag(tag, s, e)
					}
				}
			}
		}

		textViewRowNumber.ScrollToLine(line)
		textViewRowNumber.ColorBgRange(line, line+1)
	}

	currentText = txtStr

	// All errors have been handled previously
	err = nil

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
	mainObjects.clipboard.SetText(textViewRowNumber.GetText())
}

func clipboardPastToTextWin() {
	textViewRowNumber.SetText(clipboardGet())
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
