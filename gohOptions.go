// gohOptions.go

/*
	Source file auto-generated on Fri, 02 Apr 2021 10:53:33 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 H.F.M - Search And Replace v1.9 github.com/hfmrow/search-and-replace
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
	gitvsv "github.com/hfmrow/gtk3Import/textView/sourceView"
	gits "github.com/hfmrow/gtk3Import/tools"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

// App infos
var (
	Name         = "Search And Replace"
	Vers         = "v1.9"
	Descr        = "Search and replace text sequences in one or more files.\nThe clipboard is also included inside replacing functionality.\nThere is some useful functionalities like code preview with\nsyntax highlighting, backup modified files whether is requested,\nand many more tricks."
	Creat        = "H.F.M"
	YearCreat    = "2018-21"
	LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
	LicenseAbrv  = "License (MIT)"
	Repository   = "github.com/hfmrow/search-and-replace"
)

// Vars declarations
var (
	absoluteRealPath,
	optFilename string

	mainOptions *MainOpt
	tempDir,
	currFilename,
	mainWindowTitle string

	doTempDir,
	namingWidget,
	devMode,
	acceptBinary bool
	filesCount,
	currentLine int
	textWinTextToShowBytes []byte

	// Lib mapping
	tvsTree,
	tvsList *gitw.TreeViewStructure

	statusbar                                *gimc.StatusBar
	mainWinTitle, textWinTitle, findWinTitle *gimc.TitleBar
	dnd                                      *gimc.DragNDropStruct

	// Popup menu
	PopupMenuIconStructNew = gimc.PopupMenuIconStructNew
	popMenuTextView        *gimc.PopupMenuIconStruct
	popupMenu              *gimc.PopupMenuIconStruct

	DialogMessage = gidg.DialogMessage

	pbs *gimc.ProgressBarStruct

	// Error handling
	DlgErr = func(dsc string, err error) (yes bool) {
		yes = gidg.DialogError(mainObjects.mainWin, sts["issue"], dsc, err, devMode, devMode)
		return
	}

	// GtkSourceView
	svs                 *gitvsv.SourceViewStruct
	SourceViewStructNew = gitvsv.SourceViewStructNew

	// Genlib decl
	filesFoundMulti []glfsft.SearchAndReplaceFiles
	fileFoundSingle *glfsft.SearchAndReplace
	toDispFileList  []glfssf.ScanDirFileInfos

	// Global vars
	previsouWordsPos   [][]int
	currentInFilesList []string
	scanTime,
	searchTime,
	currentText,
	dispTime string // Some times controls

	currentTextChanged,
	forbiddenFiles,
	forbiddenFilesAlreadyDone,
	fromDnD,
	alreadyPlacedPrevWin,
	alreadyPlacedFoundWin,
	btnFindInUse,
	fileChooserDoUpdt,
	displayFilesReadyToUpdate bool
)

type MainOpt struct {
	/* Public, will be saved and restored */
	AboutOptions *gidg.AboutInfos
	MainWinWidth,
	MainWinHeight,
	MainWinPosX,
	MainWinPosY,
	PanedWidth,
	SourceWinWidth,
	SourceWinHeight,
	SourceWinPosX,
	SourceWinPosY int

	FixedMapWidth bool

	LanguageFilename string

	CaseSensitive,
	UseRegex,
	WholeWord,
	ScanSubDir,
	FollowSymlinkDir,
	MakeBackup,
	UseCharacterClass,
	UseWildcard,
	WrapText,
	DispForbiddenFiles,
	UseEscapeChar,
	UseEscapeCharToReplace,
	AutoScan,
	SyntaxHighlight bool

	ExtMask []string

	ExtSep,
	Directory,
	AppLauncher,
	DefaultSourceLang,
	DefaultSourceStyle,
	HighlightUserDefined,
	TxtBgCol,
	TxtFgCol string

	CascadeDepth,
	FilePathLength,
	ScanDirDepth int

	FileMinSizeLimit,
	FileMaxSizeLimit int64

	/* Private, will NOT be saved */
	listStoreColumns,
	treeStoreColumns [][]string
}

