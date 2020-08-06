// gohStartGtk.go

/*
	Source file auto-generated on Thu, 06 Aug 2020 20:25:48 using Gotk3ObjHandler v1.5 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-20 H.F.M - Search And Replace v1.8 github.com/hfmrow/sAndReplace
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

/******************************/
/* Gtk3 Window Initialisation */
/******************************/
func mainStartGtk(winTitle string, width, height int, center bool) {
	mainObjects = new(MainControlsObj)
	gtk.Init(nil)
	if newBuilder(mainGlade) == nil {
		// Init tempDir and Remove it on quit if requested.
		if doTempDir {
			tempDir = tempMake(Name)
			defer os.RemoveAll(tempDir)
		}
		// Parse Gtk objects
		gladeObjParser()
		/* Update gtk conctrols with stored values into mainOptions */
		mainOptions.UpdateObjects()
		/* Fill control with images */
		assignImages()
		// Set Window Properties
		if center {
			mainObjects.mainWin.SetPosition(gtk.WIN_POS_CENTER)
		}
		mainObjects.mainWin.SetTitle(winTitle)
		mainObjects.mainWin.SetDefaultSize(width, height)
		mainObjects.mainWin.Connect("delete-event", windowDestroy)
		mainObjects.mainWin.ShowAll()
		// Start main application ...
		mainApplication()
		// Objects Signals initialisations
		signalsPropHandler()
		//	Start Gui loop
		gtk.Main()
	} else {
		log.Fatal("Builder initialisation error.")
	}
}
