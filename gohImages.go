// gohImages.go

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
)

/**********************************************************/
/* This section preserve user modifications on update.   */
/* Images declarations, used to initialize objects with */
/* The SetPict() func, accept both kind of variables:  */
/* filename or []byte content in case of using        */
/* embedded binary data. The variables names are the */
/* same. "assetsDeclarationsUseEmbedded(bool)" func */
/* could be used to toggle between filenames and   */
/* embedded binary type. See SetPict()            */
/* declaration to learn more on how to use it.   */
/************************************************/
func assignImages() {
	SetPict(obj.btnFind, searchIcon48, 24) // www.freeiconspng.com
	SetPict(obj.btnReplaceInClipboard, clipboardRepl, 24)
	SetPict(obj.btnScan, folder48, 24)
	SetPict(obj.btnShowClipboard, clipboard, 24)
	SetPict(obj.findWin, searchFolder48)
	SetPict(obj.findWinBtnDeselect, "")
	SetPict(obj.findWinBtnInvertSel, "")
	SetPict(obj.findWinCancelBtn, crossIcon48, 24)
	SetPict(obj.findWinReplaceBtn, replace, 24)
	SetPict(obj.imgTop, searchAndReplaceTop48)
	SetPict(obj.MainButtonOptions, options, 18)
	SetPict(obj.mainWin, searchFolder48)
	SetPict(obj.mainWinBtnClose, logout48, 18)
	SetPict(obj.OptionButtonDone, tickIcon48, 24)
	SetPict(obj.OptionsImageTop, options)
	SetPict(obj.OptionsWindow, searchFolder48)
	SetPict(obj.SourceToggleButtonMapWidth, "")
	SetPict(obj.spinButtonDepth, folder48, OptPict{Position: gtk.ENTRY_ICON_SECONDARY})
	SetPict(obj.textWin, searchFolder48)
	SetPict(obj.textWinBtnDone, tickIcon48, 24)
}

/**********************************************************/
/* This section is rewritten on assets update.           */
/* Assets var declarations, this step permit to make a  */
/* bridge between the differents types used, string or */
/* []byte, and to simply switch from one to another.  */
/*****************************************************/
var mainGlade interface{}                  // assets/glade/main.glade
var clipboard interface{}                  // assets/images/clipboard.png
var clipboardRepl interface{}              // assets/images/clipboard-repl.png
var crossIcon48 interface{}                // assets/images/Cross-icon-48.png
var folder48 interface{}                   // assets/images/folder-48.png
var linearProgressHorzBlue interface{}     // assets/images/linear-progress-horz-blue.gif
var logout48 interface{}                   // assets/images/logout-48.png
var mimetypeSourceIconGolang48 interface{} // assets/images/Mimetype-source-icon-golang-48.png
var options interface{}                    // assets/images/Options.png
var replace interface{}                    // assets/images/replace.png
var searchAndReplaceTop27 interface{}      // assets/images/search-and-replace-top-27.png
var searchAndReplaceTop48 interface{}      // assets/images/search-and-replace-top-48.png
var searchFolder48 interface{}             // assets/images/search-folder-48.png
var searchIcon48 interface{}               // assets/images/search-icon-48.png
var tickIcon48 interface{}                 // assets/images/Tick-icon-48.png
