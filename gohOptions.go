// gohOptions.go

/*
	Source file auto-generated on Sat, 24 Apr 2021 04:33:11 using Gotk3 Objects Handler v1.7.8
	©2018-21 hfmrow https://hfmrow.github.io

	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-21 hfmrow - Search And Replace v1.10 github.com/hfmrow/search-and-replace

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

	glfs "github.com/hfmrow/genLib/files"
	glfsft "github.com/hfmrow/genLib/files/fileText"
	glfssf "github.com/hfmrow/genLib/files/scanFileDir"
	glss "github.com/hfmrow/genLib/slices"
	glsg "github.com/hfmrow/genLib/strings"
	glts "github.com/hfmrow/genLib/tools"
	gltsbh "github.com/hfmrow/genLib/tools/bench"
	gltsle "github.com/hfmrow/genLib/tools/log2file"
	gltsushe "github.com/hfmrow/genLib/tools/units/human_readable"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gipops "github.com/hfmrow/gtk3Import/pango/pangoSimple"
	gitvsv "github.com/hfmrow/gtk3Import/textView/sourceView"
	gitl "github.com/hfmrow/gtk3Import/tools"
	gits "github.com/hfmrow/gtk3Import/tools"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

// App infos
var (
	Name         = "Search And Replace"
	Vers         = "v1.10"
	Descr        = "Search and replace text sequences in one or more files.\nThe clipboard is also included inside replacing functionality.\nThere is some useful functionalities like code preview with\nsyntax highlighting, backup modified files whether is requested,\nand many more tricks."
	Creat        = "hfmrow"
	YearCreat    = "2018-21"
	LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
	LicenseAbrv  = "License (MIT)"
	Repository   = "github.com/hfmrow/search-and-replace"
)

// Vars declarations
var (
	absoluteRealPath,
	optFilename string

	opt *MainOpt
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

	/*
	 * Lib mapping
	 */
	// AboutBox
	AboutInfosNew = gidg.AboutInfosNew
	About         *gidg.AboutInfos

	tvsTree,
	tvsList *gitw.TreeViewStructure
	TreeViewStructureNew = gitw.TreeViewStructureNew

	statusbar                                *gimc.StatusBar
	mainWinTitle, textWinTitle, findWinTitle *gimc.TitleBar
	dnd                                      *gimc.DragNDropStruct
	DragNDropNew                             = gimc.DragNDropNew
	TitleBarStructureNew                     = gimc.TitleBarStructureNew
	StatusBarStructureNew                    = gimc.StatusBarStructureNew

	// Popup menu
	PopupMenuIconStructNew = gimc.PopupMenuIconStructNew
	popMenuTextView        *gimc.PopupMenuIconStruct
	popupMenu              *gimc.PopupMenuIconStruct

	DialogMessage = gidg.DialogMessage

	GetEntryText      = gitl.GetEntryText
	GetEntryTextAsInt = gitl.GetEntryTextAsInt

	pbs            *gimc.ProgressBarStruct
	ProgressGifNew = gimc.ProgressGifNew

	PangoColorNew = gipops.PangoColorNew

	// Error handling
	DlgErr = func(dsc string, err error) (yes bool) {
		yes = gidg.DialogError(obj.mainWin, sts["issue"], dsc, err, devMode, devMode)
		return
	}
	// Errors handling
	Log2FileStructNew = gltsle.Log2FileStructNew
	Logger            *gltsle.Log2FileStruct

	// GtkSourceView
	svs                 *gitvsv.SourceViewStruct
	SourceViewStructNew = gitvsv.SourceViewStructNew

	// Genlib decl
	filesFoundMulti         []glfsft.SearchAndReplaceFiles
	fileFoundSingle         *glfsft.SearchAndReplace
	SearchAndReplaceNew     = glfsft.SearchAndReplaceNew
	SearchAndReplaceInFiles = glfsft.SearchAndReplaceInFiles
	IsTextFile              = glfsft.IsTextFile

	IsExistSl         = glss.IsExistSl
	IsExistSlIface    = glss.IsExistSlIface
	DeleteSlIface     = glss.DeleteSlIface
	HumanReadableSize = gltsushe.HumanReadableSize
	TruncatePath      = glfs.TruncatePath
	RemoveNonNum      = glsg.RemoveNonNum
	UnEscapedStr      = func(in string) string {
		if obj.chkUseEscapeChar.GetActive() {
			return glsg.UnEscapedStr(in)
		}
		return in
	}
	// UnEscapeString = glsg.UnEscapeString

	BenchNew    = gltsbh.BenchNew
	ExecCommand = glts.ExecCommand

	// Files
	toDispFileList       []glfssf.ScanDirFileInfos
	ScanDirDepthFileInfo = glfssf.ScanDirDepthFileInfo

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
	// File signature
	FileSign []string

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
	ExpandAll,
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
	mapListStore,
	mapTreeStore map[string]int
}

