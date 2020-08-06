// gohOptions.go

/*
	Source file auto-generated on Thu, 06 Aug 2020 20:25:48 using Gotk3ObjHandler v1.5 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-20 H.F.M - Search And Replace v1.8 github.com/hfmrow/sAndReplace
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	glfsft "github.com/hfmrow/genLib/files/fileText"
	glfssf "github.com/hfmrow/genLib/files/scanFileDir"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gitv "github.com/hfmrow/gtk3Import/textView"
	gitvch "github.com/hfmrow/gtk3Import/textView/chromaHighlight"
	gitl "github.com/hfmrow/gtk3Import/tools"
	gits "github.com/hfmrow/gtk3Import/tools"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

// App infos
var Name = "Search And Replace"
var Vers = "v1.8"
var Descr = "Search and replace text sequences in one or more files.\nThe clipboard is also included inside replacing functionality.\nThere is some useful functionalities like code preview with\nsyntax highlighting, backup modified files whether is requested,\nand many more tricks."
var Creat = "H.F.M"
var YearCreat = "2018-20"
var LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
var LicenseAbrv = "License (MIT)"
var Repository = "github.com/hfmrow/sAndReplace"

// Vars declarations
var absoluteRealPath, optFilename = getAbsRealPath()
var mainOptions *MainOpt
var tempDir string
var doTempDir bool
var namingWidget bool
var devMode bool
var mainWindowTitle string

var acceptBinary bool
var filesCount int
var textWinTextToShowBytes []byte
var currFilename string

// Declarations lib mapping
var tvsTree,
	tvsList *gitw.TreeViewStructure
var statusbar *gimc.StatusBar
var mainWinTitle, textWinTitle, findWinTitle *gimc.TitleBar
var dnd *gimc.DragNDropStruct
var popMenuTextView *gimc.PopupMenu
var popupMenu *gimc.PopupMenu

// Error handling
var DlgErr = func(dsc string, err error) (yes bool) {
	yes = gidg.DialogError(mainObjects.mainWin, sts["issue"], dsc, err, devMode, true)
	return
}

// Highlight chroma
var ChromaHighlightNew = gitvch.ChromaHighlightNew
var highlighter *gitvch.ChromaHighlight
var textViewRowNumber *gitv.TextViewNumbered

var ComboBoxTextFill = gits.ComboBoxTextFill
var ComboBoxTextAddSetEntry = gits.ComboBoxTextAddSetEntry

// Global decl
var found_background_color_lightgreen = "found_background_color_lightgreen"
var filesFoundMulti []glfsft.SearchAndReplaceFiles
var fileFoundSingle *glfsft.SearchAndReplace
var toDispFileList []glfssf.ScanDirFileInfos

var previsouWordsPos [][]int
var currentInFilesList []string
var forbiddenFiles bool
var forbiddenFilesAlreadyDone bool
var scanTime, searchTime, dispTime string // Some times controls
var fromDnD bool
var lastLine int
var alreadyPlacedPrevWin bool
var alreadyPlacedFoundWin bool
var btnFindInUse bool
var pbs *gimc.ProgressBarStruct
var fileChooserDoUpdt bool
var displayFilesReadyToUpdate bool

var currentText string

type MainOpt struct {
	/* Public, will be saved and restored */
	AboutOptions     *gidg.AboutInfos
	MainWinWidth     int
	MainWinHeight    int
	MainWinPosX      int
	MainWinPosY      int
	LanguageFilename string

	CaseSensitive      bool
	UseRegex           bool
	WholeWord          bool
	ScanSubDir         bool
	FollowSymlinkDir   bool
	MakeBackup         bool
	UseCharacterClass  bool
	UseWildcard        bool
	WrapText           bool
	DispForbiddenFiles bool

	Directory   string
	AppLauncher string

	ExtMask        []string
	ExtSep         string
	CascadeDepth   int
	FilePathLength int
	ScanDirDepth   int

	FileMinSizeLimit int64
	FileMaxSizeLimit int64

	AutoScan                bool
	SyntaxHighlight         bool
	SyntaxHighlightType     string
	SyntaxHighlightLanguage string
	TxtBgCol                string
	TxtFgCol                string
	NumFgCol                string

	/* Private, will NOT be saved */
	listStoreColumns [][]string
	treeStoreColumns [][]string
}

// Main options initialisation
func (opt *MainOpt) Init() {
	opt.AboutOptions = new(gidg.AboutInfos)

	opt.MainWinWidth = 800
	opt.MainWinHeight = 600
	opt.CascadeDepth = 30
	opt.FilePathLength = 4

	opt.LanguageFilename = "assets/lang/eng.lang"
	opt.Directory = ""
	opt.ExtMask = []string{"*"}
	opt.ExtSep = ";"
	opt.listStoreColumns = [][]string{
		{"Name", "text"},
		{"Size", "text"},
		{"Date", "text"},
		{"Path", "markup"},
		{"sizeSort", "int64"}, // This one will be invisible (int64)
		{"dateSort", "int64"}} // This one will be invisible (int64)}
	opt.treeStoreColumns = [][]string{{"Filename(s)", "markup"}}

	opt.ScanDirDepth = -1
	opt.AutoScan = true
	opt.MakeBackup = true

	opt.AppLauncher = "xdg-open"

	opt.SyntaxHighlightType = "hfmrow"
	opt.SyntaxHighlightLanguage = "Go"
	opt.TxtBgCol = "#F8F8F8"
	opt.TxtFgCol = "#1A1A1A"
	opt.NumFgCol = "#008B8B"

	opt.FileMinSizeLimit = 16      // mean 16 bytes
	opt.FileMaxSizeLimit = 1048576 // mean 1M/b
}

