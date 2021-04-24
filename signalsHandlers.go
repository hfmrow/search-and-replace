// signalsHandlers.go

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
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hfmrow/gotk3_gtksource/source"
)

// Signal handler dblClick ... doudble click on found result to display text preview
func findTreeViewDblClick(tw *gtk.TreeView) {
	var (
		parentContent string
		text          []byte
		err           error
		selContent    string
		iters         []*gtk.TreeIter
		path          *gtk.TreePath
	)

	obj.textWinChkShowModifications.SetActive(false)
	obj.textWinChkWrap.SetActive(false)

	if iters = tvsTree.GetSelectedIters(); len(iters) > 0 {
		selContent = tvsTree.GetColValue(iters[0], opt.mapTreeStore["Filename"]).(string)

		// TODO Use treeview stored line number instead of the one part of the displayed row

		if path, err = tvsTree.TreeStore.GetPath(iters[0]); err == nil {
			var regLine = regexp.MustCompile(`[><]`)
			newCurrentLine := regLine.Split(selContent, -1)

			// Retrieve line number if exist ...
			if len(newCurrentLine) > 1 {
				currentLine, _ = strconv.Atoi(strings.TrimSpace(newCurrentLine[2]))
			}
			// Get the parent node of the current node or itself if it already is
			path.Up()
			if len(path.String()) == 0 {
				parentContent = selContent
			} else {
				parentContent = tvsTree.GetColValuePath(path, opt.mapTreeStore["Filename"]).(string)
			}

			if text, err = ioutil.ReadFile(parentContent); err == nil {

				currFilename = parentContent // Filename passed to popup menu of the TextView

				showTextWin(string(text), TruncatePath(parentContent,
					opt.FilePathLength))
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
		if isTxt, gtLimit, err = IsTextFile(
			filename,
			opt.FileMinSizeLimit,
			opt.FileMaxSizeLimit); err == nil {

			if isTxt && gtLimit {
				if textWinTextToShowBytes, err = ioutil.ReadFile(filename); err == nil {

					currFilename = filename // Filename passed to popup menu of the TextView

					showTextWin(string(textWinTextToShowBytes), TruncatePath(filename, opt.FilePathLength))
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

		bench := BenchNew(false) // TODO

		displayProgressBar(true)

		pbs.Init(func() (err error) {

			bench.Lapse("Search") // TODO

			filesFoundMulti, occurrences, err = searchAndReplace(treeviewSelectedRows,
				entrySearchText, entryReplaceText, removeEmptyResult, false)
			return
		},
			func() (err error) {

				bench.Lapse("Disp") // TODO

				displayResults(occurrences)
				displayProgressBar(false)

				bench.Stop() // TODO

				searchTime = bench.Lapses[0].StringShort
				dispTime = bench.Lapses[1].StringShort
				updateStatusBar()
				return
			})

		go func() { err = pbs.StartGif() }()
	}

	DlgErr(sts["missing"], err)

	btnFindInUse = false
}

// Signal handler toggled ...
func textWinChkShowModificationsToggled() {

	var err error

	if obj.textWinChkShowModifications.GetActive() {
		textWinTextToShowBytes = []byte(svs.GetText())

		if len(GetEntryText(obj.entrySearch)) == 0 {

			DlgErr(sts["missing"], fmt.Errorf(sts["nothingToSearch"]))

			obj.textWinChkShowModifications.SetActive(false)
			return
		}
	}

	_, err = highlightText(string(textWinTextToShowBytes), obj.textWinChkShowModifications.GetActive())

	if DlgErr("textWinChkShowModificationsToggled:highlightText", err) {
		obj.textWinChkShowModifications.SetActive(false)
	}
}

// fileChooserBtnSelectionChanged
func fileChooserBtnSelectionChanged(fc *gtk.FileChooserButton) {
	if obj.switchFileChooserButton.GetActive() {
		filename := obj.fileChooserBtn.GetFilename()
		if _, err := os.Stat(filename); err == nil {
			currentInFilesList = append(currentInFilesList, filename)
			updateTreeViewFilesDisplay()
		}
	}
}

func findWinBtnDeselectClicked(btn *gtk.Button) {
	tvsTree.ChangeCheckState(opt.mapTreeStore["Toggle"], false, false)
}

func findWinBtnInvertSelClicked(btn *gtk.Button) {
	tvsTree.ChangeCheckState(opt.mapTreeStore["Toggle"], false, true)
}

/*
 * Checkbuttons
 */
func findWinChkExpandAllToggled(chk *gtk.CheckButton) {
	opt.ExpandAll = chk.GetActive()
	tvsTree.ExpandAll(!opt.ExpandAll)
}

// findWinChkDispForbFilesToggled
func findWinChkDispForbFilesToggled(chk *gtk.CheckButton) {
	if obj.findWin.GetVisible() {
		showResults(&filesFoundMulti)
	}
}

// btnScanClicked:
func btnScanClicked() {
	filename := obj.fileChooserBtn.GetFilename()
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

// entryExtMaskEnterKeyPressed: Update data and disp on enter pressed
func entryExtMaskEnterKeyPressed(e *gtk.Entry) {
	ExtSliceToOpt()
	updateTreeViewFilesDisplay()
}

// textWinChkWrapToggled:  (Wrap text)
func textWinChkWrapToggled(chk *gtk.CheckButton) {

	if svs != nil {
		if chk.GetActive() {
			svs.View.SetWrapMode(gtk.WRAP_WORD_CHAR)
		} else {
			svs.View.SetWrapMode(gtk.WRAP_NONE)
		}
	}
}

func textWinChkSyntxHighlightToggled(chk *gtk.CheckButton) {
	svs.Buffer.SetHighlightSyntax(chk.GetActive())
}

// Signal handler toggled ...
func chkCharacterClassToggled() {
	if obj.chkCharacterClass.GetActive() {
		obj.chkCaseSensitive.SetSensitive(false)
		obj.chkWholeWord.SetSensitive(false)
		obj.chkUseEscapeChar.SetSensitive(false)
		obj.chkRegex.SetSensitive(false)
		obj.chkWildcard.SetSensitive(false)
		obj.chkCaseSensitive.SetActive(false)
		obj.chkWholeWord.SetActive(false)
		obj.chkUseEscapeChar.SetActive(false)
		obj.chkRegex.SetActive(false)
		obj.chkWildcard.SetActive(false)
		obj.chkCharacterClassStrictMode.SetVisible(true)
	} else {
		obj.chkCaseSensitive.SetSensitive(true)
		obj.chkWholeWord.SetSensitive(true)
		obj.chkUseEscapeChar.SetSensitive(true)
		obj.chkRegex.SetSensitive(true)
		obj.chkWildcard.SetSensitive(true)
		obj.chkCharacterClassStrictMode.SetVisible(false)
	}
}

// Signal handler toggled ...
func chkUseEscapeCharToggled() {
	if obj.chkUseEscapeChar.GetActive() {
		obj.chkCharacterClass.SetSensitive(false)
		obj.chkRegex.SetSensitive(false)
		obj.chkWholeWord.SetSensitive(false)
		obj.chkCharacterClass.SetActive(false)
		obj.chkRegex.SetActive(false)
		obj.chkWholeWord.SetActive(false)
	} else {
		obj.chkCharacterClass.SetSensitive(true)
		obj.chkRegex.SetSensitive(true)
		obj.chkWholeWord.SetSensitive(true)
	}
}

// Signal handler toggled ...
func chkRegexToggled() {
	if obj.chkRegex.GetActive() {
		obj.chkCharacterClass.SetSensitive(false)
		obj.chkUseEscapeChar.SetSensitive(false)
		obj.chkCaseSensitive.SetSensitive(false)
		obj.chkWholeWord.SetSensitive(false)
		obj.chkWildcard.SetSensitive(false)
		obj.chkCharacterClass.SetActive(false)
		obj.chkUseEscapeChar.SetActive(false)
		obj.chkCaseSensitive.SetActive(false)
		obj.chkWholeWord.SetActive(false)
		obj.chkWildcard.SetActive(false)
	} else {
		obj.chkCharacterClass.SetSensitive(true)
		obj.chkUseEscapeChar.SetSensitive(true)
		obj.chkCaseSensitive.SetSensitive(true)
		obj.chkWholeWord.SetSensitive(true)
		obj.chkWildcard.SetSensitive(true)
	}
}

// Signal handler toggled ...
func findWinChkBackUpToggled() {
	opt.MakeBackup = obj.findWinChkBackUp.GetActive()
}

// Signal handler toggled ...
func chkWholeWordToggled() {
	if obj.chkWholeWord.GetActive() {
		obj.chkRegex.SetActive(false)
	}
}

// Signal handler toggled ...
func chkCaseSensitiveToggled() {
	if obj.chkCaseSensitive.GetActive() {
		obj.chkRegex.SetActive(false)
	}
}

// Signal handler toggled ...
func chkWildcardToggled() {
	if obj.chkWildcard.GetActive() {
		obj.chkRegex.SetActive(false)
		obj.chkCharacterClass.SetActive(false)
	}
}

// Signal handler clicked ...
func findWinReplaceBtnClicked() {
	var (
		err error
		occurences,
		occReplaced int
	)

	// Confirmation message
	if DialogMessage(obj.mainWin, "error", sts["confirm"], "\n\n"+sts["proceed"], "", sts["ok"], sts["cancel"]) == 1 {
		return
	}

	// if treeviewSelectedRows, entrySearchText, entryReplaceText, err := getArguments(); err == nil {
	// 	_, occurences, err = searchAndReplace(treeviewSelectedRows, entrySearchText, entryReplaceText, true, true) // Do replace in files
	// }
	if _, _, entryReplaceText, err := getArguments(); err == nil {
		for _, fileFound := range filesFoundMulti {
			if !fileFound.NotTextFile {

				fileFound.SearchAndRepl.ReplaceWith = entryReplaceText
				fileFound.SearchAndRepl.DoBackup = obj.findWinChkBackUp.GetActive()
				err = fileFound.SearchAndRepl.ReplaceInFile()
			}
			occReplaced += fileFound.SearchAndRepl.OccReplaced
			occurences += fileFound.SearchAndRepl.Occurrences

			fmt.Println(string(fileFound.SearchAndRepl.TextBytes))
		}
	}

	// TODO adjust occurances modified !!!

	if DlgErr(sts["unexpected"], err) {
		return
	}
	genericHideWindow(obj.findWin) // Hide find window
	genericHideWindow(obj.textWin) // Hide preview text window

	DialogMessage(obj.mainWin, "info", sts["done"],
		fmt.Sprintf("\n\n%s: %d/%d\n", sts["totalModified"], occReplaced, occurences), "", "Ok")

	if occurences > 0 {
		// Reset Init single found structure to avoid always display found pattern(s)
		fileFoundSingle = SearchAndReplaceNew("", []byte{}, "", "")
	}
}

// btnShowClipboardClicked:
func btnShowClipboardClicked() {

	if obj.entrySearch.GetTextLength() != 0 {

		textWinTextToShowBytes = []byte(clipboardGet())
		showTextWin(string(textWinTextToShowBytes), sts["clpbrdPreview"])
	}
}

// btnReplaceInClipboardClicked:
func btnReplaceInClipboardClicked() {

	if len(GetEntryText(obj.entrySearch)) == 0 {

		DialogMessage(obj.mainWin,
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
	genericHideWindow(obj.findWin)
	genericHideWindow(obj.textWin)
}

// Signal handler clicked ...
func textWinBtnDoneClicked() {
	obj.textWinChkShowModifications.SetActive(false)

	opt.SourceWinWidth, opt.SourceWinHeight = obj.textWin.GetSize()
	opt.SourceWinPosX, opt.SourceWinPosY = obj.textWin.GetPosition()
	opt.PanedWidth = opt.SourceWinWidth - obj.Paned.GetPosition()

	genericHideWindow(obj.textWin)

	currentText = ""
	currentLine = -1
}

// ViewButtonPressEvent:
func ViewButtonPressEvent(sv *source.SourceView, event *gdk.Event) {
	bEvent := gdk.EventButtonNewFromEvent(event)
	if bEvent.Button() == gdk.BUTTON_PRIMARY {
		currentLine = svs.GetCurrentLineNb()
	}
}

// SourceToggleButtonMapWidthToggled:
func SourceToggleButtonMapWidthToggled() {
	opt.PanedWidth = opt.SourceWinWidth - obj.Paned.GetPosition()
}

// WindowSourceCheckResize:
func WindowSourceCheckResize(w *gtk.Window) {

	if obj.SourceToggleButtonMapWidth.GetActive() {
		opt.SourceWinWidth, opt.SourceWinHeight = w.GetSize()
		obj.Paned.SetPosition(opt.SourceWinWidth - opt.PanedWidth)
	}
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
		if obj.textWin.GetVisible() {
			listViewFilesRowActivated(obj.listViewFiles)
		}
		// obj.textWin.SetVisible(true)
	}
}

// imgTop handler release signal (Display about box)
func imgTopReleaseEvent() {
	About.Width = 400
	About.ImageOkButtonSize = 24
	About.Show()
}

// Signal handler on Exit window ... Saving options before quit
func mainWinOnExit() {
	var err error
	opt.UpdateOptions()
	if err = opt.Write(); err != nil {
		log.Fatalf("Error writing option file: %s\n", err.Error())
	}
}

// OptionButtonDoneClicked
func OptionButtonDoneClicked() {
	genericHideWindow(obj.OptionsWindow)
}

// OptionsEntryFileSizeFocusOutEvent:
func OptionsEntryFileSizeChanged(e *gtk.Entry) {
	var err error
	var entry string
	if entry, err = e.GetText(); err == nil && len(entry) != 0 {
		entry = strings.TrimSpace(entry)
		if _, err = strconv.Atoi(entry); err == nil {
			// Value are sanitized, so, we only need to retrieve them
			opt.FileMaxSizeLimit = int64(GetEntryTextAsInt(obj.OptionsEntryMaxFileSize))
			opt.FileMinSizeLimit = int64(GetEntryTextAsInt(obj.OptionsEntryMinFileSize))
			obj.OptionButtonDone.SetSensitive(true)
		}
	}
	if err != nil {
		DialogMessage(obj.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
	}
	e.SetText(RemoveNonNum(entry))
	return
}

// MainButtonOptionsClicked:
func MainButtonOptionsClicked() {
	obj.OptionsWindow.ShowAll()
	obj.OptionsWindow.SetKeepAbove(true)
}
