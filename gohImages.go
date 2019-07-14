// gohImages.go

// Source file auto-generated on Sun, 14 Jul 2019 16:40:30 using Gotk3ObjHandler v1.3 ©2019 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/************************************************************/
/* Images declarations, used to initialize objects with it */
/* The functions: setImage, setWinIcon and setButtonImage */
/* accept both kind of variables: filename or []byte     */
/* content in case of using embedded binary data. The   */
/* variables names are the same. You can use function  */
/* "func assetsDeclarationsUseEmbedded(bool)"         */
/* to toggle between filenames and embedded binary.  */
/****************************************************/
func assignImages() {
	setButtonImage(mainObjects.btnFind, view18x18)
	setButtonImage(mainObjects.btnReplaceInClipboard, viewDetails18x18)
	setButtonImage(mainObjects.btnShowClipboard, scorecard18x18)
	setWinIcon(mainObjects.findWin, findMultipleDocuments48x48)
	setButtonImage(mainObjects.findWinCancelBtn, cancel18x18)
	setButtonImage(mainObjects.findWinReplaceBtn, checked18x18)
	setImage(mainObjects.imgTop, sanderSearchAndReplace700x48)
	setWinIcon(mainObjects.mainWin, findMultipleDocuments48x48)
	setWinIcon(mainObjects.textWin, findMultipleDocuments48x48)
	setButtonImage(mainObjects.textWinBtnDone, checked18x18)
}

// Assets var declarations, this step permit to make a "bridge" between the differents
// types used: (string or []byte) and to simply switch from one to another.
var cancel18x18 interface{}                  // assets/images/cancel-18x18.png
var checked18x18 interface{}                 // assets/images/checked-18x18.png
var findMultipleDocuments48x48 interface{}   // assets/images/find-multiple-documents-48x48.png
var mainGlade interface{}                    // assets/glade/main.glade
var sanderSearchAndReplace400x27 interface{} // assets/images/Sander-search-and-replace-400x27.png
var sanderSearchAndReplace700x48 interface{} // assets/images/Sander-search-and-replace-700x48.png
var scorecard18x18 interface{}               // assets/images/scorecard-18x18.png
var view18x18 interface{}                    // assets/images/view-18x18.png
var viewDetails18x18 interface{}             // assets/images/view-details-18x18.png