// Main options initialisation
func (opt *MainOpt) Init() {

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
		{"pathReal", "text"},  // This one will be invisible (int64)
		{"sizeSort", "int64"}, // This one will be invisible (int64)
		{"dateSort", "int64"}} // This one will be invisible (int64)}
	opt.mapListStore = map[string]int{
		"Name":     0,
		"Size":     1,
		"Date":     2,
		"PathDisp": 3,
		"pathReal": 4,
		"sizeSort": 5,
		"dateSort": 6}
	opt.treeStoreColumns = [][]string{
		{"", "active"},
		{"Filename(s)", "markup"},
		{"fileIdx", "int64"},
		{"lineIdx", "int64"}}
	opt.mapTreeStore = map[string]int{
		"Toggle":   0,
		"Filename": 1,
		"fileIdx":  2, // This one will be invisible (int64)
		"lineIdx":  3} // This one will be invisible (int64)

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
	obj.mainWin.Resize(opt.MainWinWidth, opt.MainWinHeight)
	obj.mainWin.Move(opt.MainWinPosX, opt.MainWinPosY)

	obj.textWin.Resize(opt.SourceWinWidth, opt.SourceWinHeight)
	obj.textWin.Move(opt.SourceWinPosX, opt.SourceWinPosY)
	obj.Paned.SetPosition(opt.SourceWinWidth - opt.PanedWidth)

	OptToExtSlice()
	obj.fileChooserBtn.SetCurrentFolder(opt.Directory)
	obj.chkCaseSensitive.SetActive(opt.CaseSensitive)
	obj.chkRegex.SetActive(opt.UseRegex)
	obj.chkWholeWord.SetActive(opt.WholeWord)
	obj.chkWildcard.SetActive(opt.UseWildcard)
	obj.chkCharacterClass.SetActive(opt.UseCharacterClass)
	obj.textWinChkWrap.SetActive(opt.WrapText)
	obj.findWinChkBackUp.SetActive(opt.MakeBackup)
	obj.findWinChkDispForbFiles.SetActive(opt.DispForbiddenFiles)
	obj.chkFollowSymlinkDir.SetActive(opt.FollowSymlinkDir)
	obj.textWinChkSyntxHighlight.SetActive(opt.SyntaxHighlight)

	obj.chkUseEscapeChar.SetActive(opt.UseEscapeChar)
	obj.chkUseEscapeCharToReplace.SetActive(opt.UseEscapeCharToReplace)

	obj.OptionsEntryMaxFileSize.SetText(fmt.Sprintf("%d", opt.FileMaxSizeLimit))
	obj.OptionsEntryMinFileSize.SetText(fmt.Sprintf("%d", opt.FileMinSizeLimit))
	obj.switchFileChooserButton.SetActive(opt.AutoScan)
	obj.SourceToggleButtonMapWidth.SetActive(opt.FixedMapWidth)

	obj.findWinChkExpandAll.SetActive(opt.ExpandAll)

	/* Init spinButton */
	gits.SpinScaleSetNew(obj.spinButtonDepth, -1, 128, float64(opt.ScanDirDepth), 1)
}

// Objects -> Variables.
func (opt *MainOpt) UpdateOptions() {
	opt.MainWinWidth, opt.MainWinHeight = obj.mainWin.GetSize()
	opt.MainWinPosX, opt.MainWinPosY = obj.mainWin.GetPosition()

	opt.SourceWinWidth, opt.SourceWinHeight = obj.textWin.GetSize()
	opt.SourceWinPosX, opt.SourceWinPosY = obj.textWin.GetPosition()
	opt.PanedWidth = opt.SourceWinWidth - obj.Paned.GetPosition()
	opt.FixedMapWidth = obj.SourceToggleButtonMapWidth.GetActive()

	ExtSliceToOpt() // Store extensions list
	opt.Directory = obj.fileChooserBtn.FileChooser.GetFilename()
	opt.CaseSensitive = obj.chkCaseSensitive.GetActive()
	opt.UseRegex = obj.chkRegex.GetActive()
	opt.WholeWord = obj.chkWholeWord.GetActive()
	opt.UseWildcard = obj.chkWildcard.GetActive()
	opt.UseCharacterClass = obj.chkCharacterClass.GetActive()
	opt.WrapText = obj.textWinChkWrap.GetActive()
	opt.MakeBackup = obj.findWinChkBackUp.GetActive()
	opt.DispForbiddenFiles = obj.findWinChkDispForbFiles.GetActive()
	opt.FollowSymlinkDir = obj.chkFollowSymlinkDir.GetActive()
	opt.SyntaxHighlight = obj.textWinChkSyntxHighlight.GetActive()

	opt.UseEscapeChar = obj.chkUseEscapeChar.GetActive()
	opt.UseEscapeCharToReplace = obj.chkUseEscapeCharToReplace.GetActive()

	// Value are sanitized, so, we only need to retrieve them
	opt.FileMaxSizeLimit = int64(gits.GetEntryTextAsInt(obj.OptionsEntryMaxFileSize))
	opt.FileMinSizeLimit = int64(gits.GetEntryTextAsInt(obj.OptionsEntryMinFileSize))
	opt.AutoScan = obj.switchFileChooserButton.GetActive()

	opt.ExpandAll = obj.findWinChkExpandAll.GetActive()

	opt.ScanDirDepth = obj.spinButtonDepth.GetValueAsInt()
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
	opt.FileSign = []string{Name, Vers, "©" + YearCreat, Creat, Repository, LicenseAbrv}
	if jsonData, err = json.Marshal(&opt); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(optFilename, out.Bytes(), os.ModePerm)
		}
	}
	return err
}
