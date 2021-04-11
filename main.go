// main.go

/*
	Source file auto-generated on Fri, 09 Apr 2021 03:01:52 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 hfmrow - Search And Replace v1.10 github.com/hfmrow/search-and-replace
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

	/* Logger init. */
	Logger = Log2FileStructNew(optFilename, devMode)
	defer Logger.CloseLogger()

	// Initialization of assets according to the chosen mode (devMode).
	// you can set this flag to your liking without reference to devMode.
	assetsDeclarationsUseEmbedded(!devMode)

	// Create temp directory .. or not
	doTempDir = false

	/* Init & read options file */
	opt = new(MainOpt) // Assignate options' structure.
	opt.Read()         // Read values from options' file if exists.

	/* Init gtk display */
	mainWindowTitle = fmt.Sprintf("%s %s  %s %s %s",
		Name,
		Vers,
		"©"+YearCreat,
		Creat,
		LicenseAbrv)

	mainStartGtk(mainWindowTitle,
		opt.MainWinWidth,
		opt.MainWinHeight, true)
}

func mainApplication() {

	var (
		err    error
		tmpErr string
	)

	/* Init AboutBox */
	opt.AboutOptions.InitFillInfos(
		obj.mainWin,
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
	mainWinTitle = TitleBarStructureNew(obj.mainWin)
	findWinTitle = TitleBarStructureNew(obj.findWin, sts["titleSearchResults"])
	textWinTitle = TitleBarStructureNew(obj.textWin, sts["titlePreviewText"])

	/* Translate init. */
	translate = MainTranslateNew(absoluteRealPath+opt.LanguageFilename, devMode)

	/*	Init Statusbar	*/
	statusbar = StatusBarStructureNew(obj.statusbar, []string{
		sts["sbFiles"], sts["sbFilesSel"], sts["scanTime"],
		sts["searchTime"], sts["dispTime"], sts["status"]})

	/* Init ListStore */
	if tvsList, err = TreeViewStructureNew(obj.listViewFiles, true, false); err == nil {

		tvsList.AddColumns(opt.listStoreColumns, false, true, true, true, false, true)

		// Define selection changed function .
		tvsList.SelectionChangedFunc = func() {
			if tvsList.Selection.CountSelectedRows() > 0 {
				updateStatusBar()
			}
		}
		tvsList.Columns[opt.mapListStore["pathReal"]].Visible = false // Hold filename not midified, hidden

		if err = tvsList.StoreSetup(new(gtk.ListStore)); err == nil {
			// Assign sorted column
			tvsList.Columns[opt.mapListStore["Size"]].Column.SetSortColumnID(opt.mapListStore["sizeSort"]) // (col 1 will be sorted using values of col 4)
			tvsList.Columns[opt.mapListStore["Date"]].Column.SetSortColumnID(opt.mapListStore["dateSort"]) // (col 2 will be sorted using values of col 5)
			/* Init treeView popup menu */
			initTreeViewPopupMenu()

			/* Init TreeStore (found results) */
			if tvsTree, err = TreeViewStructureNew(obj.findWinTreeView, true, false); err == nil {
				tvsTree.AddColumns(opt.treeStoreColumns, false, true, false, false, false, true)
				tvsTree.Columns[opt.mapTreeStore["Toggle"]].Editable = true
				tvsTree.Columns[opt.mapTreeStore["Toggle"]].ReadOnly = false
				tvsTree.Columns[opt.mapTreeStore["fileIdx"]].Visible = false // File index storage
				tvsTree.Columns[opt.mapTreeStore["lineIdx"]].Visible = false // Line Index storage

				tvsTree.CallbackOnSetColValue = func(iter *gtk.TreeIter, col int, value interface{}) {

					if col == opt.mapTreeStore["Toggle"] {
						fileIdx := int(tvsTree.GetColValue(iter, opt.mapTreeStore["fileIdx"]).(int64))
						lineIdx := int(tvsTree.GetColValue(iter, opt.mapTreeStore["lineIdx"]).(int64))

						if lineIdx > -1 {
							idxExist := IsExistSlIface(
								&filesFoundMulti[fileIdx].SearchAndRepl.Pos.UntouchedLines,
								filesFoundMulti[fileIdx].SearchAndRepl.Pos.FoundLinesIdx[lineIdx].Number)

							if !value.(bool) && idxExist == -1 {
								filesFoundMulti[fileIdx].SearchAndRepl.Pos.UntouchedLines = append(
									filesFoundMulti[fileIdx].SearchAndRepl.Pos.UntouchedLines,
									filesFoundMulti[fileIdx].SearchAndRepl.Pos.FoundLinesIdx[lineIdx].Number)

							} else if idxExist > -1 {
								for idxExist > -1 {

									DeleteSlIface(&filesFoundMulti[fileIdx].SearchAndRepl.Pos.UntouchedLines, idxExist)
									idxExist = IsExistSlIface(
										&filesFoundMulti[fileIdx].SearchAndRepl.Pos.UntouchedLines,
										filesFoundMulti[fileIdx].SearchAndRepl.Pos.FoundLinesIdx[lineIdx].Number)
								}
							}
						}
					}
				}
			}
			err = tvsTree.StoreSetup(new(gtk.TreeStore))
		}
	}
	if err != nil {
		DlgErr("mainApplication:TreeViewStructureNew", err)
		return
	}

	/* Init Drag and drop */
	dnd = DragNDropNew(obj.listViewFiles, &currentInFilesList,
		func() {
			// Callaback function on Drag and drop operations
			fromDnD = true
			obj.switchFileChooserButton.SetActive(false)
			if len(*dnd.FilesList) > 0 {
				l := *dnd.FilesList
				dir := filepath.Dir(l[0])
				obj.fileChooserBtn.SetFilename(dir)
				updateTreeViewFilesDisplay()
			}
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
			obj.switchFileChooserButton.SetActive(false)
			obj.fileChooserBtn.SetFilename("None")
		} else {
			if len(tmpErr) != 0 {
				err = errors.New(tmpErr)
			}
			DlgErr("Retreive command line arguments", err)
			return
		}
	}

	/* TextView with line number init.*/
	svs, err = SourceViewStructNew(obj.View, obj.Map, obj.textWin)
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
	fileFoundSingle = SearchAndReplaceNew("", []byte{}, "", "")

	/* Initialize comboboxEntry */
	svs.UserLanguagePath = filepath.Join(absoluteRealPath, opt.HighlightUserDefined)
	svs.UserStylePath = filepath.Join(absoluteRealPath, opt.HighlightUserDefined)
	svs.DefaultLanguageId = opt.DefaultSourceLang
	svs.DefaultStyleShemeId = opt.DefaultSourceStyle

	svs.ComboboxHandling(
		obj.textWinComboBoxLanguage,
		obj.textWinComboBoxTextStyleChooser,
		&opt.DefaultSourceLang,
		&opt.DefaultSourceStyle)

	/* Focus search entry */
	obj.entrySearch.GrabFocus()

	/* Progressbar Init */
	pbs = ProgressGifNew(linearProgressHorzBlue, obj.mainBox, 1)
}

/*************************************\
/* Executed just before closing all. */
/************************************/
func onShutdown() bool {
	var err error
	// Update opt with GtkObjects and save it
	if err = opt.Write(); err == nil {
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
