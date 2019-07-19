// signalsHandlers.go

package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	g "github.com/hfmrow/sAndReplace/genLib"
	gi "github.com/hfmrow/sAndReplace/gtk3Import"
)

// SwitchFileChooserButtonStateSet:
func SwitchFileChooserButtonStateSet() {
	mainObjects.fileChooserBtn.SetSensitive(mainObjects.SwitchFileChooserButton.GetActive())
}

// treeviewFilesReceived:
func treeviewFilesReceived(tw *gtk.TreeView, context *gdk.DragContext, x, y int, data_ptr uintptr, info, time uint) {
	var err error
	mainOptions.currentInFilesList = mainOptions.currentInFilesList[:0]
	// Convert pointer to datas
	data := gtk.GetData(data_ptr)
	list := strings.Split(string(data), g.GetTextEOL(data))
	for _, file := range list {
		if len(file) != 0 {
			if u, err := url.PathUnescape(file); err == nil {
				mainOptions.currentInFilesList = append(mainOptions.currentInFilesList, strings.TrimPrefix(u, "file://"))
			}
		}
	}
	err = getFilesSelection(mainOptions.currentInFilesList)
	if err != nil {
		gi.DlgMessage(mainObjects.mainWin, "error", mainOptions.TxtAlert, err.Error(), "", "Ok")
	}
}

// Signal handler changed ... FocusOut
func entrySearchFocusOut(e *gtk.Entry) {
	genericEntryFocusOut(e)
}

// Signal handler changed ... FocusOut
func entryReplaceFocusOut(e *gtk.Entry) {
	genericEntryFocusOut(e)
}

// Signal handler changed ... FocusOut - to trim entry
func entryExtMaskFocusOut() bool {

	var dispError = func() {
		// idle run to permit gtk3 working right during goroutine
		_, err := glib.IdleAdd(func() {
			gi.DlgMessage(mainObjects.mainWin, "error", sts["file-rem"], err.Error(), "", "Ok")
		})
		Check(err)
	}

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

	// Check if directory exist ...
	if _, err = os.Stat(mainOptions.Directory); os.IsNotExist(err) && !cmdLineArg {
		gi.DlgMessage(mainObjects.mainWin, "error", mainOptions.TxtAlert, "\n"+sts["dir-rem"]+"\n\n"+err.Error(), "", "Ok")
		// mainOptions.Directory = filepath.Dir(os.Args[0])
		scanFilesAndDisp()
	} else {
		err = getFilesSelection(mainOptions.currentInFilesList)
		if err != nil {
			// Tips, calling goroutine allow to finish this function and
			// returning the "GDK_EVENT_PROPAGATE" signal to avoid GDK error
			go dispError()
		}
	}
	return false // GDK_EVENT_PROPAGATE signal
}

// Signal handler changed ...
func treeViewSelectionChanged(ts *gtk.TreeSelection) {
	// Returns glib.List of gtk.TreePath pointers
	rows := ts.GetSelectedRows(mainObjects.listStore)
	items := make([]string, 0, rows.Length())

	for l := rows; l != nil; l = l.Next() {
		path := l.Data().(*gtk.TreePath)
		iter, _ := mainObjects.listStore.GetIter(path)
		value, _ := mainObjects.listStore.GetValue(iter, 3) // Field 3: get full path
		str, _ := value.GetString()
		items = append(items, str)
	}
	treeviewSelectedRows = items
	filesSelected = len(treeviewSelectedRows)
	updateStatusBar()
}

// Signal handler dblClick ... doudble click on ound result to display text preview
func findTreeViewDblClick(tw *gtk.TreeView) {
	ts, _ := tw.GetSelection()
	rows := ts.GetSelectedRows(mainObjects.treeStore)
	path := rows.Data().(*gtk.TreePath)
	tIter, _ := mainObjects.treeStore.GetIter(path)
	value, _ := mainObjects.treeStore.GetValue(tIter, 0)
	str, _ := value.GetString()
	if g.FileExist(str) {
		text, err := g.ReadFile(str)
		if err != nil {
			gi.DlgMessage(mainObjects.mainWin, "error", mainOptions.TxtAlert, err.Error(), "", "Ok")
		} else {
			mainObjects.textWin.SetModal(true)
			showTextWin(string(text))
		}
	}
}

// Signal handler dblClick ...
func treeViewDblClick(tw *gtk.TreeView) {
	var err error
	ts, _ := tw.GetSelection()
	rows := ts.GetSelectedRows(mainObjects.listStore)
	path := rows.Data().(*gtk.TreePath)
	iter, _ := mainObjects.listStore.GetIter(path)
	value, _ := mainObjects.listStore.GetValue(iter, 3) // Field 3: get full path
	// Get selected file path
	filename, _ := value.GetString()

	// Check for text file
	isTxt, _ := g.IsTextFile(filename)
	if strings.Contains(isTxt, "utf-") {
		textWinTextToShowBytes, err = g.ReadFile(filename)
		if err != nil {
			gi.DlgMessage(mainObjects.mainWin, "error", mainOptions.TxtAlert, err.Error(), "", "Ok")
		} else {
			// Is there any words out there ...
			if len(textWinTextToShowBytes) != 0 {
				mainObjects.textWin.SetModal(false)
				showTextWin(string(textWinTextToShowBytes))
			}
		}
	}
}

