// gohStartGtk.go

// Source file auto-generated on Sun, 15 Sep 2019 08:30:03 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
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
		// Naming widgets as Gtk objects
		if namingWidget {
			widgetNaming()
		}
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
		if devMode {
			fmt.Printf("** %s %s **\n", "Development mode activated ...", "Assets used directly from the sources files.")
		}
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
