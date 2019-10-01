// signalsHandlers.go

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	glfsft "github.com/hfmrow/genLib/files/fileText"
	gltsbh "github.com/hfmrow/genLib/tools/bench"

	gidg "github.com/hfmrow/gtk3Import/dialog"
)

// Signal handler dblClick ... doudble click on found result to display text preview
func findTreeViewDblClick(tw *gtk.TreeView) {
	var parentContent string
	var text []byte
	var line int
	var err error
	var value *glib.Value
	var selContent string
	var iters []*gtk.TreeIter
	var iter *gtk.TreeIter
	var path *gtk.TreePath

	if iters, err = tvsTree.GetSelectedIters(); err == nil {
		if value, err = tvsTree.TreeStore.GetValue(iters[0], 0); err == nil { // Field 0: content
			selContent, err = value.GetString()
			if path, err = tvsTree.TreeStore.GetPath(iters[0]); err == nil {

				var regLine = regexp.MustCompile(`[><]`)
				currentLine := regLine.Split(selContent, -1)

				// Retrieve line number if xist ...
				if len(currentLine) > 1 {
					line, _ = strconv.Atoi(strings.TrimSpace(currentLine[2]))
				}

				// Get the parent node of the current node or itself if it already is
				path.Up()
				if len(path.String()) == 0 {
					parentContent = selContent
				} else {
					iter, _ = tvsTree.TreeStore.GetIter(path)
					value, _ = tvsTree.TreeStore.GetValue(iter, 0)
					parentContent, _ = value.GetString()
				}

				if text, err = ioutil.ReadFile(parentContent); err == nil {
					showTextWin(string(text), truncatePath(parentContent,
						mainOptions.FilePathLength), line)
				}
			}
		}
	}

	if err != nil {
		gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
	}
}

// Signal handler dblClick on main listView ...
func listViewFilesRowActivated(tw *gtk.TreeView) {
	var err error
	var value *glib.Value
	var filename string
	var iters []*gtk.TreeIter
	var isTxt, gtLimit bool

	if iters, err = tvsList.GetSelectedIters(); err == nil {
		if value, err = tvsList.ListStore.GetValue(iters[0], 3); err == nil { // Field 3: get full path
			filename, err = value.GetString() // Get selected file path
		}
		// Check for text file
		if isTxt, gtLimit, err = glfsft.IsTextFileSimple(filename,
			mainOptions.FileSizeLimit,
			mainOptions.LineEndThreshold,
			mainOptions.OverCharsThreshold); err == nil {
			if isTxt && gtLimit {
				textWinTextToShowBytes, err = ioutil.ReadFile(filename)
				// mainObjects.textWin.SetModal(true)
				showTextWin(string(textWinTextToShowBytes), truncatePath(filename, mainOptions.FilePathLength), 0)
			}
		}
	}
	if err != nil {
		gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
	}
}

// Signal handler clicked ...
func btnFindClicked() {
	// TODO /* BENCH Search Time */
	bench := new(gltsbh.Bench)
	bench.Lapse("Start")
	/* BENCH */

	var err error
	var found bool
	var occurrences, countFiles int
	removeEmptyResult := true
	entrySearchText := getEntryText(mainObjects.entrySearch)
	entryReplaceText := getEntryText(mainObjects.entryReplace)

	if len(entrySearchText) != 0 {

		mainOptions.mainFound, occurrences, err = Find(entrySearchText,
			entryReplaceText, removeEmptyResult)
		if err != nil {
			gidg.DialogMessage(mainObjects.mainWin,
				"warning", sts["alert"], "\n\n"+err.Error(), "", "Ok")
		} else { // Something was found ?
			if occurrences > 0 {
				countFiles = showResults(&mainOptions.mainFound)

				findWinTitle.Update([]string{fmt.Sprintf("%s %d %s %d %s", sts["totalOccurrences"], occurrences, sts["in"], countFiles, sts["file"])})

				x, y := mainObjects.mainWin.GetPosition()
				mainObjects.findWin.SetModal(true)
				mainObjects.findWin.Move(x+mainOptions.CascadeDepth, y+mainOptions.CascadeDepth)
				mainObjects.findWin.Resize(mainObjects.mainWin.GetSize())
				mainObjects.findWin.Show()
				found = true

				bench.Stop()
				searchTime = fmt.Sprintf("%dm %ds %dms", bench.NumTime[0].Min, bench.NumTime[0].Sec, bench.NumTime[0].Ms)

			} // Nope !
			if !found && !applyChanges {
				gidg.DialogMessage(mainObjects.mainWin,
					"warning", sts["missing"],
					"\n\n"+sts["notFound"], "", "Ok")
			}
		}
	} else {
		gidg.DialogMessage(mainObjects.mainWin,
			"warning", sts["missing"],
			"\n\n"+sts["nothingToSearch"], "", "Ok")
	}
	updateStatusBar()
	searchTime = ""
}

