// gohObjects.go

/*
	Source file auto-generated on Thu, 13 Aug 2020 04:08:31 using Gotk3ObjHandler v1.5 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-20 H.F.M - Search And Replace v1.8 github.com/hfmrow/sAndReplace
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
	BoxTextViewPreview              *gtk.Grid
	btnFind                         *gtk.Button
	btnReplaceInClipboard           *gtk.Button
	btnScan                         *gtk.Button
	btnShowClipboard                *gtk.Button
	chkCaseSensitive                *gtk.CheckButton
	chkCharacterClass               *gtk.CheckButton
	chkCharacterClassStrictMode     *gtk.CheckButton
	chkFollowSymlinkDir             *gtk.CheckButton
	chkRegex                        *gtk.CheckButton
	chkUseEscapeChar                *gtk.CheckButton
	chkWholeWord                    *gtk.CheckButton
	chkWildcard                     *gtk.CheckButton
	clipboard                       *gtk.Clipboard /*MANUAL*/
	entryExtMask                    *gtk.Entry
	entryReplace                    *gtk.Entry
	entrySearch                     *gtk.Entry
	fileChooserBtn                  *gtk.FileChooserButton
	findWin                         *gtk.Window
	findWinCancelBtn                *gtk.Button
	findWinChkBackUp                *gtk.CheckButton
	findWinChkDispForbFiles         *gtk.CheckButton
	findWinGrid                     *gtk.Grid
	findWinGridBtn                  *gtk.Grid
	findWinReplaceBtn               *gtk.Button
	findWinScrollWin                *gtk.ScrolledWindow
	findWinTreeView                 *gtk.TreeView
	imgTop                          *gtk.Image
	ImgTopEventbox                  *gtk.EventBox
	lblExtMask                      *gtk.Label
	lblReplace                      *gtk.Label
	lblSearch                       *gtk.Label
	listViewFiles                   *gtk.TreeView
	mainBox                         *gtk.Box
	MainButtonOptions               *gtk.Button
	MainGrid                        *gtk.Grid
	MainTopGrig                     *gtk.Grid
	mainUiBuilder                   *gtk.Builder
	mainWin                         *gtk.Window
	mainWinBtnClose                 *gtk.Button
	OptionButtonDone                *gtk.Button
	OptionsEntryMaxFileSize         *gtk.Entry
	OptionsEntryMinFileSize         *gtk.Entry
	OptionsImageTop                 *gtk.Image
	OptionsLabelChooseDepth         *gtk.Label
	OptionsLabelEnableDirChooser    *gtk.Label
	OptionsLabelMaxFileSize         *gtk.Label
	OptionsLabelMinFileSize         *gtk.Label
	OptionsWindow                   *gtk.Window
	replaceGrid                     *gtk.Grid
	scrolledWindowTreeView          *gtk.ScrolledWindow
	spinButtonDepth                 *gtk.SpinButton
	statusbar                       *gtk.Statusbar
	switchFileChooserButton         *gtk.Switch
	textWin                         *gtk.Window
	textWinBox                      *gtk.Box
	textWinBtnDone                  *gtk.Button
	textWinChkShowModifications     *gtk.CheckButton
	textWinChkSyntxHighlight        *gtk.CheckButton
	textWinChkWrap                  *gtk.CheckButton
	textWinComboBoxLanguage         *gtk.ComboBoxText
	textWinComboBoxTextStyleChooser *gtk.ComboBoxText
	textWinGridTop                  *gtk.Grid
	treeviewSelection               *gtk.TreeSelection
}

/******************************/
/* GtkObjects  Initialisation */
/******************************/
// gladeObjParser: Initialise Gtk3 Objects into mainObjects structure.
func gladeObjParser() {
	mainObjects.BoxTextViewPreview = loadObject("BoxTextViewPreview").(*gtk.Grid)
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
	mainObjects.chkWholeWord = loadObject("chkWholeWord").(*gtk.CheckButton)
	mainObjects.chkWildcard = loadObject("chkWildcard").(*gtk.CheckButton)
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
	mainObjects.MainButtonOptions = loadObject("MainButtonOptions").(*gtk.Button)
	mainObjects.MainGrid = loadObject("MainGrid").(*gtk.Grid)
	mainObjects.MainTopGrig = loadObject("MainTopGrig").(*gtk.Grid)
	mainObjects.mainWin = loadObject("mainWin").(*gtk.Window)
	mainObjects.mainWinBtnClose = loadObject("mainWinBtnClose").(*gtk.Button)
	mainObjects.OptionButtonDone = loadObject("OptionButtonDone").(*gtk.Button)
	mainObjects.OptionsEntryMaxFileSize = loadObject("OptionsEntryMaxFileSize").(*gtk.Entry)
	mainObjects.OptionsEntryMinFileSize = loadObject("OptionsEntryMinFileSize").(*gtk.Entry)
	mainObjects.OptionsImageTop = loadObject("OptionsImageTop").(*gtk.Image)
	mainObjects.OptionsLabelChooseDepth = loadObject("OptionsLabelChooseDepth").(*gtk.Label)
	mainObjects.OptionsLabelEnableDirChooser = loadObject("OptionsLabelEnableDirChooser").(*gtk.Label)
	mainObjects.OptionsLabelMaxFileSize = loadObject("OptionsLabelMaxFileSize").(*gtk.Label)
	mainObjects.OptionsLabelMinFileSize = loadObject("OptionsLabelMinFileSize").(*gtk.Label)
	mainObjects.OptionsWindow = loadObject("OptionsWindow").(*gtk.Window)
	mainObjects.replaceGrid = loadObject("replaceGrid").(*gtk.Grid)
	mainObjects.scrolledWindowTreeView = loadObject("scrolledWindowTreeView").(*gtk.ScrolledWindow)
	mainObjects.spinButtonDepth = loadObject("spinButtonDepth").(*gtk.SpinButton)
	mainObjects.statusbar = loadObject("statusbar").(*gtk.Statusbar)
	mainObjects.switchFileChooserButton = loadObject("switchFileChooserButton").(*gtk.Switch)
	mainObjects.textWin = loadObject("textWin").(*gtk.Window)
	mainObjects.textWinBox = loadObject("textWinBox").(*gtk.Box)
	mainObjects.textWinBtnDone = loadObject("textWinBtnDone").(*gtk.Button)
	mainObjects.textWinChkShowModifications = loadObject("textWinChkShowModifications").(*gtk.CheckButton)
	mainObjects.textWinChkSyntxHighlight = loadObject("textWinChkSyntxHighlight").(*gtk.CheckButton)
	mainObjects.textWinChkWrap = loadObject("textWinChkWrap").(*gtk.CheckButton)
	mainObjects.textWinComboBoxLanguage = loadObject("textWinComboBoxLanguage").(*gtk.ComboBoxText)
	mainObjects.textWinComboBoxTextStyleChooser = loadObject("textWinComboBoxTextStyleChooser").(*gtk.ComboBoxText)
	mainObjects.textWinGridTop = loadObject("textWinGridTop").(*gtk.Grid)
	mainObjects.treeviewSelection = loadObject("treeviewSelection").(*gtk.TreeSelection)
}