// Main options initialisation
func (opt *MainOpt) Init() {
	opt.AboutOptions = new(gidg.AboutInfos)

	opt.MainWinWidth = 800
	opt.MainWinHeight = 600
	opt.CascadeDepth = 30
	opt.FilePathLength = 4

	opt.PanedWidth = 120
	opt.SourceWinWidth = 800
	opt.SourceWinHeight = 480
	opt.SourceWinPosX = -1
	opt.SourceWinPosY = -1

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

	opt.DefaultSourceLang = "go-hfmrow"
	opt.DefaultSourceStyle = "hfmrow"
	opt.HighlightUserDefined = "assets/langAndstyle"
	opt.SyntaxHighlight = true

	opt.TxtBgCol = "#F8F8F8"
	opt.TxtFgCol = "#1A1A1A"

	opt.FileMinSizeLimit = 16      // means 16 bytes
	opt.FileMaxSizeLimit = 1048576 // means 1M/b
}

// Variables -> Objects.
func (opt *MainOpt) UpdateObjects() {
	mainObjects.mainWin.Resize(opt.MainWinWidth, opt.MainWinHeight)
	mainObjects.mainWin.Move(opt.MainWinPosX, opt.MainWinPosY)

	mainObjects.textWin.Resize(opt.SourceWinWidth, opt.SourceWinHeight)
	mainObjects.textWin.Move(opt.SourceWinPosX, opt.SourceWinPosY)
	mainObjects.Paned.SetPosition(opt.SourceWinWidth - opt.PanedWidth)

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

	mainObjects.chkUseEscapeChar.SetActive(opt.UseEscapeChar)
	mainObjects.chkUseEscapeCharToReplace.SetActive(opt.UseEscapeCharToReplace)

	mainObjects.OptionsEntryMaxFileSize.SetText(fmt.Sprintf("%d", opt.FileMaxSizeLimit))
	mainObjects.OptionsEntryMinFileSize.SetText(fmt.Sprintf("%d", opt.FileMinSizeLimit))
	mainObjects.switchFileChooserButton.SetActive(opt.AutoScan)
	mainObjects.SourceToggleButtonMapWidth.SetActive(opt.FixedMapWidth)

	/* Init spinButton */
	gits.SpinScaleSetNew(mainObjects.spinButtonDepth, -1, 128, float64(opt.ScanDirDepth), 1)
}

// Objects -> Variables.
func (opt *MainOpt) UpdateOptions() {
	opt.MainWinWidth, opt.MainWinHeight = mainObjects.mainWin.GetSize()
	opt.MainWinPosX, opt.MainWinPosY = mainObjects.mainWin.GetPosition()

	opt.SourceWinWidth, opt.SourceWinHeight = mainObjects.textWin.GetSize()
	opt.SourceWinPosX, opt.SourceWinPosY = mainObjects.textWin.GetPosition()
	opt.PanedWidth = opt.SourceWinWidth - mainObjects.Paned.GetPosition()
	opt.FixedMapWidth = mainObjects.SourceToggleButtonMapWidth.GetActive()

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

	opt.UseEscapeChar = mainObjects.chkUseEscapeChar.GetActive()
	opt.UseEscapeCharToReplace = mainObjects.chkUseEscapeCharToReplace.GetActive()

	// Value are sanitized, so, we only need to retrieve them
	opt.FileMaxSizeLimit = int64(gits.GetEntryTextAsInt(mainObjects.OptionsEntryMaxFileSize))
	opt.FileMinSizeLimit = int64(gits.GetEntryTextAsInt(mainObjects.OptionsEntryMinFileSize))
	opt.AutoScan = mainObjects.switchFileChooserButton.GetActive()

	opt.ScanDirDepth = mainObjects.spinButtonDepth.GetValueAsInt()
}

// Read Options from file
func (opt *MainOpt) Read() (err error) {
	var textFileBytes []byte
	opt.Init() // Init with default values.
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

	opt.AboutOptions.DlgBoxStruct = nil // remove dialog object before saving

	if jsonData, err = json.Marshal(&opt); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(optFilename, out.Bytes(), os.ModePerm)
		}
	}
	return err
}
