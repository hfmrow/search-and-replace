// gohObjects.go

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
	"github.com/gotk3/gotk3/gtk"

	"github.com/hfmrow/gotk3_gtksource/source"
)

// Control over all used objects from glade.
var obj *MainControlsObj

/******************************/
/* Main structure Declaration */
/******************************/
type MainControlsObj struct {
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
	chkUseEscapeCharToReplace       *gtk.CheckButton
	chkWholeWord                    *gtk.CheckButton
	chkWildcard                     *gtk.CheckButton
	clipboard                       *gtk.Clipboard /*MANUAL*/
	entryExtMask                    *gtk.Entry
	entryReplace                    *gtk.Entry
	entrySearch                     *gtk.Entry
	fileChooserBtn                  *gtk.FileChooserButton
	findWin                         *gtk.Window
	findWinBtnDeselect              *gtk.Button
	findWinBtnInvertSel             *gtk.Button
	findWinCancelBtn                *gtk.Button
	findWinChkBackUp                *gtk.CheckButton
	findWinChkDispForbFiles         *gtk.CheckButton
	findWinChkExpandAll             *gtk.CheckButton
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
	MainTopGrig                     *gtk.Grid
	mainUiBuilder                   *gtk.Builder
	mainWin                         *gtk.Window
	mainWinBtnClose                 *gtk.Button
	Map                             *source.SourceMap
	OptionButtonDone                *gtk.Button
	OptionsEntryMaxFileSize         *gtk.Entry
	OptionsEntryMinFileSize         *gtk.Entry
	OptionsImageTop                 *gtk.Image
	OptionsLabelChooseDepth         *gtk.Label
	OptionsLabelEnableDirChooser    *gtk.Label
	OptionsLabelMaxFileSize         *gtk.Label
	OptionsLabelMinFileSize         *gtk.Label
	OptionsWindow                   *gtk.Window
	Paned                           *gtk.Paned
	replaceGrid                     *gtk.Grid
	scrolledWindowTreeView          *gtk.ScrolledWindow
	SourceToggleButtonMapWidth      *gtk.ToggleButton
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
	View                            *source.SourceView
}

