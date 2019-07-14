// main.go

// Source file auto-generated on Sun, 14 Jul 2019 16:40:30 using Gotk3ObjHandler v1.3 ©2019 H.F.M

/*
	SearchAndReplace v1.7.2 ©2018 H.F.M

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:

	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
	associated documentation files (the "Software"), to dealin the Software without restriction,
	including without limitation the rights to use, copy, modify, merge, publish, distribute,
	sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all copies or
	substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
	NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
	NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
	DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT
	OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"errors"
	"fmt"
	"os"

	gi "github.com/hfmrow/sAndReplace/gtk3Import"
)

func main() {
	devMode = true

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
		checked18x18)

	/* Init gtk display */
	mainStartGtk(fmt.Sprintf("%s %s  %s %s %s",
		mainOptions.AboutOptions.AppName,
		mainOptions.AboutOptions.AppVers,
		mainOptions.AboutOptions.YearCreat,
		mainOptions.AboutOptions.AppCreats,
		mainOptions.AboutOptions.LicenseAbrv),
		mainOptions.MainWinWidth,
		mainOptions.MainWinHeight, true)
}

func mainApplication() {
	var err error
	var tmpErr string
	var fileInfo os.FileInfo

	/* Translate init. */
	translate = MainTranslateNew(absoluteRealPath+mainOptions.LanguageFilename, devMode)
	sts = translate.Sentences

	/* Update gtk conctrols with stored values into mainOptions */
	mainOptions.UpdateObjects()

	/* Init listStore & treeStore */
	mainObjects.listStore = gi.TreeViewListStoreSetup(mainObjects.treeviewFiles, true, mainOptions.TreeViewColumns)
	mainObjects.treeStore = gi.TreeViewTreeStoreSetup(mainObjects.findWinTreeView, true, mainOptions.TreeStoreColumns)

	/* Init Clipboard*/
	clipboardInit()

	/* Retreive command line arguments */
	switch {
	// case len(os.Args) == 1:
	// 	if mainOptions.Directory == "" {
	// 		mainOptions.Directory = filepath.Dir(os.Args[0])
	// 	}
	// 	scanFilesAndDisp()
	default:
		if len(os.Args) >= 2 {
			if fileInfo, err = os.Stat(os.Args[1]); err == nil {
				if fileInfo.IsDir() {
					mainOptions.Directory = os.Args[1]
					scanFilesAndDisp()
				} else {
					mainOptions.currentInFilesList = mainOptions.currentInFilesList[:0]
					for _, file := range os.Args {
						if _, err := os.Stat(os.Args[1]); err == nil {
							mainOptions.currentInFilesList = append(mainOptions.currentInFilesList, file)
						} else {
							tmpErr += err.Error() + "\n"
						}
					}
					err = getFilesSelection(mainOptions.currentInFilesList[1:])
					if err != nil {
						tmpErr += err.Error() + "\n"
					}
				}
			}
		}
	}
	if len(tmpErr) != 0 {
		err = errors.New(tmpErr)
	}
	if err != nil {
		gi.DlgMessage(mainObjects.mainWin, "error", mainOptions.TxtAlert, err.Error(), "", "Ok")
	}
}
