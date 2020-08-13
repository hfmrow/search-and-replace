// main.go

/*
	Source file auto-generated on Thu, 06 Aug 2020 20:25:48 using Gotk3ObjHandler v1.5 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-20 H.F.M - Search And Replace v1.8 github.com/hfmrow/sAndReplace
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/*
	This software use Chroma:
	- Chroma — A general purpose syntax highlighter in pure Go, under the MIT License:
	  https://github.com/alecthomas/chroma/LICENSE
*/

import (
	"errors"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"

	glfsft "github.com/hfmrow/genLib/files/fileText"

	gimc "github.com/hfmrow/gtk3Import/misc"
	gitv "github.com/hfmrow/gtk3Import/textView"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

func main() {
	var err error

	/* Be or not to be ... in dev mode ... */
	devMode = false

	/* Build directory for tempDir */
	doTempDir = false

	/* Set to true when you choose using embedded assets functionality */
	assetsDeclarationsUseEmbedded(!devMode)

	/* Init Options */
	mainOptions = new(MainOpt)
	mainOptions.Init()

	/* Read Options */
	err = mainOptions.Read()
	if err != nil {
		fmt.Printf("%s\n%v\n", "Reading options error.", err)

	}

	/* Init gtk display */
	mainWindowTitle = fmt.Sprintf("%s %s  %s %s %s",
		Name,
		Vers,
		YearCreat,
		Creat,
		LicenseAbrv)

	mainStartGtk(mainWindowTitle,
		mainOptions.MainWinWidth,
		mainOptions.MainWinHeight, true)
}

func mainApplication() {
	var err error
	var tmpErr string

	/* Init AboutBox */
	mainOptions.AboutOptions.InitFillInfos(
		mainObjects.mainWin,
		"About "+Name,
		Name,
		Vers,
		Creat,
		YearCreat,
		LicenseAbrv,
		LicenseShort,
		Repository,
		Descr,
		"",
		tickIcon48)

	/* Specific UI tuning */
	chkCharacterClassToggled()

	/* Init Windows title for easy updating */
	mainWinTitle = gimc.TitleBarStructureNew(mainObjects.mainWin)
	findWinTitle = gimc.TitleBarStructureNew(mainObjects.findWin, sts["titleSearchResults"])
	textWinTitle = gimc.TitleBarStructureNew(mainObjects.textWin, sts["titlePreviewText"])

	/* Translate init. */
	translate = MainTranslateNew(absoluteRealPath+mainOptions.LanguageFilename, devMode)

	/*	Init Statusbar	*/
	statusbar = gimc.StatusBarStructureNew(mainObjects.statusbar, []string{
		sts["sbFiles"], sts["sbFilesSel"], sts["scanTime"],
		sts["searchTime"], sts["dispTime"], sts["status"]})

	/* Init ListStore */
	if tvsList, err = gitw.TreeViewStructureNew(mainObjects.listViewFiles, true, false); err == nil {

		tvsList.AddColumns(mainOptions.listStoreColumns, false, true, true, true, false, true)

		// Define selection changed function .
		tvsList.SelectionChangedFunc = func() {
			if tvsList.Selection.CountSelectedRows() > 0 {
				updateStatusBar()
			}
		}

		if err = tvsList.StoreSetup(new(gtk.ListStore)); err == nil {
			// Assign sorted column
			tvsList.Columns[1].Column.SetSortColumnID(4) // (col 1 will be sorted using values of col 4)
			tvsList.Columns[2].Column.SetSortColumnID(5) // (col 2 will be sorted using values of col 5)

			/* Init treeView popup menu */
			initTreeViewPopupMenu()

			/* Init TreeStore (found results) */
			if tvsTree, err = gitw.TreeViewStructureNew(mainObjects.findWinTreeView, true, false); err == nil {
				tvsTree.AddColumns(mainOptions.treeStoreColumns, false, true, false, false, false, true)
				err = tvsTree.StoreSetup(new(gtk.TreeStore))
			}
		}
	}
	if err != nil {
		DlgErr("mainApplication:TreeViewStructureNew", err)
		return
	}

	/* Init Drag and drop */
	dnd = gimc.DragNDropNew(mainObjects.listViewFiles, &currentInFilesList,
		func() {
			// Callaback function on Drag and drop operations
			fromDnD = true
			mainObjects.switchFileChooserButton.SetActive(false)
			mainObjects.fileChooserBtn.SetFilename("None")
			updateTreeViewFilesDisplay()
		})

	/* Init Clipboard */
	clipboardInit()

	/* Retreive command line arguments */
	if len(os.Args) > 1 {
		var countedError int
		var flagError = true
		for _, file := range os.Args[1:] {
			if _, err = os.Stat(file); err == nil {
				currentInFilesList = append(currentInFilesList, file)
			} else {
				if countedError < 14 && flagError {
					tmpErr += err.Error() + "\n"
					countedError++
				} else {
					flagError = false
					tmpErr += "too many errors ...\n"
				}
			}
		}
		if err == nil {
			fromDnD = true
			mainObjects.switchFileChooserButton.SetActive(false)
			mainObjects.fileChooserBtn.SetFilename("None")
		} else {
			if len(tmpErr) != 0 {
				err = errors.New(tmpErr)
			}
			DlgErr("Retreive command line arguments", err)
			return
		}
	}

	/* TextView with line number init.*/
	textViewRowNumber, err = gitv.TextViewNumberedNew(mainObjects.BoxTextViewPreview.Container)
	DlgErr("mainApplication:TextViewNumberedNew", err)
	textViewRowNumber.Editable = false

	/* Handling "populate-popup" signal to add some personal entries */
	textViewRowNumber.TextView.Connect("populate-popup", popupTextViewPopulateMenu)

	/* Hide ProgressBar */
	displayProgressBar(false)

	/* Display files list */
	updateTreeViewFilesDisplay()

	/* Initialize syntax highlighter to get list of languages and styles */
	highlighter, _ = ChromaHighlightNew(nil)

	/*Init single found structure */
	fileFoundSingle = glfsft.SearchAndReplaceNew([]byte{}, "", "")

	/* Initialize comboboxEntry */
	ComboBoxTextFill(mainObjects.textWinComboBoxTextStyleChooser, highlighter.Styles)
	ComboBoxTextFill(mainObjects.textWinComboBoxLanguage, highlighter.Lexers)
	ComboBoxTextAddSetEntry(mainObjects.textWinComboBoxTextStyleChooser, mainOptions.SyntaxHighlightType)
	ComboBoxTextAddSetEntry(mainObjects.textWinComboBoxLanguage, mainOptions.SyntaxHighlightLanguage)

	/* Focus search entry */
	mainObjects.entrySearch.GrabFocus()

	// mainObjects.entrySearch.SetText(`(\b)(func|return|exit\(\)|if|else|then|switch|case)`)
	// mainObjects.chkRegex.SetActive(true)
}