// Signal handler toggled ...
func textWinChkShowModificationsToggled() {
	var occurrences int
	var err error
	var buff *gtk.TextBuffer
	var txt string

	if buff, err = mainObjects.textWinTextview.GetBuffer(); err == nil {
		if mainObjects.textWinChkShowModifications.GetActive() {
			entrySearchText := getEntryText(mainObjects.entrySearch)
			if len(entrySearchText) == 0 {
				gidg.DialogMessage(mainObjects.mainWin,
					"warning", sts["missing"],
					"\n\n"+sts["nothingToSearch"], "", "Ok")
				mainObjects.textWinChkShowModifications.SetActive(false)
				return // No action, back to previous state
			} else {
				if txt, err = buff.GetText(buff.GetStartIter(), buff.GetEndIter(), true); err == nil {
					textWinTextToShowBytes = []byte(txt)
					var outText []byte
					if err = onTheFlyReplace(textWinTextToShowBytes, &outText, entrySearchText); err == nil {
						buff.SetText(string(outText))
					}
				}
			}
		} else {
			mark := buff.GetInsert()
			iter := buff.GetIterAtMark(mark)

			if occurrences, err = highlightText(buff, string(textWinTextToShowBytes), iter.GetLine()); err == nil {
				textWinTitle.Update([]string{fmt.Sprintf("%s %d", sts["totalOccurrences"], occurrences)})
			}
		}
	}
	if err != nil {
		gidg.DialogMessage(mainObjects.mainWin, "warning", sts["alert"], "\n\n"+err.Error(), "", "Ok")
		mainObjects.textWinChkShowModifications.SetActive(false)
	}
}

// fileChooserBtnSelectionChanged
func fileChooserBtnSelectionChanged(fc *gtk.FileChooserButton) {
	// fmt.Println(mainObjects.fileChooserBtn.GetFilename())
	if mainObjects.switchFileChooserButton.GetActive() {
		filename := mainObjects.fileChooserBtn.GetFilename()
		if _, err := os.Stat(filename); err == nil {
			currentInFilesList = append(currentInFilesList, filename)
			updateTreeViewFilesDisplay()
		}
	}
}

// findWinChkDispForbFilesToggled
func findWinChkDispForbFilesToggled(chk *gtk.CheckButton) {
	if mainObjects.findWin.GetVisible() {
		showResults(&mainOptions.mainFound)
	}
}

// btnScanClicked:
func btnScanClicked() {
	filename := mainObjects.fileChooserBtn.GetFilename()
	if _, err := os.Stat(filename); err == nil {
		currentInFilesList = currentInFilesList[:0]
		currentInFilesList = append(currentInFilesList, filename)
	}
	updateTreeViewFilesDisplay()
}

// spinButtonDepthChangeValue:
func spinButtonDepthValueChanged() {
	updateTreeViewFilesDisplay()
}

// Signal handler changed ... FocusOut - to trim entry
func entryExtMaskFocusOut() bool {
	updateTreeViewFilesDisplay()
	return false // GDK_EVENT_PROPAGATE signal
}

// switchFileChooserButtonStateSet:
func switchFileChooserButtonStateSet() bool {
	// mainObjects.fileChooserBtn.SetSensitive(mainObjects.switchFileChooserButton.GetActive())
	// mainObjects.btnScan.SetSensitive(mainObjects.switchFileChooserButton.GetActive())
	// return fileChooserDoUpdt
	return false
}

// switchFileChooserButtonEventAfter
func switchFileChooserButtonEventAfter() {
	// if !fromDnD {
	// 	mainObjects.switchFileChooserButton.SetActive(true)
	// }
}

// chkFollowSymlinkDirToggled:
func chkFollowSymlinkDirToggled() {
	updateTreeViewFilesDisplay()
}

// Signal handler changed ... FocusOut
func entrySearchFocusOut(e *gtk.Entry) {
	genericEntryFocusOut(e)
}

// Signal handler changed ... FocusOut
func entryReplaceFocusOut(e *gtk.Entry) {
	genericEntryFocusOut(e)
}

// entryExtMaskEnterKeyPressed: Update data and disp on enter pressed
func entryExtMaskEnterKeyPressed(e *gtk.Entry) {
	ExtSliceToOpt()
	updateTreeViewFilesDisplay()
}

// Signal handler toggled ... (Wrap text)
func textWinChkWrapToggled() {
	if mainObjects.textWinChkWrap.GetActive() {
		mainObjects.textWinScrolledwindowNumbers.SetVisible(false)
		mainObjects.textWinTextview.SetProperty("left-margin", 6)
		mainObjects.textWinTextview.SetWrapMode(gtk.WRAP_WORD)
	} else {
		mainObjects.textWinScrolledwindowNumbers.SetVisible(true)
		mainObjects.textWinTextview.SetWrapMode(gtk.WRAP_NONE)
	}
}

