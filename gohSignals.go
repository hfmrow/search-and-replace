// gohSignals.go

/*
	Source file auto-generated on Fri, 09 Apr 2021 03:01:52 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 hfmrow - Search And Replace v1.10 github.com/hfmrow/search-and-replace
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
	obj.btnFind.Connect("clicked", btnFindClicked)
	obj.btnReplaceInClipboard.Connect("clicked", btnReplaceInClipboardClicked)
	obj.btnScan.Connect("clicked", btnScanClicked)
	obj.btnShowClipboard.Connect("clicked", btnShowClipboardClicked)
	obj.chkCaseSensitive.Connect("toggled", chkCaseSensitiveToggled)
	obj.chkCharacterClass.Connect("toggled", chkCharacterClassToggled)
	obj.chkFollowSymlinkDir.Connect("toggled", chkFollowSymlinkDirToggled)
	obj.chkRegex.Connect("toggled", chkRegexToggled)
	obj.chkUseEscapeChar.Connect("toggled", chkUseEscapeCharToggled)
	obj.chkWholeWord.Connect("toggled", chkWholeWordToggled)
	obj.chkWildcard.Connect("toggled", chkWildcardToggled)
	obj.entryExtMask.Connect("activate", entryExtMaskFocusOut)
	obj.entryReplace.Connect("activate", genericEntryFocusOut)
	obj.entrySearch.Connect("activate", genericEntryFocusOut)
	obj.fileChooserBtn.Connect("selection-changed", fileChooserBtnSelectionChanged)
	obj.findWin.Connect("delete_event", genericHideWindow)
	obj.findWinBtnDeselect.Connect("clicked", findWinBtnDeselectClicked)
	obj.findWinBtnInvertSel.Connect("clicked", findWinBtnInvertSelClicked)
	obj.findWinCancelBtn.Connect("clicked", findWinCancelBtnClicked)
	obj.findWinChkBackUp.Connect("toggled", findWinChkBackUpToggled)
	obj.findWinChkDispForbFiles.Connect("toggled", findWinChkDispForbFilesToggled)
	obj.findWinChkExpandAll.Connect("toggled", findWinChkExpandAllToggled)
	obj.findWinReplaceBtn.Connect("clicked", findWinReplaceBtnClicked)
	obj.findWinTreeView.Connect("row-activated", findTreeViewDblClick)
	obj.ImgTopEventbox.Connect("button-release-event", imgTopReleaseEvent)
	obj.listViewFiles.Connect("row-activated", listViewFilesRowActivated)
	obj.listViewFiles.Connect("button-press-event", popupMenu.CheckRMBFromTreeView)
	obj.MainButtonOptions.Connect("clicked", MainButtonOptionsClicked)
	obj.mainWin.Connect("delete-event", mainWinOnExit)
	obj.mainWinBtnClose.Connect("clicked", windowDestroy)
	obj.OptionButtonDone.Connect("clicked", OptionButtonDoneClicked)
	obj.OptionsEntryMaxFileSize.Connect("changed", OptionsEntryFileSizeChanged)
	obj.OptionsEntryMinFileSize.Connect("changed", OptionsEntryFileSizeChanged)
	obj.SourceToggleButtonMapWidth.Connect("toggled", SourceToggleButtonMapWidthToggled)
	obj.spinButtonDepth.Connect("value-changed", spinButtonDepthValueChanged)
	obj.textWin.Connect("delete_event", genericHideWindow)
	obj.textWin.Connect("check-resize", WindowSourceCheckResize)
	obj.textWinBtnDone.Connect("clicked", textWinBtnDoneClicked)
	obj.textWinChkShowModifications.Connect("toggled", textWinChkShowModificationsToggled)
	obj.textWinChkSyntxHighlight.Connect("toggled", textWinChkSyntxHighlightToggled)
	obj.textWinChkWrap.Connect("toggled", textWinChkWrapToggled)
	obj.View.Connect("button-press-event", ViewButtonPressEvent)
}
