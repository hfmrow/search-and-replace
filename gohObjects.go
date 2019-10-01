// gohObjects.go

// Source file auto-generated on Tue, 01 Oct 2019 16:36:13 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

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
	btnFind                      *gtk.Button
	btnReplaceInClipboard        *gtk.Button
	btnScan                      *gtk.Button
	btnShowClipboard             *gtk.Button
	chkCaseSensitive             *gtk.CheckButton
	chkCharacterClass            *gtk.CheckButton
	chkCharacterClassStrictMode  *gtk.CheckButton
	chkFollowSymlinkDir          *gtk.CheckButton
	chkRegex                     *gtk.CheckButton
	chkUseEscapeChar             *gtk.CheckButton
	chkWildcard                  *gtk.CheckButton
	chkWoleWord                  *gtk.CheckButton
	clipboard                    *gtk.Clipboard /*MANUAL*/
	entryExtMask                 *gtk.Entry
	entryReplace                 *gtk.Entry
	entrySearch                  *gtk.Entry
	fileChooserBtn               *gtk.FileChooserButton
	findWin                      *gtk.Window
	findWinCancelBtn             *gtk.Button
	findWinChkBackUp             *gtk.CheckButton
	findWinChkDispForbFiles      *gtk.CheckButton
	findWinGrid                  *gtk.Grid
	findWinGridBtn               *gtk.Grid
	findWinReplaceBtn            *gtk.Button
	findWinScrollWin             *gtk.ScrolledWindow
	findWinTreeView              *gtk.TreeView
	imgTop                       *gtk.Image
	ImgTopEventbox               *gtk.EventBox
	lblExtMask                   *gtk.Label
	lblReplace                   *gtk.Label
	lblSearch                    *gtk.Label
	listViewFiles                *gtk.TreeView
	mainBox                      *gtk.Box
	mainGrid                     *gtk.Grid
	mainUiBuilder                *gtk.Builder
	mainWin                      *gtk.Window
	mainWinBtnClose              *gtk.Button
	replaceGrid                  *gtk.Grid
	scrolledWindowTreeView       *gtk.ScrolledWindow
	spinButtonDepth              *gtk.SpinButton
	statusbar                    *gtk.Statusbar
	switchFileChooserButton      *gtk.Switch
	textWin                      *gtk.Window
	textWinBox                   *gtk.Box
	textWinBtnDone               *gtk.Button
	textWinChkShowModifications  *gtk.CheckButton
	textWinChkWrap               *gtk.CheckButton
	textWinGridTop               *gtk.Grid
	textWinScrolledwindow        *gtk.ScrolledWindow
	textWinScrolledwindowNumbers *gtk.ScrolledWindow
	textWinTextview              *gtk.TextView
	textWinTextviewNumbers       *gtk.TextView
	treeviewSelection            *gtk.TreeSelection
}