/******************************/
/* GtkObjects  Initialisation */
/******************************/
// gladeObjParser: Initialise Gtk3 Objects into obj structure.
func gladeObjParser() {
	obj.btnFind = loadObject("btnFind").(*gtk.Button)
	obj.btnReplaceInClipboard = loadObject("btnReplaceInClipboard").(*gtk.Button)
	obj.btnScan = loadObject("btnScan").(*gtk.Button)
	obj.btnShowClipboard = loadObject("btnShowClipboard").(*gtk.Button)
	obj.chkCaseSensitive = loadObject("chkCaseSensitive").(*gtk.CheckButton)
	obj.chkCharacterClass = loadObject("chkCharacterClass").(*gtk.CheckButton)
	obj.chkCharacterClassStrictMode = loadObject("chkCharacterClassStrictMode").(*gtk.CheckButton)
	obj.chkFollowSymlinkDir = loadObject("chkFollowSymlinkDir").(*gtk.CheckButton)
	obj.chkRegex = loadObject("chkRegex").(*gtk.CheckButton)
	obj.chkUseEscapeChar = loadObject("chkUseEscapeChar").(*gtk.CheckButton)
	obj.chkUseEscapeCharToReplace = loadObject("chkUseEscapeCharToReplace").(*gtk.CheckButton)
	obj.chkWholeWord = loadObject("chkWholeWord").(*gtk.CheckButton)
	obj.chkWildcard = loadObject("chkWildcard").(*gtk.CheckButton)
	obj.entryExtMask = loadObject("entryExtMask").(*gtk.Entry)
	obj.entryReplace = loadObject("entryReplace").(*gtk.Entry)
	obj.entrySearch = loadObject("entrySearch").(*gtk.Entry)
	obj.fileChooserBtn = loadObject("fileChooserBtn").(*gtk.FileChooserButton)
	obj.findWin = loadObject("findWin").(*gtk.Window)
	obj.findWinBtnDeselect = loadObject("findWinBtnDeselect").(*gtk.Button)
	obj.findWinBtnInvertSel = loadObject("findWinBtnInvertSel").(*gtk.Button)
	obj.findWinCancelBtn = loadObject("findWinCancelBtn").(*gtk.Button)
	obj.findWinChkBackUp = loadObject("findWinChkBackUp").(*gtk.CheckButton)
	obj.findWinChkDispForbFiles = loadObject("findWinChkDispForbFiles").(*gtk.CheckButton)
	obj.findWinChkExpandAll = loadObject("findWinChkExpandAll").(*gtk.CheckButton)
	obj.findWinGrid = loadObject("findWinGrid").(*gtk.Grid)
	obj.findWinGridBtn = loadObject("findWinGridBtn").(*gtk.Grid)
	obj.findWinReplaceBtn = loadObject("findWinReplaceBtn").(*gtk.Button)
	obj.findWinScrollWin = loadObject("findWinScrollWin").(*gtk.ScrolledWindow)
	obj.findWinTreeView = loadObject("findWinTreeView").(*gtk.TreeView)
	obj.imgTop = loadObject("imgTop").(*gtk.Image)
	obj.ImgTopEventbox = loadObject("ImgTopEventbox").(*gtk.EventBox)
	obj.lblExtMask = loadObject("lblExtMask").(*gtk.Label)
	obj.lblReplace = loadObject("lblReplace").(*gtk.Label)
	obj.lblSearch = loadObject("lblSearch").(*gtk.Label)
	obj.listViewFiles = loadObject("listViewFiles").(*gtk.TreeView)
	obj.mainBox = loadObject("mainBox").(*gtk.Box)
	obj.MainButtonOptions = loadObject("MainButtonOptions").(*gtk.Button)
	obj.MainTopGrig = loadObject("MainTopGrig").(*gtk.Grid)
	obj.mainWin = loadObject("mainWin").(*gtk.Window)
	obj.mainWinBtnClose = loadObject("mainWinBtnClose").(*gtk.Button)
	obj.Map = loadObject("Map").(*source.SourceMap)
	obj.OptionButtonDone = loadObject("OptionButtonDone").(*gtk.Button)
	obj.OptionsEntryMaxFileSize = loadObject("OptionsEntryMaxFileSize").(*gtk.Entry)
	obj.OptionsEntryMinFileSize = loadObject("OptionsEntryMinFileSize").(*gtk.Entry)
	obj.OptionsImageTop = loadObject("OptionsImageTop").(*gtk.Image)
	obj.OptionsLabelChooseDepth = loadObject("OptionsLabelChooseDepth").(*gtk.Label)
	obj.OptionsLabelEnableDirChooser = loadObject("OptionsLabelEnableDirChooser").(*gtk.Label)
	obj.OptionsLabelMaxFileSize = loadObject("OptionsLabelMaxFileSize").(*gtk.Label)
	obj.OptionsLabelMinFileSize = loadObject("OptionsLabelMinFileSize").(*gtk.Label)
	obj.OptionsWindow = loadObject("OptionsWindow").(*gtk.Window)
	obj.Paned = loadObject("Paned").(*gtk.Paned)
	obj.replaceGrid = loadObject("replaceGrid").(*gtk.Grid)
	obj.scrolledWindowTreeView = loadObject("scrolledWindowTreeView").(*gtk.ScrolledWindow)
	obj.SourceToggleButtonMapWidth = loadObject("SourceToggleButtonMapWidth").(*gtk.ToggleButton)
	obj.spinButtonDepth = loadObject("spinButtonDepth").(*gtk.SpinButton)
	obj.statusbar = loadObject("statusbar").(*gtk.Statusbar)
	obj.switchFileChooserButton = loadObject("switchFileChooserButton").(*gtk.Switch)
	obj.textWin = loadObject("textWin").(*gtk.Window)
	obj.textWinBox = loadObject("textWinBox").(*gtk.Box)
	obj.textWinBtnDone = loadObject("textWinBtnDone").(*gtk.Button)
	obj.textWinChkShowModifications = loadObject("textWinChkShowModifications").(*gtk.CheckButton)
	obj.textWinChkSyntxHighlight = loadObject("textWinChkSyntxHighlight").(*gtk.CheckButton)
	obj.textWinChkWrap = loadObject("textWinChkWrap").(*gtk.CheckButton)
	obj.textWinComboBoxLanguage = loadObject("textWinComboBoxLanguage").(*gtk.ComboBoxText)
	obj.textWinComboBoxTextStyleChooser = loadObject("textWinComboBoxTextStyleChooser").(*gtk.ComboBoxText)
	obj.textWinGridTop = loadObject("textWinGridTop").(*gtk.Grid)
	obj.treeviewSelection = loadObject("treeviewSelection").(*gtk.TreeSelection)
	obj.View = loadObject("View").(*source.SourceView)
}
