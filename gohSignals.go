// gohSignals.go

/*
	Source file auto-generated on Fri, 02 Apr 2021 10:53:33 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 H.F.M - Search And Replace v1.9 github.com/hfmrow/search-and-replace
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/*********************************************************/
/* Signals Implementations:                             */
/* initialise signals used by gtk objects ...          */
/* This section preserve user modification on update. */
/*****************************************************/
func signalsPropHandler() {
	mainObjects.btnFind.Connect("clicked", btnFindClicked)
	mainObjects.btnReplaceInClipboard.Connect("clicked", btnReplaceInClipboardClicked)
	mainObjects.btnScan.Connect("clicked", btnScanClicked)
	mainObjects.btnShowClipboard.Connect("clicked", btnShowClipboardClicked)
	mainObjects.chkCaseSensitive.Connect("toggled", chkCaseSensitiveToggled)
	mainObjects.chkCharacterClass.Connect("toggled", chkCharacterClassToggled)
	mainObjects.chkFollowSymlinkDir.Connect("toggled", chkFollowSymlinkDirToggled)
	mainObjects.chkRegex.Connect("toggled", chkRegexToggled)
	mainObjects.chkUseEscapeChar.Connect("toggled", chkUseEscapeCharToggled)
	mainObjects.chkWholeWord.Connect("toggled", chkWholeWordToggled)
	mainObjects.chkWildcard.Connect("toggled", chkWildcardToggled)
	mainObjects.entryExtMask.Connect("activate", entryExtMaskFocusOut)
	mainObjects.entryReplace.Connect("activate", genericEntryFocusOut)
	mainObjects.entrySearch.Connect("activate", genericEntryFocusOut)
	mainObjects.fileChooserBtn.Connect("selection-changed", fileChooserBtnSelectionChanged)
	mainObjects.findWin.Connect("delete_event", genericHideWindow)
	mainObjects.findWinCancelBtn.Connect("clicked", findWinCancelBtnClicked)
	mainObjects.findWinChkBackUp.Connect("toggled", findWinChkBackUpToggled)
	mainObjects.findWinChkDispForbFiles.Connect("toggled", findWinChkDispForbFilesToggled)
	mainObjects.findWinReplaceBtn.Connect("clicked", findWinReplaceBtnClicked)
	mainObjects.findWinTreeView.Connect("row-activated", findTreeViewDblClick)
	mainObjects.ImgTopEventbox.Connect("button-release-event", imgTopReleaseEvent)
	mainObjects.listViewFiles.Connect("row-activated", listViewFilesRowActivated)
	mainObjects.listViewFiles.Connect("button-press-event", popupMenu.CheckRMBFromTreeView)
	mainObjects.MainButtonOptions.Connect("clicked", MainButtonOptionsClicked)
	mainObjects.mainWin.Connect("delete-event", mainWinOnExit)
	mainObjects.mainWinBtnClose.Connect("clicked", windowDestroy)
	mainObjects.OptionButtonDone.Connect("clicked", OptionButtonDoneClicked)
	mainObjects.OptionsEntryMaxFileSize.Connect("changed", OptionsEntryFileSizeChanged)
	mainObjects.OptionsEntryMinFileSize.Connect("changed", OptionsEntryFileSizeChanged)
	mainObjects.SourceToggleButtonMapWidth.Connect("toggled", SourceToggleButtonMapWidthToggled)
	mainObjects.spinButtonDepth.Connect("value-changed", spinButtonDepthValueChanged)
	mainObjects.textWin.Connect("delete_event", genericHideWindow)
	mainObjects.textWin.Connect("check-resize", WindowSourceCheckResize)
	mainObjects.textWinBtnDone.Connect("clicked", textWinBtnDoneClicked)
	mainObjects.textWinChkShowModifications.Connect("toggled", textWinChkShowModificationsToggled)
	mainObjects.textWinChkSyntxHighlight.Connect("toggled", textWinChkSyntxHighlightToggled)
	mainObjects.textWinChkWrap.Connect("toggled", textWinChkWrapToggled)
	mainObjects.View.Connect("button-press-event", ViewButtonPressEvent)
}