// Signal handler toggled ...
func chkCharacterClassToggled() {
	if mainObjects.chkCharacterClass.GetActive() {
		mainObjects.chkUseEscapeChar.SetActive(false)
		mainObjects.chkRegex.SetActive(false)
		mainObjects.chkWildcard.SetActive(false)
		mainObjects.chkCharacterClassStrictMode.SetVisible(true)
	} else {
		mainObjects.chkCharacterClassStrictMode.SetVisible(false)
	}
}

// Signal handler toggled ...
func chkUseEscapeCharToggled() {
	if mainObjects.chkUseEscapeChar.GetActive() {
		mainObjects.chkCharacterClass.SetActive(false)
		mainObjects.chkRegex.SetActive(false)
	}
}

// Signal handler toggled ...
func chkRegexToggled() {
	if mainObjects.chkRegex.GetActive() {
		mainObjects.chkCharacterClass.SetActive(false)
		mainObjects.chkUseEscapeChar.SetActive(false)
		mainObjects.chkCaseSensitive.SetActive(false)
		mainObjects.chkWoleWord.SetActive(false)
		mainObjects.chkWildcard.SetActive(false)
	}
}

// Signal handler toggled ...
func findWinChkBackUpToggled() {
	mainOptions.MakeBackup = mainObjects.findWinChkBackUp.GetActive()
}

// Signal handler toggled ...
func chkWoleWordToggled() {
	if mainObjects.chkWoleWord.GetActive() {
		mainObjects.chkRegex.SetActive(false)
	}
}

// Signal handler toggled ...
func chkCaseSensitiveToggled() {
	if mainObjects.chkCaseSensitive.GetActive() {
		mainObjects.chkRegex.SetActive(false)
	}
}

// Signal handler toggled ...
func chkWildcardToggled() {
	if mainObjects.chkWildcard.GetActive() {
		mainObjects.chkRegex.SetActive(false)
		mainObjects.chkCharacterClass.SetActive(false)
	}
}

// Signal handler clicked ...
func findWinReplaceBtnClicked() {
	applyChanges = true
	btnFindClicked() // Do replace in files
	applyChanges = false
	genericHideWindow(mainObjects.findWin) // Hide find window
	// fileChooserBtnClicked()                // Reload files list
}

// Signal handler clicked ...
func btnShowClipboardClicked() {
	if mainObjects.entrySearch.GetTextLength() != 0 {
		textWinTextToShowBytes = []byte(clipboardGet())
		showTextWin(string(textWinTextToShowBytes), sts["clpbrdPreview"], 0)
	}
}

// Signal handler clicked ...
func btnReplaceInClipboardClicked() {
	entrySearchText, err := mainObjects.entrySearch.GetText()
	if err != nil {
		gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
	} else if len(entrySearchText) == 0 {
		gidg.DialogMessage(mainObjects.mainWin,
			"warning", sts["missing"],
			"\n\n"+sts["nothingToSearch"], "", "Ok")
	} else {
		outText := new([]byte)
		err := onTheFlyReplace([]byte(clipboardGet()), outText, entrySearchText)
		if err != nil {
			gidg.DialogMessage(mainObjects.mainWin, "warning", sts["alert"], "\n\n"+err.Error(), "", "Ok")
		} else {
			clipboardSet(string(*outText))
		}
	}
}

// Signal handler clicked ...
func findWinCancelBtnClicked() {
	genericHideWindow(mainObjects.findWin)
	genericHideWindow(mainObjects.textWin)
}

// Signal handler clicked ...
func textWinBtnDoneClicked() {
	mainObjects.textWinChkShowModifications.SetActive(false)
	genericHideWindow(mainObjects.textWin)
}

// Signal handler delete_event (hidding window)
func genericHideWindow(w *gtk.Window) bool {
	if w.GetVisible() {
		w.SetModal(false)
		w.Hide()
	}
	return true
}

// Used with gtk.Entry objects.
func genericEntryFocusOut(e *gtk.Entry) {
	if entry, err := e.GetText(); err == nil {
		entry = strings.Replace(fmt.Sprintf("%q", entry), `"`, ``, -1)
		e.SetText(entry)
	} else {
		gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
	}
}

// imgTop handler release signal (Display about box)
func imgTopReleaseEvent() {
	mainOptions.AboutOptions.Show()
}

// Signal handler on Exit window ... Saving options before quit
func mainWinOnExit() {
	var err error
	mainOptions.UpdateOptions()
	if err = mainOptions.Write(); err != nil {
		log.Fatalf("Error writing option file: %s\n", err.Error())
	}
}
