// gohSignals.go

// Source file auto-generated on Sun, 14 Jul 2019 16:40:30 using Gotk3ObjHandler v1.3 ©2019 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/***************************/
/* Signals Implementations */
/***************************/
// signalsPropHandler: initialise signals used by gtk objects ...
func signalsPropHandler() {
	mainObjects.btnFind.Connect("clicked", btnFindClicked)                             /*MANUAL*/
	mainObjects.btnReplaceInClipboard.Connect("clicked", btnReplaceInClipboardClicked) /*MANUAL*/
	mainObjects.btnShowClipboard.Connect("clicked", btnShowClipboardClicked)           /*MANUAL*/
	mainObjects.chkCaseSensitive.Connect("toggled", chkCaseSensitiveToggled)           /*MANUAL*/
	mainObjects.chkCharacterClass.Connect("toggled", chkCharacterClassToggled)
	mainObjects.chkFollowSymlinkDir.Connect("toggled", fileChooserBtnClicked) /*MANUAL*/
	mainObjects.chkRegex.Connect("toggled", chkRegexToggled)                  /*MANUAL*/
	mainObjects.chkSubDir.Connect("toggled", fileChooserBtnClicked)           /*MANUAL*/
	mainObjects.chkUseEscapeChar.Connect("toggled", chkUseEscapeCharToggled)  /*MANUAL*/
	mainObjects.chkWildcard.Connect("toggled", chkWildcardToggled)            /*MANUAL*/
	mainObjects.chkWoleWord.Connect("toggled", chkWoleWordToggled)            /*MANUAL*/
	mainObjects.entryExtMask.Connect("focus-out-event", entryExtMaskFocusOut) /*MANUAL*/
	// mainObjects.entryExtMask.Connect("focus-out-event", fileChooserBtnClicked) /*MANUAL*/
	mainObjects.fileChooserBtn.Connect("selection-changed", fileChooserBtnClicked) /*MANUAL*/
	mainObjects.findWinCancelBtn.Connect("clicked", findWinCancelBtnClicked)       /*MANUAL*/
	mainObjects.findWinChkBackUp.Connect("toggled", findWinChkBackUpToggled)       /*MANUAL*/
	mainObjects.findWinReplaceBtn.Connect("clicked", findWinReplaceBtnClicked)     /*MANUAL*/
	mainObjects.findWinTreeView.Connect("row-activated", findTreeViewDblClick)     /*MANUAL*/
	mainObjects.ImgTopEventbox.Connect("button-release-event", imgTopReleaseEvent) /*MANUAL*/
	mainObjects.mainWin.Connect("delete-event", mainWinOnExit)
	mainObjects.textWinBtnDone.Connect("clicked", textWinBtnDoneClicked)                           /*MANUAL*/
	mainObjects.textWinChkShowModifications.Connect("toggled", textWinChkShowModificationsToggled) /*MANUAL*/
	mainObjects.textWinChkWrap.Connect("toggled", textWinChkWrapToggled)                           /*MANUAL*/
	initDropSets()                                                                                 /*MANUAL*/
	mainObjects.treeviewFiles.Connect("row-activated", treeViewDblClick)                           /*MANUAL*/
	mainObjects.treeviewFiles.Connect("drag-data-received", treeviewFilesReceived)                 /*MANUAL*/
	mainObjects.treeviewSelection.Connect("changed", treeViewSelectionChanged)                     /*MANUAL*/
}
