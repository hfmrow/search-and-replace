// gohSignals.go

/*
	Source file auto-generated on Thu, 13 Aug 2020 04:08:31 using Gotk3ObjHandler v1.5 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-20 H.F.M - Search And Replace v1.8 github.com/hfmrow/sAndReplace
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
	mainObjects.entryExtMask.Connect("focus-out-event", entryExtMaskFocusOut)
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
	mainObjects.spinButtonDepth.Connect("value-changed", spinButtonDepthValueChanged)
	mainObjects.textWin.Connect("delete_event", genericHideWindow)
	mainObjects.textWinBtnDone.Connect("clicked", textWinBtnDoneClicked)
	mainObjects.textWinChkShowModifications.Connect("toggled", textWinChkShowModificationsToggled)
	mainObjects.textWinChkSyntxHighlight.Connect("toggled", textWinComboBoxTextStyleChooserChanged)
	mainObjects.textWinChkWrap.Connect("toggled", textWinChkWrapToggled)
	mainObjects.textWinComboBoxLanguage.Connect("changed", textWinComboBoxTextStyleChooserChanged)
	mainObjects.textWinComboBoxTextStyleChooser.Connect("changed", textWinComboBoxTextStyleChooserChanged)
}
