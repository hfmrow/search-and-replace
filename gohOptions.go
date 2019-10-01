// gohOptions.go

// Source file auto-generated on Tue, 01 Oct 2019 16:36:13 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/gotk3/gotk3/glib"

	glft "github.com/hfmrow/genLib/files/fileText"

	gidgat "github.com/hfmrow/gtk3Import/dialog/about"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gitv "github.com/hfmrow/gtk3Import/textView"
	gits "github.com/hfmrow/gtk3Import/tools"
	gitw "github.com/hfmrow/gtk3Import/treeview"
)

// App infos
var Name = "Search And Replace"
var Vers = "v1.7.6"
var Descr = "Search and replace text sequences in one or more files.\nThe clipboard is also included inside replacing functionality."
var Creat = "H.F.M"
var YearCreat = "2018-19"
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

var applyChanges, acceptBinary bool
var filesCount int
var textWinTextToShowBytes []byte

var tvsTree,
	tvsList *gitw.TreeViewStructure
var statusbar *gimc.StatusBar
var dnd *gimc.DragNDropStruct

var toDispFileList []string
var currentInFilesList []string
var forbiddenFiles bool
var forbiddenFilesAlreadyDone bool
var scanTime, searchTime, dispTime string // Some times controls
var fromDnD bool

var fileChooserSigHandlerSelChanged,
	spinButtonDepthSigHandlerChanged glib.SignalHandle
var fileChooserDoUpdt bool

var displayFilesReadyToUpdate bool
var mainWinTitle, textWinTitle, findWinTitle *gimc.TitleBar

type MainOpt struct {
	/* Public, will be saved and restored */
	AboutOptions     *gidgat.AboutInfos
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

	ScanButtonEnabled bool

	Directory string

	ExtMask        []string
	ExtSep         string
	CascadeDepth   int
	FilePathLength int
	ScanDirDepth   int

	FileSizeLimit      int64
	LineEndThreshold   float64
	OverCharsThreshold float64

	/* Private, will NOT be saved */
	mainFound         []glft.Find_s
	textViewRowNumber *gitv.TextViewRowNumber
	listStoreColumns  [][]string
	treeStoreColumns  [][]string
}

// Main options initialisation
func (opt *MainOpt) Init() {
	opt.AboutOptions = new(gidgat.AboutInfos)

	opt.MainWinWidth = 800
	opt.MainWinHeight = 600
	opt.CascadeDepth = 30
	opt.FilePathLength = 4

	opt.LanguageFilename = "assets/lang/eng.lang"
	opt.CaseSensitive = false
	opt.Directory = ""
	opt.FollowSymlinkDir = false
	opt.ExtMask = []string{"*"}
	opt.ExtSep = ";"
	opt.ScanSubDir = false
	opt.listStoreColumns = [][]string{{"Name", "text"}, {"Size", "text"}, {"Date", "text"}, {"Path", "markup"}}
	opt.treeStoreColumns = [][]string{{"Filename(s)", "markup"}}
	opt.UseRegex = false
	opt.WholeWord = false
	opt.MakeBackup = true
	opt.UseCharacterClass = false
	opt.UseWildcard = false
	opt.WrapText = false

	opt.ScanButtonEnabled = true
	opt.ScanDirDepth = 0

	opt.FileSizeLimit = 16 // mean 16 bytes
	opt.LineEndThreshold = 0.6
	opt.OverCharsThreshold = 5 // mean 5%
}

// Variables -> Objects.
func (opt *MainOpt) UpdateObjects() {
	mainObjects.mainWin.Resize(opt.MainWinWidth, opt.MainWinHeight)
	mainObjects.mainWin.Move(opt.MainWinPosX, opt.MainWinPosY)

	OptToExtSlice()
	mainObjects.fileChooserBtn.SetCurrentFolder(opt.Directory)
	mainObjects.chkCaseSensitive.SetActive(opt.CaseSensitive)
	mainObjects.chkRegex.SetActive(opt.UseRegex)
	mainObjects.chkWoleWord.SetActive(opt.WholeWord)
	mainObjects.chkWildcard.SetActive(opt.UseWildcard)
	mainObjects.chkCharacterClass.SetActive(opt.UseCharacterClass)
	mainObjects.textWinChkWrap.SetActive(opt.WrapText)
	mainObjects.findWinChkBackUp.SetActive(mainOptions.MakeBackup)
	mainObjects.findWinChkDispForbFiles.SetActive(opt.DispForbiddenFiles)
	mainObjects.chkFollowSymlinkDir.SetActive(opt.FollowSymlinkDir)

	/* Init spinButton */
	gits.SpinbuttonSetValues(mainObjects.spinButtonDepth, -1, 128, opt.ScanDirDepth)
}

// Objects -> Variables.
func (opt *MainOpt) UpdateOptions() {
	opt.MainWinWidth, opt.MainWinHeight = mainObjects.mainWin.GetSize()
	opt.MainWinPosX, opt.MainWinPosY = mainObjects.mainWin.GetPosition()

	ExtSliceToOpt() // Store extensions list
	opt.Directory = mainObjects.fileChooserBtn.FileChooser.GetFilename()
	opt.CaseSensitive = mainObjects.chkCaseSensitive.GetActive()
	opt.UseRegex = mainObjects.chkRegex.GetActive()
	opt.WholeWord = mainObjects.chkWoleWord.GetActive()
	opt.UseWildcard = mainObjects.chkWildcard.GetActive()
	opt.UseCharacterClass = mainObjects.chkCharacterClass.GetActive()
	opt.WrapText = mainObjects.textWinChkWrap.GetActive()
	opt.MakeBackup = mainObjects.findWinChkBackUp.GetActive()
	opt.DispForbiddenFiles = mainObjects.findWinChkDispForbFiles.GetActive()
	opt.FollowSymlinkDir = mainObjects.chkFollowSymlinkDir.GetActive()

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
			err = ioutil.WriteFile(optFilename, out.Bytes(), 0644)
		}
	}
	return err
}
