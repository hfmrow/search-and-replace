// gohStartGtk.go

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
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

/******************************/
/* Gtk3 Window Initialisation */
/******************************/
func mainStartGtk(winTitle string, width, height int, center bool) {
	obj = new(MainControlsObj)
	gtk.Init(nil)
	if newBuilder(mainGlade) == nil {
		// Init tempDir and Remove it on quit if requested.
		if doTempDir {
			tempDir = tempMake(Name)
			defer os.RemoveAll(tempDir)
		}
		// Parse Gtk objects
		gladeObjParser()
		/* Update gtk conctrols with stored values into opt */
		opt.UpdateObjects()
		/* Fill control with images */
		assignImages()
		// Set Window Properties
		if center {
			obj.mainWin.SetPosition(gtk.WIN_POS_CENTER)
		}
		obj.mainWin.SetTitle(winTitle)
		obj.mainWin.SetDefaultSize(width, height)
		obj.mainWin.Connect("delete-event", windowDestroy)
		obj.mainWin.ShowAll()
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

// windowDestroy: on closing/destroying the gui window.
func windowDestroy() {
	if onShutdown() {
		gtk.MainQuit()
	}
}
