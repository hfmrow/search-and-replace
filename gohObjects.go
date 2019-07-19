// gohObjects.go

// Source file auto-generated on Fri, 19 Jul 2019 03:46:10 using Gotk3ObjHandler v1.3 ©2019 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"github.com/gotk3/gotk3/gtk"
)

// Control over all used objects from glade.
var mainObjects *MainControlsObj

/******************************/
/* Main structure Declaration */
/******************************/
type MainControlsObj struct {
	btnFind                     *gtk.Button
	btnReplaceInClipboard       *gtk.Button
	btnShowClipboard            *gtk.Button
	chkCaseSensitive            *gtk.CheckButton
	chkCharacterClass           *gtk.CheckButton
	chkCharacterClassStrictMode *gtk.CheckButton
	chkFollowSymlinkDir         *gtk.CheckButton
	chkRegex                    *gtk.CheckButton
	chkSubDir                   *gtk.CheckButton
	chkUseEscapeChar            *gtk.CheckButton
	chkWildcard                 *gtk.CheckButton
	chkWoleWord                 *gtk.CheckButton
	clipboard                   *gtk.Clipboard /*MANUAL*/
	entryExtMask                *gtk.Entry
	entryReplace1               *gtk.Entry
	entrySearch1                *gtk.Entry
	fileChooserBtn              *gtk.FileChooserButton
	findWin                     *gtk.Window
	findWinCancelBtn            *gtk.Button
	findWinChkBackUp            *gtk.CheckButton
	findWinGrid                 *gtk.Grid
	findWinGridBtn              *gtk.Grid
	findWinLabel                *gtk.Label
	findWinLabel1               *gtk.Label
	findWinLabel2               *gtk.Label
	findWinReplaceBtn           *gtk.Button
	findWinScrollWin            *gtk.ScrolledWindow
	findWinTreeView             *gtk.TreeView
	findWinTreeviewSelection    *gtk.TreeSelection
	grid1                       *gtk.Grid
	gridAction                  *gtk.Grid
	gridOptionsRxWc             *gtk.Grid
	gridOptSearch               *gtk.Grid
	imgTop                      *gtk.Image
	ImgTopEventbox              *gtk.EventBox
	lblExtMask                  *gtk.Label
	lblReplace                  *gtk.Label
	lblSearch                   *gtk.Label
	listStore                   *gtk.ListStore /*MANUAL*/
	mainBox                     *gtk.Box
	mainGrid                    *gtk.Grid
	mainUiBuilder               *gtk.Builder
	mainWin                     *gtk.Window
	posixGrid                   *gtk.Grid
	replaceGrid                 *gtk.Grid
	scrolledWindowTreeView      *gtk.ScrolledWindow
	statusbar                   *gtk.Statusbar
	SwitchFileChooserButton     *gtk.Switch
	textWin                     *gtk.Window
	textWinBox                  *gtk.Box
	textWinBtnDone              *gtk.Button
	textWinChkShowModifications *gtk.CheckButton
	textWinChkWrap              *gtk.CheckButton
	textWinGridBottom           *gtk.Grid
	textWinGridTop              *gtk.Grid
	textWinScrolledwindow       *gtk.ScrolledWindow
	textWinTextview             *gtk.TextView
	treeStore                   *gtk.TreeStore /*MANUAL*/
	treeviewFiles               *gtk.TreeView
	treeviewSelection           *gtk.TreeSelection
}