// Variables -> Objects.
func (opt *MainOpt) UpdateObjects() {
	mainObjects.mainWin.Resize(opt.MainWinWidth, opt.MainWinHeight)
	mainObjects.mainWin.Move(opt.MainWinPosX, opt.MainWinPosY)

	OptToExtSlice()
	mainObjects.fileChooserBtn.SetCurrentFolder(opt.Directory)
	mainObjects.chkCaseSensitive.SetActive(opt.CaseSensitive)
	mainObjects.chkRegex.SetActive(opt.UseRegex)
	mainObjects.chkWholeWord.SetActive(opt.WholeWord)
	mainObjects.chkWildcard.SetActive(opt.UseWildcard)
	mainObjects.chkCharacterClass.SetActive(opt.UseCharacterClass)
	mainObjects.textWinChkWrap.SetActive(opt.WrapText)
	mainObjects.findWinChkBackUp.SetActive(mainOptions.MakeBackup)
	mainObjects.findWinChkDispForbFiles.SetActive(opt.DispForbiddenFiles)
	mainObjects.chkFollowSymlinkDir.SetActive(opt.FollowSymlinkDir)
	mainObjects.textWinChkSyntxHighlight.SetActive(opt.SyntaxHighlight)

	mainObjects.OptionsEntryMaxFileSize.SetText(fmt.Sprintf("%d", opt.FileMaxSizeLimit))
	mainObjects.OptionsEntryMinFileSize.SetText(fmt.Sprintf("%d", opt.FileMinSizeLimit))
	mainObjects.switchFileChooserButton.SetActive(opt.AutoScan)
	ComboBoxTextAddSetEntry(mainObjects.textWinComboBoxTextStyleChooser, opt.SyntaxHighlightType)
	ComboBoxTextAddSetEntry(mainObjects.textWinComboBoxLanguage, opt.SyntaxHighlightLanguage)

	/* Init spinButton */
	gitl.SpinbuttonSetValues(mainObjects.spinButtonDepth, -1, 128, opt.ScanDirDepth)
}

// Objects -> Variables.
func (opt *MainOpt) UpdateOptions() {
	opt.MainWinWidth, opt.MainWinHeight = mainObjects.mainWin.GetSize()
	opt.MainWinPosX, opt.MainWinPosY = mainObjects.mainWin.GetPosition()

	ExtSliceToOpt() // Store extensions list
	opt.Directory = mainObjects.fileChooserBtn.FileChooser.GetFilename()
	opt.CaseSensitive = mainObjects.chkCaseSensitive.GetActive()
	opt.UseRegex = mainObjects.chkRegex.GetActive()
	opt.WholeWord = mainObjects.chkWholeWord.GetActive()
	opt.UseWildcard = mainObjects.chkWildcard.GetActive()
	opt.UseCharacterClass = mainObjects.chkCharacterClass.GetActive()
	opt.WrapText = mainObjects.textWinChkWrap.GetActive()
	opt.MakeBackup = mainObjects.findWinChkBackUp.GetActive()
	opt.DispForbiddenFiles = mainObjects.findWinChkDispForbFiles.GetActive()
	opt.FollowSymlinkDir = mainObjects.chkFollowSymlinkDir.GetActive()
	opt.SyntaxHighlight = mainObjects.textWinChkSyntxHighlight.GetActive()

	// Value are sanitized, so, we only need to retrieve them
	opt.FileMaxSizeLimit = int64(gitl.GetEntryTextAsInt(mainObjects.OptionsEntryMaxFileSize))
	opt.FileMinSizeLimit = int64(gitl.GetEntryTextAsInt(mainObjects.OptionsEntryMinFileSize))
	opt.AutoScan = mainObjects.switchFileChooserButton.GetActive()
	opt.SyntaxHighlightType = mainObjects.textWinComboBoxTextStyleChooser.GetActiveText()
	opt.SyntaxHighlightLanguage = mainObjects.textWinComboBoxLanguage.GetActiveText()

	opt.ScanDirDepth = mainObjects.spinButtonDepth.GetValueAsInt()
}

// Read Options from file
func (opt *MainOpt) Read() (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(optFilename); err == nil {
		err = json.Unmarshal(textFileBytes, &opt)
	}
	return err
}

// Write Options to file
func (opt *MainOpt) Write() (err error) {
	var out bytes.Buffer
	var jsonData []byte
	opt.UpdateOptions()
	if jsonData, err = json.Marshal(&opt); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(optFilename, out.Bytes(), os.ModePerm)
		}
	}
	return err
}
