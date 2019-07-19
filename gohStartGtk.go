// gohStartGtk.go

// Source file auto-generated on Fri, 19 Jul 2019 03:46:10 using Gotk3ObjHandler v1.3 ©2019 H.F.M

/*
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
		// Init tempDir and Remove tempDirectory on exit
		tempDir = tempMake(Name)
		defer os.RemoveAll(tempDir)
		// Parse Gtk objects
		gladeObjParser()
		// Objects Signals initialisations
		signalsPropHandler()
		// Fill control with images
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
		//	Start Gui loop
		gtk.Main()
	} else {
		log.Fatal("Builder initialisation error.")
	}
}