/******************************/
/* GtkObjects  Initialisation */
/******************************/
// gladeObjParser: Initialise Gtk3 Objects into mainObjects structure.
func gladeObjParser() {
	mainObjects.btnFind = loadObject("btnFind").(*gtk.Button)
	mainObjects.btnReplaceInClipboard = loadObject("btnReplaceInClipboard").(*gtk.Button)
	mainObjects.btnScan = loadObject("btnScan").(*gtk.Button)
	mainObjects.btnShowClipboard = loadObject("btnShowClipboard").(*gtk.Button)
	mainObjects.chkCaseSensitive = loadObject("chkCaseSensitive").(*gtk.CheckButton)
	mainObjects.chkCharacterClass = loadObject("chkCharacterClass").(*gtk.CheckButton)
	mainObjects.chkCharacterClassStrictMode = loadObject("chkCharacterClassStrictMode").(*gtk.CheckButton)
	mainObjects.chkFollowSymlinkDir = loadObject("chkFollowSymlinkDir").(*gtk.CheckButton)
	mainObjects.chkRegex = loadObject("chkRegex").(*gtk.CheckButton)
	mainObjects.chkUseEscapeChar = loadObject("chkUseEscapeChar").(*gtk.CheckButton)
	mainObjects.chkWildcard = loadObject("chkWildcard").(*gtk.CheckButton)
	mainObjects.chkWoleWord = loadObject("chkWoleWord").(*gtk.CheckButton)
	mainObjects.entryExtMask = loadObject("entryExtMask").(*gtk.Entry)
	mainObjects.entryReplace = loadObject("entryReplace").(*gtk.Entry)
	mainObjects.entrySearch = loadObject("entrySearch").(*gtk.Entry)
	mainObjects.fileChooserBtn = loadObject("fileChooserBtn").(*gtk.FileChooserButton)
	mainObjects.findWin = loadObject("findWin").(*gtk.Window)
	mainObjects.findWinCancelBtn = loadObject("findWinCancelBtn").(*gtk.Button)
	mainObjects.findWinChkBackUp = loadObject("findWinChkBackUp").(*gtk.CheckButton)
	mainObjects.findWinChkDispForbFiles = loadObject("findWinChkDispForbFiles").(*gtk.CheckButton)
	mainObjects.findWinGrid = loadObject("findWinGrid").(*gtk.Grid)
	mainObjects.findWinGridBtn = loadObject("findWinGridBtn").(*gtk.Grid)
	mainObjects.findWinReplaceBtn = loadObject("findWinReplaceBtn").(*gtk.Button)
	mainObjects.findWinScrollWin = loadObject("findWinScrollWin").(*gtk.ScrolledWindow)
	mainObjects.findWinTreeView = loadObject("findWinTreeView").(*gtk.TreeView)
	mainObjects.imgTop = loadObject("imgTop").(*gtk.Image)
	mainObjects.ImgTopEventbox = loadObject("ImgTopEventbox").(*gtk.EventBox)
	mainObjects.lblExtMask = loadObject("lblExtMask").(*gtk.Label)
	mainObjects.lblReplace = loadObject("lblReplace").(*gtk.Label)
	mainObjects.lblSearch = loadObject("lblSearch").(*gtk.Label)
	mainObjects.listViewFiles = loadObject("listViewFiles").(*gtk.TreeView)
	mainObjects.mainBox = loadObject("mainBox").(*gtk.Box)
	mainObjects.mainGrid = loadObject("mainGrid").(*gtk.Grid)
	mainObjects.mainWin = loadObject("mainWin").(*gtk.Window)
	mainObjects.mainWinBtnClose = loadObject("mainWinBtnClose").(*gtk.Button)
	mainObjects.replaceGrid = loadObject("replaceGrid").(*gtk.Grid)
	mainObjects.scrolledWindowTreeView = loadObject("scrolledWindowTreeView").(*gtk.ScrolledWindow)
	mainObjects.spinButtonDepth = loadObject("spinButtonDepth").(*gtk.SpinButton)
	mainObjects.statusbar = loadObject("statusbar").(*gtk.Statusbar)
	mainObjects.switchFileChooserButton = loadObject("switchFileChooserButton").(*gtk.Switch)
	mainObjects.textWin = loadObject("textWin").(*gtk.Window)
	mainObjects.textWinBox = loadObject("textWinBox").(*gtk.Box)
	mainObjects.textWinBtnDone = loadObject("textWinBtnDone").(*gtk.Button)
	mainObjects.textWinChkShowModifications = loadObject("textWinChkShowModifications").(*gtk.CheckButton)
	mainObjects.textWinChkWrap = loadObject("textWinChkWrap").(*gtk.CheckButton)
	mainObjects.textWinGridTop = loadObject("textWinGridTop").(*gtk.Grid)
	mainObjects.textWinScrolledwindow = loadObject("textWinScrolledwindow").(*gtk.ScrolledWindow)
	mainObjects.textWinScrolledwindowNumbers = loadObject("textWinScrolledwindowNumbers").(*gtk.ScrolledWindow)
	mainObjects.textWinTextview = loadObject("textWinTextview").(*gtk.TextView)
	mainObjects.textWinTextviewNumbers = loadObject("textWinTextviewNumbers").(*gtk.TextView)
	mainObjects.treeviewSelection = loadObject("treeviewSelection").(*gtk.TreeSelection)

}

