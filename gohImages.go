// gohImages.go

// Source file auto-generated on Tue, 01 Oct 2019 16:36:13 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/*************************************************************/
/* This section preserve user modifications on update.      */
/* Images declarations, used to initialize objects with it */
/* The functions: setImage, setWinIcon and setButtonImage */
/* accept both kind of variables: filename or []byte     */
/* content in case of using embedded binary data. The   */
/* variables names are the same. You can use function  */
/* "func assetsDeclarationsUseEmbedded(bool)"         */
/* to toggle between filenames and embedded binary.  */
/* The "set...Image/Icon" functions could receive   */
/* a display size if necessary or a pos. argument. */
/**************************************************/
func assignImages() {
	setButtonImage(mainObjects.btnFind, searchB36, 32)
	setButtonImage(mainObjects.btnReplaceInClipboard, clipboardPasteIcon18)
	setButtonImage(mainObjects.btnScan, searchHardDisk20, 18)
	setButtonImage(mainObjects.btnShowClipboard, clipboardIcon18)
	setWinIcon(mainObjects.findWin, findMultipleDocuments48x48)
	setButtonImage(mainObjects.findWinCancelBtn, signCancel20)
	setButtonImage(mainObjects.findWinReplaceBtn, textReplace20)
	setImage(mainObjects.imgTop, sanderSearchAndReplace700x48)
	setWinIcon(mainObjects.mainWin, findMultipleDocuments48x48)
	setButtonImage(mainObjects.mainWinBtnClose, signCancel20)
	setSpinButtonImage(mainObjects.spinButtonDepth, filesystemsFolder20, "left")
	setWinIcon(mainObjects.textWin, findMultipleDocuments48x48)
	setButtonImage(mainObjects.textWinBtnDone, signSelect20)

}

/**********************************************************/
/* This section is rewritten on assets update.           */
/* Assets var declarations, this step permit to make a  */
/* bridge between the differents types used, string or */
/* []byte, and to simply switch from one to another.  */
/*****************************************************/
var clipboardIcon18 interface{}              // assets/images/Clipboard-icon-18.png
var clipboardPasteIcon18 interface{}         // assets/images/Clipboard-Paste-icon-18.png
var filesystemsFolder20 interface{}          // assets/images/Filesystems-folder-20.png
var findMultipleDocuments48x48 interface{}   // assets/images/find-multiple-documents-48x48.png
var mainGlade interface{}                    // assets/glade/main.glade
var sanderSearchAndReplace400x27 interface{} // assets/images/Sander-search-and-replace-400x27.png
var sanderSearchAndReplace700x48 interface{} // assets/images/Sander-search-and-replace-700x48.png
var searchB36 interface{}                    // assets/images/search-b-36.png
var searchB38 interface{}                    // assets/images/search-b-38.png
var searchHardDisk20 interface{}             // assets/images/Search-Hard-Disk-20.png
var signCancel20 interface{}                 // assets/images/Sign-cancel-20.png
var signSelect20 interface{}                 // assets/images/Sign-Select-20.png
var textReplace20 interface{}                // assets/images/text-replace-20.png
