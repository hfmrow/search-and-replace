// main.go

/*
	Source file auto-generated on Fri, 02 Apr 2021 10:53:33 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 H.F.M - Search And Replace v1.9 github.com/hfmrow/search-and-replace
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"

	glfsft "github.com/hfmrow/genLib/files/fileText"

	gimc "github.com/hfmrow/gtk3Import/misc"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

func main() {

	/* Build options */
	// devMode: is used in some functions to control the behavior of the program
	// When software is ready to be published, this flag must be set at "false"
	// that means:
	// - options file will be stored in $HOME/.config/[Creat]/[softwareName],
	// - translate function if used, will no more auto-update "sts" map sentences,
	// - all built-in assets will be used instead of the files themselves.
	//   Be aware to update assets via "Goh" and translations with "Got" before all.
	devMode = true
	absoluteRealPath, optFilename = getAbsRealPath()

	// Initialization of assets according to the chosen mode (devMode).
	// you can set this flag to your liking without reference to devMode.
	assetsDeclarationsUseEmbedded(!devMode)

	// Create temp directory .. or not
	doTempDir = false

	/* Init & read options file */
	mainOptions = new(MainOpt) // Assignate options' structure.
	mainOptions.Read()         // Read values from options' file if exists.

	/* Init gtk display */
	mainWindowTitle = fmt.Sprintf("%s %s  %s %s %s",
		Name,
		Vers,
		"©"+YearCreat,
		Creat,
		LicenseAbrv)

	mainStartGtk(mainWindowTitle,
		mainOptions.MainWinWidth,
		mainOptions.MainWinHeight, true)
}

func mainApplication() {

	var (
		err    error
		tmpErr string
	)

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
	svs, err = SourceViewStructNew(mainObjects.View, mainObjects.Map, mainObjects.textWin)
	DlgErr("mainApplication:SourceViewStructNew", err)

	svs.View.SetEditable(false)
	svs.View.SetHighlightCurrentLine(true)
	svs.Buffer.SetHighlightMatchingBrackets(true)

	/* Handling "populate-popup" signal to add some personal entries */
	svs.View.Connect("populate-popup", popupTextViewPopulateMenu)

	/* Hide ProgressBar */
	displayProgressBar(false)

	/* Display files list */
	updateTreeViewFilesDisplay()

	/*Init single found structure */
	fileFoundSingle = glfsft.SearchAndReplaceNew([]byte{}, "", "")

	/* Initialize comboboxEntry */
	svs.UserLanguagePath = filepath.Join(absoluteRealPath, mainOptions.HighlightUserDefined)
	svs.UserStylePath = filepath.Join(absoluteRealPath, mainOptions.HighlightUserDefined)
	svs.DefaultLanguageId = mainOptions.DefaultSourceLang
	svs.DefaultStyleShemeId = mainOptions.DefaultSourceStyle

	svs.ComboboxHandling(
		mainObjects.textWinComboBoxLanguage,
		mainObjects.textWinComboBoxTextStyleChooser,
		&mainOptions.DefaultSourceLang,
		&mainOptions.DefaultSourceStyle)

	/* Focus search entry */
	mainObjects.entrySearch.GrabFocus()
}

/*************************************\
/* Executed just before closing all. */
/************************************/
func onShutdown() bool {
	var err error
	// Update mainOptions with GtkObjects and save it
	if err = mainOptions.Write(); err == nil {
		// What you want to execute before closing the app.
		// Return:
		// true for exit applicaton
		// false does not exit application
	}
	if err != nil {
		log.Fatalf("Unexpected error on exit: %s", err.Error())
	}
	return true
}
