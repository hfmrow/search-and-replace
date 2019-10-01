// gohSignals.go

// Source file auto-generated on Tue, 01 Oct 2019 16:36:13 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
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
	mainObjects.chkWildcard.Connect("toggled", chkWildcardToggled)
	mainObjects.chkWoleWord.Connect("toggled", chkWoleWordToggled)
	mainObjects.entryExtMask.Connect("focus-out-event", entryExtMaskFocusOut)
	mainObjects.entryExtMask.Connect("activate", entryExtMaskFocusOut)
	mainObjects.entryReplace.Connect("notify", blankNotify)
	mainObjects.entrySearch.Connect("activate", entrySearchFocusOut)
	fileChooserSigHandlerSelChanged, _ = mainObjects.fileChooserBtn.Connect("selection-changed", fileChooserBtnSelectionChanged)
	mainObjects.findWin.Connect("delete_event", genericHideWindow)
	mainObjects.findWinCancelBtn.Connect("clicked", findWinCancelBtnClicked)
	mainObjects.findWinChkBackUp.Connect("toggled", findWinChkBackUpToggled)
	mainObjects.findWinChkDispForbFiles.Connect("toggled", findWinChkDispForbFilesToggled)
	mainObjects.findWinReplaceBtn.Connect("clicked", findWinReplaceBtnClicked)
	mainObjects.findWinTreeView.Connect("row-activated", findTreeViewDblClick)
	mainObjects.ImgTopEventbox.Connect("button-release-event", imgTopReleaseEvent)
	mainObjects.listViewFiles.Connect("row-activated", listViewFilesRowActivated)
	mainObjects.mainWin.Connect("delete-event", mainWinOnExit)
	mainObjects.mainWinBtnClose.Connect("clicked", windowDestroy)
	spinButtonDepthSigHandlerChanged, _ = mainObjects.spinButtonDepth.Connect("value-changed", spinButtonDepthValueChanged)
	mainObjects.switchFileChooserButton.Connect("state-set", switchFileChooserButtonStateSet)
	mainObjects.switchFileChooserButton.Connect("event-after", switchFileChooserButtonEventAfter)
	mainObjects.textWin.Connect("delete_event", genericHideWindow)
	mainObjects.textWinBtnDone.Connect("clicked", textWinBtnDoneClicked)
	mainObjects.textWinChkShowModifications.Connect("toggled", textWinChkShowModificationsToggled)
	mainObjects.textWinChkWrap.Connect("toggled", textWinChkWrapToggled)
	mainObjects.textWinScrolledwindowNumbers.Connect("notify", blankNotify)
	mainObjects.textWinTextviewNumbers.Connect("notify", blankNotify)

}