// Signal handler toggled ... (Wrap text)
func textWinChkWrapToggled() {
	if mainObjects.textWinChkWrap.GetActive() {
		mainObjects.textWinTextview.SetWrapMode(gtk.WRAP_WORD)
	} else {
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

// Signal handler toggled ...
func textWinChkShowModificationsToggled() {
	buff, err := mainObjects.textWinTextview.GetBuffer()
	if err != nil {
		gi.DlgMessage(mainObjects.mainWin, "warning", mainOptions.TxtAlert,
			"\n\n"+g.FormatText(err.Error(), 69, true), "", "Ok")
		mainObjects.textWinChkShowModifications.SetActive(false)
		return // No action, back to previous state
	}
	if mainObjects.textWinChkShowModifications.GetActive() {
		entrySearchText := getEntryText(mainObjects.entrySearch1)
		if len(entrySearchText) == 0 {
			gi.DlgMessage(mainObjects.mainWin,
				"warning", mainOptions.TxtSomethingMissing,
				"\n\n"+mainOptions.TxtNothingToSearch, "", "Ok")
			mainObjects.textWinChkShowModifications.SetActive(false)
			return // No action, back to previous state
		} else {

			txt, _ := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), false)
			textWinTextToShowBytes = []byte(txt)

			var outText []byte
			err := onTheFlyReplace(textWinTextToShowBytes, &outText, entrySearchText)
			if err != nil {
				gi.DlgMessage(mainObjects.mainWin, "warning", mainOptions.TxtAlert,
					"\n\n"+g.FormatText(err.Error(), 69, true), "", "Ok")
				mainObjects.textWinChkShowModifications.SetActive(false)
				return // No action, back to previous state
			} else {
				buff.SetText(string(outText))
			}
		}
	} else {
		buff.SetText(string(textWinTextToShowBytes))
	}
}

// Signal handler clicked ...
func findWinReplaceBtnClicked() {
	applyChanges = true
	btnFindClicked() // Do replace in files
	applyChanges = false
	genericHideWindow(mainObjects.findWin) // Hide find window
	fileChooserBtnClicked()                // Reload files list
}

// Signal handler clicked ...
func btnFindClicked() {
	removeEmptyResult := true
	entrySearchText := getEntryText(mainObjects.entrySearch1)
	entryReplaceText := getEntryText(mainObjects.entryReplace1)

	if len(entrySearchText) != 0 {
		mainFound, err := Find(entrySearchText, entryReplaceText, removeEmptyResult)
		if err != nil {
			gi.DlgMessage(mainObjects.mainWin,
				"warning", mainOptions.TxtAlert,
				err.Error(), "", "Ok")
		} else {
			if len(mainFound) != 0 {
				if !mainObjects.findWin.GetVisible() {
					showResults(&mainFound)
					mainObjects.findWin.Show()
				}
			} else {
				gi.DlgMessage(mainObjects.mainWin,
					"warning", mainOptions.TxtSomethingMissing,
					"\n\n"+mainOptions.TxtNothingFound, "", "Ok")
			}
		}
	} else {
		gi.DlgMessage(mainObjects.mainWin,
			"warning", mainOptions.TxtSomethingMissing,
			"\n\n"+mainOptions.TxtNothingToSearch, "", "Ok")
	}
}

// Signal handler clicked ...
func btnShowClipboardClicked() {
	textWinTextToShowBytes = []byte(clipboardGet())
	showTextWin(string(textWinTextToShowBytes))
}

// Signal handler clicked ...
func btnReplaceInClipboardClicked() {
	entrySearchText, err := mainObjects.entrySearch1.GetText()
	if err != nil {
		gi.DlgMessage(mainObjects.mainWin, "error", mainOptions.TxtAlert,
			"\n\n"+g.FormatText(err.Error(), 69, true), "", "Ok")
	} else if len(entrySearchText) == 0 {
		gi.DlgMessage(mainObjects.mainWin,
			"warning", mainOptions.TxtSomethingMissing,
			"\n\n"+mainOptions.TxtNothingToSearch, "", "Ok")
	} else {
		outText := new([]byte)
		err := onTheFlyReplace([]byte(clipboardGet()), outText, entrySearchText)
		if err != nil {
			gi.DlgMessage(mainObjects.mainWin, "warning", mainOptions.TxtAlert,
				"\n\n"+g.FormatText(err.Error(), 69, true), "", "Ok")
		} else {
			clipboardSet(string(*outText))
		}
	}
}

// Signal handler clicked ...
func fileChooserBtnClicked() {
	mainOptions.Directory = mainObjects.fileChooserBtn.GetFilename()
	cmdLineArg = false
	scanFilesAndDisp()
}

// Signal handler clicked ...
func findWinCancelBtnClicked() {
	genericHideWindow(mainObjects.findWin)
}

// Signal handler clicked ...
func textWinBtnDoneClicked() {
	mainObjects.textWinChkShowModifications.SetActive(false)
	genericHideWindow(mainObjects.textWin)
}

// Signal handler delete_event (hidding window)
func genericHideWindow(w *gtk.Window) bool {
	if w.GetVisible() {
		w.Hide()
	}
	return true
}

// Used with gtk.Entry objects.
func genericEntryFocusOut(e *gtk.Entry) {
	entry, err := e.GetText()
	if err != nil {
		gi.DlgMessage(mainObjects.mainWin, "error", mainOptions.TxtAlert, "\n\n"+g.FormatText(err.Error(), 69, true), "", "Ok")
	}
	entry = strings.Replace(fmt.Sprintf("%q", entry), `"`, ``, -1)
	e.SetText(entry)
}

// imgTop handler release signal (Display about box)
func imgTopReleaseEvent() {
	mainOptions.AboutOptions.Show()
}

// Signal handler on Exit window ... Saving options before quit
func mainWinOnExit() {
	if !devMode {
		mainOptions.UpdateOptions()
		err = mainOptions.Write()
		g.CheckE(err, "WriteJson", optFilename)
	}
}