/*************************************/
/* GtkObjects Widget naming. Usualy */
/* used in css to identify objects */
/**********************************/
func widgetNaming() {
	mainObjects.btnFind.SetName("btnFind")
	mainObjects.btnReplaceInClipboard.SetName("btnReplaceInClipboard")
	mainObjects.btnScan.SetName("btnScan")
	mainObjects.btnShowClipboard.SetName("btnShowClipboard")
	mainObjects.chkCaseSensitive.SetName("chkCaseSensitive")
	mainObjects.chkCharacterClass.SetName("chkCharacterClass")
	mainObjects.chkCharacterClassStrictMode.SetName("chkCharacterClassStrictMode")
	mainObjects.chkFollowSymlinkDir.SetName("chkFollowSymlinkDir")
	mainObjects.chkRegex.SetName("chkRegex")
	mainObjects.chkUseEscapeChar.SetName("chkUseEscapeChar")
	mainObjects.chkWildcard.SetName("chkWildcard")
	mainObjects.chkWoleWord.SetName("chkWoleWord")
	mainObjects.entryExtMask.SetName("entryExtMask")
	mainObjects.entryReplace.SetName("entryReplace")
	mainObjects.entrySearch.SetName("entrySearch")
	mainObjects.fileChooserBtn.SetName("fileChooserBtn")
	mainObjects.findWin.SetName("findWin")
	mainObjects.findWinCancelBtn.SetName("findWinCancelBtn")
	mainObjects.findWinChkBackUp.SetName("findWinChkBackUp")
	mainObjects.findWinChkDispForbFiles.SetName("findWinChkDispForbFiles")
	mainObjects.findWinGrid.SetName("findWinGrid")
	mainObjects.findWinGridBtn.SetName("findWinGridBtn")
	mainObjects.findWinReplaceBtn.SetName("findWinReplaceBtn")
	mainObjects.findWinScrollWin.SetName("findWinScrollWin")
	mainObjects.findWinTreeView.SetName("findWinTreeView")
	mainObjects.imgTop.SetName("imgTop")
	mainObjects.ImgTopEventbox.SetName("ImgTopEventbox")
	mainObjects.lblExtMask.SetName("lblExtMask")
	mainObjects.lblReplace.SetName("lblReplace")
	mainObjects.lblSearch.SetName("lblSearch")
	mainObjects.listViewFiles.SetName("listViewFiles")
	mainObjects.mainBox.SetName("mainBox")
	mainObjects.mainGrid.SetName("mainGrid")
	mainObjects.mainWin.SetName("mainWin")
	mainObjects.mainWinBtnClose.SetName("mainWinBtnClose")
	mainObjects.replaceGrid.SetName("replaceGrid")
	mainObjects.scrolledWindowTreeView.SetName("scrolledWindowTreeView")
	mainObjects.spinButtonDepth.SetName("spinButtonDepth")
	mainObjects.statusbar.SetName("statusbar")
	mainObjects.switchFileChooserButton.SetName("switchFileChooserButton")
	mainObjects.textWin.SetName("textWin")
	mainObjects.textWinBox.SetName("textWinBox")
	mainObjects.textWinBtnDone.SetName("textWinBtnDone")
	mainObjects.textWinChkShowModifications.SetName("textWinChkShowModifications")
	mainObjects.textWinChkWrap.SetName("textWinChkWrap")
	mainObjects.textWinGridTop.SetName("textWinGridTop")
	mainObjects.textWinScrolledwindow.SetName("textWinScrolledwindow")
	mainObjects.textWinScrolledwindowNumbers.SetName("textWinScrolledwindowNumbers")
	mainObjects.textWinTextview.SetName("textWinTextview")
	mainObjects.textWinTextviewNumbers.SetName("textWinTextviewNumbers")

}