/******************************/
/* GtkObjects  Initialisation */
/******************************/
// gladeObjParser: Initialise Gtk3 Objects into mainObjects structure.
func gladeObjParser() {
	mainObjects.btnFind = loadObject("btnFind").(*gtk.Button)
	mainObjects.btnReplaceInClipboard = loadObject("btnReplaceInClipboard").(*gtk.Button)
	mainObjects.btnShowClipboard = loadObject("btnShowClipboard").(*gtk.Button)
	mainObjects.chkCaseSensitive = loadObject("chkCaseSensitive").(*gtk.CheckButton)
	mainObjects.chkCharacterClass = loadObject("chkCharacterClass").(*gtk.CheckButton)
	mainObjects.chkCharacterClassStrictMode = loadObject("chkCharacterClassStrictMode").(*gtk.CheckButton)
	mainObjects.chkFollowSymlinkDir = loadObject("chkFollowSymlinkDir").(*gtk.CheckButton)
	mainObjects.chkRegex = loadObject("chkRegex").(*gtk.CheckButton)
	mainObjects.chkSubDir = loadObject("chkSubDir").(*gtk.CheckButton)
	mainObjects.chkUseEscapeChar = loadObject("chkUseEscapeChar").(*gtk.CheckButton)
	mainObjects.chkWildcard = loadObject("chkWildcard").(*gtk.CheckButton)
	mainObjects.chkWoleWord = loadObject("chkWoleWord").(*gtk.CheckButton)
	mainObjects.entryExtMask = loadObject("entryExtMask").(*gtk.Entry)
	mainObjects.entryReplace1 = loadObject("entryReplace1").(*gtk.Entry)
	mainObjects.entrySearch1 = loadObject("entrySearch1").(*gtk.Entry)
	mainObjects.fileChooserBtn = loadObject("fileChooserBtn").(*gtk.FileChooserButton)
	mainObjects.findWin = loadObject("findWin").(*gtk.Window)
	mainObjects.findWin.SetTitle(mainOptions.TxtSearchResult)           /*MANUAL*/
	mainObjects.findWin.SetModal(true)                                  /*MANUAL*/
	mainObjects.findWin.SetKeepAbove(mainOptions.ChildWindowsKeepAbove) /*MANUAL*/
	mainObjects.findWin.SetSizeRequest(640, 480)                        /*MANUAL*/
	mainObjects.findWin.Connect("delete_event", genericHideWindow)      /*MANUAL*/
	mainObjects.findWinCancelBtn = loadObject("findWinCancelBtn").(*gtk.Button)
	mainObjects.findWinChkBackUp = loadObject("findWinChkBackUp").(*gtk.CheckButton)
	mainObjects.findWinGrid = loadObject("findWinGrid").(*gtk.Grid)
	mainObjects.findWinGridBtn = loadObject("findWinGridBtn").(*gtk.Grid)
	mainObjects.findWinLabel = loadObject("findWinLabel").(*gtk.Label)
	mainObjects.findWinLabel1 = loadObject("findWinLabel1").(*gtk.Label)
	mainObjects.findWinLabel2 = loadObject("findWinLabel2").(*gtk.Label)
	mainObjects.findWinReplaceBtn = loadObject("findWinReplaceBtn").(*gtk.Button)
	mainObjects.findWinScrollWin = loadObject("findWinScrollWin").(*gtk.ScrolledWindow)
	mainObjects.findWinTreeView = loadObject("findWinTreeView").(*gtk.TreeView)
	mainObjects.findWinTreeviewSelection = loadObject("findWinTreeviewSelection").(*gtk.TreeSelection)
	mainObjects.grid1 = loadObject("grid1").(*gtk.Grid)
	mainObjects.gridAction = loadObject("gridAction").(*gtk.Grid)
	mainObjects.gridOptionsRxWc = loadObject("gridOptionsRxWc").(*gtk.Grid)
	mainObjects.gridOptSearch = loadObject("gridOptSearch").(*gtk.Grid)
	mainObjects.imgTop = loadObject("imgTop").(*gtk.Image)
	mainObjects.ImgTopEventbox = loadObject("ImgTopEventbox").(*gtk.EventBox)
	mainObjects.lblExtMask = loadObject("lblExtMask").(*gtk.Label)
	mainObjects.lblReplace = loadObject("lblReplace").(*gtk.Label)
	mainObjects.lblSearch = loadObject("lblSearch").(*gtk.Label)
	mainObjects.mainBox = loadObject("mainBox").(*gtk.Box)
	mainObjects.mainGrid = loadObject("mainGrid").(*gtk.Grid)
	mainObjects.mainWin = loadObject("mainWin").(*gtk.Window)
	mainObjects.posixGrid = loadObject("posixGrid").(*gtk.Grid)
	mainObjects.replaceGrid = loadObject("replaceGrid").(*gtk.Grid)
	mainObjects.scrolledWindowTreeView = loadObject("scrolledWindowTreeView").(*gtk.ScrolledWindow)
	mainObjects.statusbar = loadObject("statusbar").(*gtk.Statusbar)
	mainObjects.SwitchFileChooserButton = loadObject("SwitchFileChooserButton").(*gtk.Switch)
	mainObjects.textWin = loadObject("textWin").(*gtk.Window)
	mainObjects.textWin.SetTitle(mainOptions.TxtTextWin)                /*MANUAL*/
	mainObjects.textWin.SetModal(false)                                 /*MANUAL*/
	mainObjects.textWin.SetKeepAbove(mainOptions.ChildWindowsKeepAbove) /*MANUAL*/
	mainObjects.textWin.SetSizeRequest(640, 480)                        /*MANUAL*/
	mainObjects.textWin.Connect("delete_event", genericHideWindow)      /*MANUAL*/
	mainObjects.textWinBox = loadObject("textWinBox").(*gtk.Box)
	mainObjects.textWinBtnDone = loadObject("textWinBtnDone").(*gtk.Button)
	mainObjects.textWinChkShowModifications = loadObject("textWinChkShowModifications").(*gtk.CheckButton)
	mainObjects.textWinChkWrap = loadObject("textWinChkWrap").(*gtk.CheckButton)
	mainObjects.textWinGridBottom = loadObject("textWinGridBottom").(*gtk.Grid)
	mainObjects.textWinGridTop = loadObject("textWinGridTop").(*gtk.Grid)
	mainObjects.textWinScrolledwindow = loadObject("textWinScrolledwindow").(*gtk.ScrolledWindow)
	mainObjects.textWinTextview = loadObject("textWinTextview").(*gtk.TextView)
	mainObjects.treeviewFiles = loadObject("treeviewFiles").(*gtk.TreeView)
	mainObjects.treeviewSelection = loadObject("treeviewSelection").(*gtk.TreeSelection)
}
