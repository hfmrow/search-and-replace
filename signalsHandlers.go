// signalsHandlers.go

/*
	Source file auto-generated on using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-19 H.F.M - Search And Replace
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

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
	glsg "github.com/hfmrow/genLib/strings"
	gltsbh "github.com/hfmrow/genLib/tools/bench"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gipf "github.com/hfmrow/gtk3Import/pixbuff"
	gitl "github.com/hfmrow/gtk3Import/tools"
)

// Signal handler dblClick ... doudble click on found result to display text preview
func findTreeViewDblClick(tw *gtk.TreeView) {
	var parentContent string
	var text []byte
	var line int
	var err error
	// var value *glib.Value
	var selContent string
	var iters []*gtk.TreeIter
	// var iter *gtk.TreeIter
	var path *gtk.TreePath

	mainObjects.textWinChkShowModifications.SetActive(false)
	mainObjects.textWinChkWrap.SetActive(false)

	if iters = tvsTree.GetSelectedIters(); len(iters) > 0 {
		selContent = tvsTree.GetColValue(iters[0], 0).(string)

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
				parentContent = tvsTree.GetColValuePath(path, 0).(string)
			}

			if text, err = ioutil.ReadFile(parentContent); err == nil {

				currFilename = parentContent // Filename passed to popup menu of the TextView

				showTextWin(string(text), glsg.TruncatePath(parentContent,
					mainOptions.FilePathLength), line-1)
			} else if os.IsNotExist(err) { // Skip error if not exist
				err = nil
			}
		}
	}
	DlgErr(sts["alert"], err)
}

// Signal handler dblClick on main listView ...
func listViewFilesRowActivated(tw *gtk.TreeView) {
	var err error
	var value *glib.Value
	var filename string
	var iters []*gtk.TreeIter
	var isTxt, gtLimit bool

	if iters = tvsList.GetSelectedIters(); len(iters) > 0 {

		if value, err = tvsList.ListStore.GetValue(iters[0], 3); err == nil { // Field 3: get full path
			filename, err = value.GetString() // Get selected file path
		}
		// Check for text file
		if isTxt, gtLimit, err = glfsft.IsTextFile(
			filename,
			mainOptions.FileMinSizeLimit,
			mainOptions.FileMaxSizeLimit); err == nil {

			if isTxt && gtLimit {
				if textWinTextToShowBytes, err = ioutil.ReadFile(filename); err == nil {

					currFilename = filename // Filename passed to popup menu of the TextView

					showTextWin(string(textWinTextToShowBytes), glsg.TruncatePath(filename, mainOptions.FilePathLength), 0)
				}
			}
		}
	}
	DlgErr(sts["missing"], err)
}

// Signal handler clicked ...
func btnFindClicked() {
	var err error
	var treeviewSelectedRows []string
	var entrySearchText, entryReplaceText string
	var removeEmptyResult = true
	var occurrences int

	// Only one at a time
	if btnFindInUse {
		return
	}
	btnFindInUse = true
	// Block operations while event pending
	for gtk.EventsPending() {
		gtk.MainIterationDo(true)
	}
	if treeviewSelectedRows, entrySearchText, entryReplaceText, err = getArguments(); err == nil {

		bench := gltsbh.BenchNew(false) // TODO

		displayProgressBar(true)
		var s = func() (err error) {

			bench.Lapse("Search") // TODO

			filesFoundMulti, occurrences, err = searchAndReplace(treeviewSelectedRows,
				entrySearchText, entryReplaceText, removeEmptyResult, false)
			return
		}

		var d = func() (err error) {

			bench.Lapse("Disp") // TODO

			displayResults(occurrences)
			displayProgressBar(false)

			bench.Stop() // TODO

			searchTime = bench.Lapses[0].StringShort
			dispTime = bench.Lapses[1].StringShort
			updateStatusBar()
			return
		}

		anim, err := gipf.GetPixBufAnimation(linearProgressHorzBlue)

		if err != nil {
			log.Fatalf("GetPixBufAnimation: %s\n", err.Error())
		}

		gifImage, err := gtk.ImageNewFromAnimation(anim)

		if err != nil {
			log.Fatalf("ImageNewFromAnimation: %s\n", err.Error())
		}

		pbs = gimc.ProgressGifNew(gifImage, mainObjects.mainBox, 1, s, d)
		go func() { err = pbs.StartGif() }()

	}

	DlgErr(sts["missing"], err)

	btnFindInUse = false
}

// textWinComboBoxTextStyleChooserChanged:
func textWinComboBoxTextStyleChooserChanged() {
	fileFoundSingle.HasBeenDisplayed(false)
	_, err := highlightText(currentText,
		textViewRowNumber.GetCurrentLineNb(),
		mainObjects.textWinChkShowModifications.GetActive())

	DlgErr("textWinComboBoxTextStyleChooserChanged:highlightText", err)
}

// Signal handler toggled ...
func textWinChkShowModificationsToggled() {
	var occurrences int
	var err error

	lastLine = textViewRowNumber.GetCurrentLineNb()

	if mainObjects.textWinChkShowModifications.GetActive() {
		textWinTextToShowBytes = []byte(textViewRowNumber.GetText())

		if len(gitl.GetEntryText(mainObjects.entrySearch)) == 0 {

			DlgErr(sts["missing"], fmt.Errorf(sts["nothingToSearch"]))

			mainObjects.textWinChkShowModifications.SetActive(false)
			return
		}
	}
	if occurrences, err = highlightText(string(textWinTextToShowBytes), lastLine, mainObjects.textWinChkShowModifications.GetActive()); err == nil {
		textWinTitle.Update([]string{fmt.Sprintf("%s %d", sts["totalOccurrences"], occurrences)})
	}
	if DlgErr("textWinChkShowModificationsToggled:highlightText", err) {
		mainObjects.textWinChkShowModifications.SetActive(false)
	}
}

// fileChooserBtnSelectionChanged
func fileChooserBtnSelectionChanged(fc *gtk.FileChooserButton) {
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
		showResults(&filesFoundMulti)
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

// chkFollowSymlinkDirToggled:
func chkFollowSymlinkDirToggled() {
	updateTreeViewFilesDisplay()
}

// // Signal handler changed ... FocusOut
// func entrySearchFocusOut(e *gtk.Entry) {
// 	genericEntryFocusOut(e)
// 	// return false // GDK_EVENT_PROPAGATE signal
// }

// // Signal handler changed ... FocusOut
// func entryReplaceFocusOut(e *gtk.Entry) {
// 	genericEntryFocusOut(e)
// 	// return false // GDK_EVENT_PROPAGATE signal
// }

// entryExtMaskEnterKeyPressed: Update data and disp on enter pressed
func entryExtMaskEnterKeyPressed(e *gtk.Entry) {
	ExtSliceToOpt()
	updateTreeViewFilesDisplay()
}

// Signal handler toggled ... (Wrap text)
func textWinChkWrapToggled() {
	if mainObjects.textWinChkWrap.GetActive() {
		textViewRowNumber.TextView.SetWrapMode(gtk.WRAP_WORD)
	} else {
		textViewRowNumber.TextView.SetWrapMode(gtk.WRAP_NONE)
	}
	textViewRowNumber.Update()
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
		mainObjects.chkWholeWord.SetActive(false)
		mainObjects.chkWildcard.SetActive(false)
	}
}

// Signal handler toggled ...
func findWinChkBackUpToggled() {
	mainOptions.MakeBackup = mainObjects.findWinChkBackUp.GetActive()
}

// Signal handler toggled ...
func chkWholeWordToggled() {
	if mainObjects.chkWholeWord.GetActive() {
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
	var err error
	var occurences int

	if gidg.DialogMessage(mainObjects.mainWin, "error", sts["confirm"], "\n\n"+sts["proceed"], "", sts["ok"], sts["cancel"]) == 1 {
		return
	}

	if treeviewSelectedRows, entrySearchText, entryReplaceText, err := getArguments(); err == nil {
		_, occurences, err = searchAndReplace(treeviewSelectedRows, entrySearchText, entryReplaceText, true, true) // Do replace in files
	}

	if DlgErr(sts["unexpected"], err) {
		return
	}
	genericHideWindow(mainObjects.findWin) // Hide find window
	genericHideWindow(mainObjects.textWin) // Hide preview text window

	gidg.DialogMessage(mainObjects.mainWin, "info", sts["done"],
		fmt.Sprintf("\n\n%s: %d\n", sts["totalModified"], occurences), "", "Ok")

	if occurences > 0 {
		// Reset Init single found structure to avoid always display found pattern(s)
		fileFoundSingle = glfsft.SearchAndReplaceNew([]byte{}, "", "")
	}
}

// btnShowClipboardClicked:
func btnShowClipboardClicked() {
	if mainObjects.entrySearch.GetTextLength() != 0 {
		textWinTextToShowBytes = []byte(clipboardGet())
		showTextWin(string(textWinTextToShowBytes), sts["clpbrdPreview"], 0)
	}
}

// btnReplaceInClipboardClicked:
func btnReplaceInClipboardClicked() {

	if len(gitl.GetEntryText(mainObjects.entrySearch)) == 0 {
		gidg.DialogMessage(mainObjects.mainWin,
			"warning", sts["missing"],
			"\n\n"+sts["nothingToSearch"], "", "Ok")
	} else {

		outTextBytes, err := onTheFlySearch([]byte(clipboardGet()), true)
		if err != nil {
			DlgErr("btnReplaceInClipboardClicked:onTheFlySearch", err)
		} else {
			clipboardSet(string(outTextBytes))
		}
	}
}

// findWinCancelBtnClicked:
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
	if entry, err := e.GetText(); !DlgErr(sts["alert"], err) {
		// Sanitize entry
		entry = strings.Replace(fmt.Sprintf("%q", entry), `"`, ``, -1)
		e.SetText(entry)
		if mainObjects.textWin.GetVisible() {
			listViewFilesRowActivated(mainObjects.listViewFiles)
		}
		// mainObjects.textWin.SetVisible(true)
	}
}

// imgTop handler release signal (Display about box)
func imgTopReleaseEvent() {
	mainOptions.AboutOptions.Width = 400
	mainOptions.AboutOptions.ImageOkButtonSize = 24
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

// OptionButtonDoneClicked
func OptionButtonDoneClicked() {
	genericHideWindow(mainObjects.OptionsWindow)
}

// OptionsEntryFileSizeFocusOutEvent:
func OptionsEntryFileSizeChanged(e *gtk.Entry) {
	var err error
	var entry string
	if entry, err = e.GetText(); err == nil && len(entry) != 0 {
		entry = strings.TrimSpace(entry)
		if _, err = strconv.Atoi(entry); err == nil {
			// Value are sanitized, so, we only need to retrieve them
			mainOptions.FileMaxSizeLimit = int64(gitl.GetEntryTextAsInt(mainObjects.OptionsEntryMaxFileSize))
			mainOptions.FileMinSizeLimit = int64(gitl.GetEntryTextAsInt(mainObjects.OptionsEntryMinFileSize))
			mainObjects.OptionButtonDone.SetSensitive(true)
		}
	}
	if err != nil {
		gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
	}
	e.SetText(glsg.RemoveNonNum(entry))
	return
}

// MainButtonOptionsClicked:
func MainButtonOptionsClicked() {
	mainObjects.OptionsWindow.ShowAll()
	mainObjects.OptionsWindow.SetKeepAbove(true)
}
