// main.go

// Source file auto-generated on Sun, 15 Sep 2019 08:30:03 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gitv "github.com/hfmrow/gtk3Import/textView"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

func main() {
	var err error

	/* Be or not to be ... in dev mode ... */
	devMode = true

	/* Build directory for tempDir */
	doTempDir = false

	/* Naming widgets as Gtk objects names to use in css.
	   Set to false if they already named in Glade.*/
	namingWidget = true

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

	/* Init AboutBox */
	mainOptions.AboutOptions.InitFillInfos(
		"About "+Name,
		Name,
		Vers,
		Creat,
		YearCreat,
		LicenseAbrv,
		LicenseShort,
		Repository,
		Descr,
		sanderSearchAndReplace400x27,
		signSelect20)

	/* Init gtk display */
	mainWindowTitle = fmt.Sprintf("%s %s  %s %s %s",
		mainOptions.AboutOptions.AppName,
		mainOptions.AboutOptions.AppVers,
		mainOptions.AboutOptions.YearCreat,
		mainOptions.AboutOptions.AppCreats,
		mainOptions.AboutOptions.LicenseAbrv)

	mainStartGtk(mainWindowTitle,
		mainOptions.MainWinWidth,
		mainOptions.MainWinHeight, true)
}

func mainApplication() {
	var err error
	var tmpErr string

	/* Specific UI tuning */
	mainObjects.btnScan.SetVisible(mainOptions.ScanButtonEnabled)
	mainObjects.switchFileChooserButton.SetVisible(mainOptions.ScanButtonEnabled)
	chkCharacterClassToggled()

	/* Init Windows title fo easy updating */
	title, _ := mainObjects.mainWin.GetTitle()
	mainWinTitle = gimc.TitleBarNew(mainObjects.mainWin, title)
	findWinTitle = gimc.TitleBarNew(mainObjects.findWin, sts["titleSearchResults"])
	textWinTitle = gimc.TitleBarNew(mainObjects.textWin, sts["titlePreviewText"])

	/* Translate init. */
	translate = MainTranslateNew(absoluteRealPath+mainOptions.LanguageFilename, devMode)

	/*	Init Statusbar	*/

	statusbar = new(gimc.StatusBar)
	statusbar.StructureSetup(mainObjects.statusbar, []string{
		sts["sbFiles"], sts["sbFilesSel"], sts["scanTime"],
		sts["searchTime"], sts["dispTime"], sts["status"]})

	/* Init ListStore */
	if tvsList, err = gitw.TreeViewStructureNew(mainObjects.listViewFiles, true, false); err == nil {
		for _, col := range mainOptions.listStoreColumns {
			tvsList.AddColumn(col[0], col[1], false, true, true, true, false, true)
		}
		// Define selection changed function .
		tvsList.SelectionChangedFunc = func() {
			if tvsList.Selection.CountSelectedRows() > 0 {
				// glibList =  tvsList.Selection.GetSelectedRows(mainOptions.currentTreeview.ListStore)
				// data := glibList.Data()
				// path := data.(*gtk.TreePath)
				// rowNb, _ := strconv.Atoi(path.String())
				updateStatusBar()
			}
		}

		if err = tvsList.StoreSetup(new(gtk.ListStore)); err == nil {
			/* Init TreeStore */
			if tvsTree, err = gitw.TreeViewStructureNew(mainObjects.findWinTreeView, true, false); err == nil {
				for _, col := range mainOptions.treeStoreColumns {
					tvsTree.AddColumn(col[0], col[1], false, true, false, false, false, true)
				}
				err = tvsTree.StoreSetup(new(gtk.TreeStore))
			}
		}
	}

	if err != nil {
		gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
		return
	}

	/* Init Drag and drop */
	dnd = gimc.DragNDropNew(mainObjects.listViewFiles, &currentInFilesList,
		func() { // Callaback function on Drag and drop operations

			fromDnD = true
			mainObjects.switchFileChooserButton.SetActive(false)
			mainObjects.fileChooserBtn.SetFilename("None")
			updateTreeViewFilesDisplay()

			// TODO /* DEBUG */
			// for idx, file := range currentInFilesList {
			// 	fmt.Printf("%d - %s\n", idx, file)
			// }
			/* DEBUG */

			// if err != nil {
			// 	gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
			// }
		})

	/* Init Clipboard */
	clipboardInit()

	/* Retreive command line arguments */
	switch {
	case len(os.Args) == 1:
		// entryExtMaskFocusOut()
	case len(os.Args) > 1:
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
			gidg.DialogMessage(mainObjects.mainWin, "error", sts["alert"], "\n\n"+err.Error(), "", "Ok")
			return
		}
	}

	// TextView with line number init.
	if mainOptions.textViewRowNumber, err = gitv.TextViewRowNumberNew(mainObjects.textWinTextviewNumbers, mainObjects.textWinTextview,
		mainObjects.textWinScrolledwindowNumbers, mainObjects.textWinScrolledwindow, false); err != nil {
		log.Fatalf("Could not initialize preview display TextView: %s", err.Error())
	}

	/* Css Init */
	cssStyle := ` #textWinTextview * {
	color: shade (#332211, 1.06);
	background-color: #fefefe;	
	opacity: 0.99;
}

 #textWinTextview text selection {
	background-color: #aaddff;
	color:shade (#332211, 1.06);
}

 #textWinTextviewNumbers * {
	color: shade (#0033ff, 1.06);
	background-color: #eeeeee;	
	opacity: 0.99;
}

 #textWinTextviewNumbers text {
	color: shade (#0022ff, 1.06);
	background-color: #eeeeee;	
	opacity: 0.99;
}
`
	/* Applying cssStyle */
	gimc.CssWdgScnLoad(cssStyle)

	/* Display files list */
	updateTreeViewFilesDisplay()
}
